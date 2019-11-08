package commands_test

import (
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexkarlov/15x4bot/commands"
	"github.com/alexkarlov/15x4bot/config"
	"github.com/alexkarlov/15x4bot/lang"
	"github.com/alexkarlov/15x4bot/store"
	local_time "github.com/alexkarlov/15x4bot/time"
	"github.com/antonmashko/envconf"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

var mock sqlmock.Sqlmock

func TestMain(m *testing.M) {
	var conf config.Config
	envconf.Parse(&conf)
	commands.Conf = conf.Chat
	store.Conf = conf.DB
	var err error
	mock, err = store.InitTest()
	if err != nil {
		panic("error while initing db test ")
	}
	m.Run()
}

type regexpArg struct {
	p *regexp.Regexp
}

func newRegexpArg(p string) regexpArg {
	r := regexpArg{}
	r.p = regexp.MustCompile(p)
	return r
}

func (r regexpArg) Match(v driver.Value) bool {
	s, ok := v.(string)
	if !ok {
		return false
	}
	return r.p.MatchString(s)
}

// test happy path
func createRehearsal(t *testing.T) {
	u := &store.User{
		Role: store.USER_ROLE_ADMIN,
	}
	cmdText := lang.MARKUP_BUTTON_CREATE_REHEARSAL
	typeCmd := "addRehearsal"
	cmd := commands.NewCommand(cmdText, u)
	tCmd := reflect.TypeOf(cmd)
	if tCmd.Elem().Name() != typeCmd {
		t.Fatalf("wrong type of command %s. exepected %s, actual %s", cmdText, typeCmd, tCmd.Elem().Name())
	}
	// test first step
	m, err := cmd.NextStep("")
	if err != nil {
		t.Fatal("unexpected error on first step: ", err)
	}
	// on the first step we ask rehearsal time start
	if m.Text != lang.ADD_REHEARSAL_WHEN {
		t.Errorf("wrong text of first cmd. expected %s actual %s", m.Text, lang.ADD_REHEARSAL_WHEN)
	}

	// test second step
	// prepare mock rows for places
	rows := sqlmock.NewRows([]string{"id", "name", "address"}).AddRow(1, "name1", "address1").AddRow(2, "name2", "address2")
	mock.ExpectQuery("SELECT id, name, address FROM places").WillReturnRows(rows)
	tR := "2019-08-25 19:00"
	// send start time
	m, err = cmd.NextStep(tR)
	if err != nil {
		t.Fatal("unexpected error on second step: ", err)
	}
	// on the second step we ask a rehersal place
	if m.Text != lang.PLACES_CHOSE_PLACE {
		t.Errorf("wrong text of second cmd. expected %s actual %s", lang.PLACES_CHOSE_PLACE, m.Text)
	}
	// check that we send appropriate count of places buttons
	if len(m.Buttons) != 2 {
		t.Errorf("wrong count of buttons on second step: %#v %#v", m.Buttons, rows)
	}

	// test third step
	// prepare mock row for a place
	rows = sqlmock.NewRows([]string{"id"}).AddRow(3)
	mock.ExpectQuery("SELECT id FROM places").WithArgs(2).WillReturnRows(rows)
	// prepare mock time start
	tRArg, _ := time.Parse(commands.Conf.TimeLayout, tR)
	// we expect place "2" because we pass it in the "NextStep"
	mock.ExpectQuery("INSERT INTO rehearsals \\(time, place\\)").
		WithArgs(tRArg, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	m, err = cmd.NextStep("Місце 2: name2")
	if err != nil {
		t.Fatal("unexpected error on third step: ", err)
	}
	// we expect that rehearsal will be added successfuly
	if m.Text != lang.ADD_REHEARSAL_SUCCESS_MSG {
		t.Errorf("unexpected reply text for creation rehearsal. expected %s actual %s", lang.ADD_REHEARSAL_SUCCESS_MSG, m.Text)
	}
	// at this step we ask when bot should send rehearsal reminder
	if len(m.Buttons) != 3 {
		t.Errorf("wrong count of buttons on third step: %#v expected 3", m.Buttons)
	}

	// mock time and pass it further using "local_time.SetNow"
	fakeTime := "2019-09-03 10:00"
	nF := func() time.Time {
		l, _ := time.LoadLocation(commands.Conf.Location)
		t, _ := time.ParseInLocation(commands.Conf.TimeLayout, fakeTime, l)
		return t
	}
	local_time.SetNow(nF)
	// mock rehearsal row
	mock.ExpectQuery("SELECT r.id, r.time, p.name, p.address, p.map_url FROM rehearsals r LEFT JOIN places p ON p.id = r.place WHERE r.id=").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "time", "name", "address", "map_url"}).AddRow(1, tRArg, "name2", "addr1", "url1"))
	// prepare regexp for avoid fragile tests
	// command creates two reminders - for the org chat and org channel
	// mock both tasks
	arg := newRegexpArg("name2.*2019-08-25 19:00.*addr1.*url1.*ChannelUsername.*-389484898")
	mock.ExpectExec("INSERT INTO tasks \\(type, execution_time, status, details\\) VALUES").
		WithArgs(store.TASK_TYPE_REMINDER_TG_CHANNEL, nF(), store.TASK_STATUS_NEW, arg).WillReturnResult(sqlmock.NewResult(1, 1))
	arg = newRegexpArg("name2.*2019-08-25 19:00.*addr1.*url1.*ChannelUsername.*@test15x4")
	mock.ExpectExec("INSERT INTO tasks \\(type, execution_time, status, details\\) VALUES").
		WithArgs(store.TASK_TYPE_REMINDER_TG_CHANNEL, nF(), store.TASK_STATUS_NEW, arg).WillReturnResult(sqlmock.NewResult(1, 1))
	m, err = cmd.NextStep(lang.MARKUP_BUTTON_NOTIF_REHEARSAL_NOW)
	if err != nil {
		t.Fatal("unexpected error on fourth step: ", err)
	}
	if m.Text != fmt.Sprintf(lang.ADD_REHEARSAL_REMINDER_OK, fakeTime) {
		t.Error("wrong success msg: ", m.Text)
	}
}

func nextRehearsal(t *testing.T) {
	u := &store.User{
		Role: store.USER_ROLE_GUEST,
	}
	cmdText := lang.MARKUP_BUTTON_NEXT_REHEARSAL
	typeCmd := "nextRehearsal"
	cmd := commands.NewCommand(cmdText, u)
	tCmd := reflect.TypeOf(cmd)
	if tCmd.Elem().Name() != typeCmd {
		t.Fatalf("wrong type of command %s. exepected %s, actual %s", cmdText, typeCmd, tCmd.Elem().Name())
	}
	// test first step
	// prepare mock rows for places
	fakeTime := "2019-09-03 10:00"
	rt, _ := time.Parse(commands.Conf.TimeLayout, fakeTime)
	rows := sqlmock.NewRows([]string{"time", "name", "address", "map_url"}).AddRow(rt, "name1", "address1", "map_url1")
	mock.ExpectQuery("SELECT r.time, p.name, p.address, p.map_url").WillReturnRows(rows)
	m, err := cmd.NextStep("")
	if err != nil {
		t.Fatal("unexpected error on fist step: ", err)
	}
	r := "(?m)name1[\\S\\s]*address1[\\S\\s]*" + fakeTime + "[\\S\\s]*map_url1"
	if !regexp.MustCompile(r).MatchString(m.Text) {
		t.Errorf("wrong text of next rehearsal, regexp %s text %s", r, m.Text)
	}
}

func deleteRehearsal(t *testing.T) {
	u := &store.User{
		Role: store.USER_ROLE_ADMIN,
	}
	cmdText := lang.MARKUP_BUTTON_DELETE_REHEARSAL
	typeCmd := "deleteRehearsal"
	cmd := commands.NewCommand(cmdText, u)
	tCmd := reflect.TypeOf(cmd)
	if tCmd.Elem().Name() != typeCmd {
		t.Fatalf("wrong type of command %s. exepected %s, actual %s", cmdText, typeCmd, tCmd.Elem().Name())
	}
	// test first step
	// prepare mock rows for rehearsals
	rt, _ := time.Parse(commands.Conf.TimeLayout, "2019-09-03 19:00")
	rows := sqlmock.NewRows([]string{"id", "time", "name"}).AddRow(1, rt, "Марс").AddRow(2, rt, "Венера")
	mock.ExpectQuery("SELECT r.id, r.time, p.name FROM rehearsals r LEFT JOIN places p ON p.id=r.place").WillReturnRows(rows)
	m, err := cmd.NextStep("")
	if err != nil {
		t.Fatal("unexpected error on first step: ", err)
	}
	if len(m.Buttons) != 3 {
		t.Error("wrong count of rehearsals: ", len(m.Buttons))
	}

	// test second step
	r := "Репетиція 1: де: Марс, коли: 3000-01-01"
	mock.ExpectExec("DELETE FROM rehearsals WHERE id=").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	m, err = cmd.NextStep(r)
	if err != nil {
		t.Fatal("unexpected error on second step: ", err)
	}
	if m.Text != lang.DELETE_REHEARSAL_COMPLETE {
		t.Errorf("wrong text of succcessfuly deteling rehearsal: %s", m.Text)
	}
}

func TestCommandAddRehearsalAdmin(t *testing.T) {
	createRehearsal(t)
}
func TestCommandNextRehearsal(t *testing.T) {
	nextRehearsal(t)
}

func TestCommandDeleteRehearsal(t *testing.T) {
	deleteRehearsal(t)
}

func updateLecture(t *testing.T) {
	u := &store.User{
		Role: store.USER_ROLE_ADMIN,
	}
	cmdText := lang.MARKUP_BUTTON_UPDATE_LECTURE
	typeCmd := "upsertLecture"
	rows := sqlmock.NewRows([]string{"l.id", "l.name", "l.description", "u.name", "u.username", "tg_id", "u.id", "u.role"}).
		AddRow(1, "лекція 1", "опис лекції1", "юзер 1", "username1", 123, 1, "admin").
		AddRow(2, "лекція 2", "опис лекції2", "юзер 2", "username2", 124, 2, "lector")
	mock.ExpectQuery("SELECT l.id, l.name, l.description, u.name, u.username, u.tg_id, u.id, u.role FROM lectures l").WillReturnRows(rows)
	cmd := commands.NewCommand(cmdText, u)
	tCmd := reflect.TypeOf(cmd)
	if tCmd.Elem().Name() != typeCmd {
		t.Fatalf("wrong type of command %s. exepected %s, actual %s", cmdText, typeCmd, tCmd.Elem().Name())
	}
	m, err := cmd.NextStep("")
	if err != nil {
		t.Fatal("unexpected error on first step: ", err)
	}
	if len(m.Buttons) != 3 {
		t.Errorf("wrong count of buttons: %#v", m.Buttons)
	}
	if m.Text != lang.ADD_LECTURE_DESCIRPTION_CHOSE_LECTURE {
		t.Error("wrong text on first step: ", m.Text)
	}

	// test second step
	m, err = cmd.NextStep("Лекція 1: лекція 1")
	if err != nil {
		t.Fatal("unexpected error on second step: ", err)
	}
	if m.Text != lang.UPSERT_LECTURE_STEP_LECTURE_NAME {
		t.Error("wrong text on second step: ", m.Text)
	}

	// test third step
	m, err = cmd.NextStep("нова назва")
	if err != nil {
		t.Fatal("unexpected error on third step: ", err)
	}
	if m.Text != lang.UPSERT_LECTURE_STEP_LECTURE_DESCRIPTION {
		t.Error("wrong text on second step: ", m.Text)
	}

	// test fourth step
	loc, _ := time.LoadLocation(commands.Conf.Location)
	nD := "2019-09-04 19:00"
	nDay, _ := time.ParseInLocation(commands.Conf.TimeLayout, nD, loc)
	mock.ExpectExec("UPDATE lectures set name=").
		WithArgs("нова назва", "", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	dR := newRegexpArg("1")
	mock.ExpectExec("INSERT INTO tasks \\(type, execution_time, status, details\\)").
		WithArgs(store.TASK_TYPE_REMINDER_LECTOR, nDay, store.TASK_STATUS_NEW, dR).
		WillReturnResult(sqlmock.NewResult(1, 1))
	fakeTime := "2019-09-03 10:00"
	nF := func() time.Time {
		t, _ := time.ParseInLocation(commands.Conf.TimeLayout, fakeTime, loc)
		return t
	}
	local_time.SetNow(nF)
	m, err = cmd.NextStep(lang.I_DONT_KNOW)
	if err != nil {
		t.Fatal("unexpected error on fourth step: ", err)
	}
	if m.Text != fmt.Sprintf(lang.UPSERT_LECTURE_I_WILL_REMIND, nD) {
		t.Error("wrong text on third step: ", m.Text)
	}

	// test that this is the end
	if !cmd.IsEnd() {
		t.Fatal("should be ended")
	}
}

func insertLecture(t *testing.T) {
	u := &store.User{
		Role: store.USER_ROLE_ADMIN,
	}
	cmdText := lang.MARKUP_BUTTON_CREATE_LECTURE
	typeCmd := "upsertLecture"
	rt, _ := time.Parse(commands.Conf.TimeLayout, "2019-09-03 19:00")
	rows := sqlmock.NewRows([]string{"id", "tg_id", "username", "name", "role", "fb", "vk", "picture_id", "bdate"}).
		AddRow(1, "123", "username1", "юзер 1", "admin", "fbbb", "vk11", "pic1", rt).
		AddRow(2, "321", "username2", "юзер 2", "lector", "fb11", "vk11", "pic2", rt)
	mock.ExpectQuery("SELECT id, tg_id, username, name, role, fb, vk, picture_id, bdate FROM users").WillReturnRows(rows)
	cmd := commands.NewCommand(cmdText, u)
	tCmd := reflect.TypeOf(cmd)
	if tCmd.Elem().Name() != typeCmd {
		t.Fatalf("wrong type of command %s. exepected %s, actual %s", cmdText, typeCmd, tCmd.Elem().Name())
	}
	m, err := cmd.NextStep("")
	if err != nil {
		t.Fatal("unexpected error on first step: ", err)
	}
	if len(m.Buttons) != 1 {
		t.Errorf("wrong count of buttons: %#v", m.Buttons)
	}
	if strings.Contains(m.Text, lang.UPSERT_LECTURE_STEP_SPEAKER) {
		t.Error("wrong text on first step: ", m.Text)
	}

	// test second step
	// pass user id
	rows = sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("SELECT id FROM users WHERE id=").WillReturnRows(rows)

	m, err = cmd.NextStep("1")
	if err != nil {
		t.Fatal("unexpected error on second step: ", err)
	}
	if m.Text != lang.UPSERT_LECTURE_STEP_LECTURE_NAME {
		t.Error("wrong text on second step: ", m.Text)
	}

	// test third step
	m, err = cmd.NextStep("нова назва")
	if err != nil {
		t.Fatal("unexpected error on third step: ", err)
	}
	if m.Text != lang.UPSERT_LECTURE_STEP_LECTURE_DESCRIPTION {
		t.Error("wrong text on third step: ", m.Text)
	}

	// test fourth step
	mock.ExpectQuery("INSERT INTO lectures \\(name, description, user_id\\) VALUES \\(").
		WithArgs("нова назва", "опис", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	m, err = cmd.NextStep("опис")
	if err != nil {
		t.Fatal("unexpected error on fourth step: ", err)
	}
	if m.Text != lang.UPSERT_LECTURE_SEND_TO_GRAMMAR_NAZI {
		t.Error("wrong text on fourth step: ", m.Text)
	}

	m, err = cmd.NextStep(lang.MARKUP_BUTTON_NO)
	if m.Text != lang.UPSERT_LECTURE_SUCCESS_CREATE_MSG {
		t.Error("wrong text on fifth step: ", m.Text)
	}
	// test that this is the end
	if !cmd.IsEnd() {
		t.Fatal("should be ended")
	}
}

func TestCommandUpsertLecture(t *testing.T) {
	updateLecture(t)
	insertLecture(t)
}

func TestNewCommand(t *testing.T) {
	u := &store.User{
		Role: store.USER_ROLE_GUEST,
	}
	typeCmd := "unknown"
	cmdText := "hello"
	cmd := commands.NewCommand(cmdText, u)
	tCmd := reflect.TypeOf(cmd)
	if tCmd.Elem().Name() != typeCmd {
		t.Errorf("wrong type of command %s. exepected %s, actual %s", cmdText, typeCmd, tCmd.Elem().Name())
	}
}

func TestAddDescriptionLecture(t *testing.T) {
	u := &store.User{
		Role: store.USER_ROLE_ADMIN,
	}
	cmdText := lang.MARKUP_BUTTON_ADD_DESCRIPTION
	typeCmd := "addDescriptionLecture"
	cmd := commands.NewCommand(cmdText, u)
	tCmd := reflect.TypeOf(cmd)
	if tCmd.Elem().Name() != typeCmd {
		t.Fatalf("wrong type of command %s. exepected %s, actual %s", cmdText, typeCmd, tCmd.Elem().Name())
	}
	rows := sqlmock.NewRows([]string{"l.id", "l.name", "l.description", "u.name", "u.username", "tg_id", "u.id", "u.role"}).
		AddRow(1, "лекція 1", "опис лекції1", "юзер 1", "username1", 123, 1, "admin").
		AddRow(2, "лекція 2", "", "юзер 2", "username2", 124, 2, "lector").
		AddRow(2, "лекція 3", "", "юзер 3", "username2", 125, 2, "lector")
	mock.ExpectQuery("SELECT l.id, l.name, l.description, u.name, u.username, u.tg_id, u.id, u.role FROM lectures l INNER JOIN users u ON u.id = user_id AND l.id NOT IN \\(SELECT id_lecture FROM event_lectures\\)").WillReturnRows(rows)

	m, err := cmd.NextStep("")
	if err != nil {
		t.Fatal("unexpected error on first step: ", err)
	}
	if len(m.Buttons) != 3 {
		t.Errorf("wrong count of buttons: %#v", m.Buttons)
	}
	if m.Text != lang.ADD_LECTURE_DESCIRPTION_CHOSE_LECTURE {
		t.Error("wrong text on first step: ", m.Text)
	}

	// test second step
	rows = sqlmock.NewRows([]string{"l.id", "l.name", "l.description", "u.id", "u.username", "tg_id", "u.role"}).
		AddRow(1, "лекція 2", "", 124, "username2", 1, "lector")
	mock.ExpectQuery("SELECT l.id, l.name, l.description, u.id, u.username, u.tg_id, u.role FROM lectures l LEFT JOIN users u ON u.id=l.user_id WHERE l.id=").
		WithArgs(2).
		WillReturnRows(rows)

	m, err = cmd.NextStep("Лекція 2: лекція 2")
	if err != nil {
		t.Fatal("unexpected error on second step: ", err)
	}
	if m.Text != lang.UPSERT_LECTURE_STEP_LECTURE_DESCRIPTION {
		t.Error("wrong text on second step: ", m.Text)
	}

	mock.ExpectExec("UPDATE lectures SET description=\\$1 WHERE id=\\$2").
		WithArgs("опис лекції 2", 2).
		WillReturnResult(sqlmock.NewResult(1, 1))
	m, err = cmd.NextStep("опис лекції 2")
	if err != nil {
		t.Fatal("unexpected error on third step: ", err)
	}
	if m.Text != lang.UPSERT_LECTURE_SEND_TO_GRAMMAR_NAZI {
		t.Error("wrong text on third step: ", m.Text)
	}

	m, err = cmd.NextStep(lang.MARKUP_BUTTON_NO)
	if err != nil {
		t.Fatal("unexpected error on fourth step: ", err)
	}
	if m.Text != lang.ADD_LECTURE_DESCRIPTION_COMPLETE {
		t.Error("wrong text on fourth step: ", m.Text)
	}
}
