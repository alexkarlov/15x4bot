package commands

import (
	"encoding/json"
	"fmt"
	"github.com/alexkarlov/15x4bot/lang"
	local_time "github.com/alexkarlov/15x4bot/time"
	"time"

	"github.com/alexkarlov/15x4bot/store"
)

type addRehearsal struct {
	ID    int
	when  time.Time
	where int
	u     *store.User
	stepConstructor
}

// newAddRehearsal creates addRehearsal and registers all steps
func newAddRehearsal() *addRehearsal {
	c := &addRehearsal{}
	c.RegisterSteps(c.firstStep, c.secondStep, c.thirdStep, c.fourthStep)
	return c
}

func (c *addRehearsal) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

func (c *addRehearsal) firstStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
		Text:    lang.ADD_REHEARSAL_WHEN,
	}
	return replyMarkup, nil
}

// secondStep sends a list of places and asks a place for rehearsal
func (c *addRehearsal) secondStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	t, err := time.Parse(Conf.TimeLayout, answer)
	if err != nil {
		replyMarkup.Text = lang.WRONG_DATE_TIME
		return replyMarkup, nil
	}
	c.when = t
	places, err := store.Places(store.PlaceTypes{store.PLACE_TYPE_FOR_REHEARSAL, store.PLACE_TYPE_FOR_ALL})
	if err != nil {
		return nil, err
	}
	replyMarkup.Buttons = nil
	for _, p := range places {
		b := fmt.Sprintf(lang.PLACES_LIST_BUTTONS, p.ID, p.Name)
		replyMarkup.Buttons = append(replyMarkup.Buttons, b)
	}
	replyMarkup.Text = lang.PLACES_CHOSE_PLACE
	return replyMarkup, nil
}

// thirdStep parses place and saves rehearsal into db
// also, bot asks how to send notification to the internal chat and channel
func (c *addRehearsal) thirdStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	var err error
	c.where, err = parseID(answer)
	if err != nil {
		return nil, err
	}
	ok, err := store.DoesPlaceExist(c.where)
	if err != nil {
		return nil, err
	}
	if !ok {
		replyMarkup.Text = lang.WRONG_PLACE_ID
		return replyMarkup, nil
	}
	c.ID, err = store.AddRehearsal(c.when, c.where)
	if err != nil {
		return nil, err
	}
	replyMarkup.Text = lang.ADD_REHEARSAL_SUCCESS_MSG
	replyMarkup.Buttons = MessageButtons{lang.MARKUP_BUTTON_NOTIF_REHEARSAL_NOW, lang.MARKUP_BUTTON_NOTIF_BEFORE_REHEARSAL, lang.MARKUP_BUTTON_MAIN_MENU}
	return replyMarkup, nil
}

// fourthStep adds reminder for sending notification in the chat and channel
func (c *addRehearsal) fourthStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	var err error
	var execTime time.Time
	if answer == lang.MARKUP_BUTTON_NOTIF_REHEARSAL_NOW {
		execTime = asSoonAsPossible()
	} else {
		execTime = c.when.AddDate(0, 0, -1)
	}
	// create a task for sending post in the internal channel
	err = c.addRehearsalReminder(execTime)
	if err != nil {
		return nil, err
	}
	replyMarkup.Text = fmt.Sprintf(lang.ADD_REHEARSAL_REMINDER_OK, execTime.Format(Conf.TimeLayout))
	replyMarkup.Buttons = StandardMarkup(c.u.Role)
	return replyMarkup, nil
}

func (c *addRehearsal) addRehearsalReminder(execTime time.Time) error {
	// create chat reminder
	r, err := store.LoadRehearsal(c.ID)
	if err != nil {
		return err
	}
	msg := ""
	// create chat reminder
	if Conf.OrgChatID != "" {
		msg, err = c.msgToChannel(r, Conf.OrgChatID)
		if err != nil {
			return err
		}
	}
	err = store.AddTask(store.TASK_TYPE_REMINDER_TG_CHANNEL, execTime, msg)
	if err != nil {
		return err
	}
	// create channel reminder
	msg, err = c.msgToChannel(r, Conf.OrgChannelUsername)
	if err != nil {
		return err
	}
	return store.AddTask(store.TASK_TYPE_REMINDER_TG_CHANNEL, execTime, msg)
}

func (c *addRehearsal) msgToChannel(r *store.Rehearsal, chUsername string) (string, error) {
	wd := lang.Weekdays[r.Time.Weekday().String()]
	m := lang.Months[r.Time.Month().String()]
	msg := fmt.Sprintf(lang.REHEARSAL_MSG_TO_CHANNEL, r.PlaceName, wd, r.Time.Day(), m, r.Time.Format(Conf.TimeLayout), r.Address, r.MapUrl)
	rh := &store.RemindChannel{
		Msg:             msg,
		ChannelUsername: chUsername,
	}
	details, err := json.Marshal(rh)
	return string(details), err
}

type nextRehearsal struct {
	u *store.User
}

func (c *nextRehearsal) IsEnd() bool {
	return true
}

func (c *nextRehearsal) IsAllow(u *store.User) bool {
	c.u = u
	return true
}

func (c *nextRehearsal) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	r, err := store.NextRehearsal()
	if err != nil {
		if err == store.ErrUndefinedNextRehearsal {
			replyMarkup.Text = lang.NEXT_REHEARSAL_UNDEFINED
			return replyMarkup, nil
		}
		return nil, err
	}
	replyMarkup.Text = fmt.Sprintf(lang.NEXT_REHEARSAL, r.PlaceName, r.Address, r.Time.Format(Conf.TimeLayout), r.MapUrl)
	return replyMarkup, nil
}

type deleteRehearsal struct {
	rehearsalID int
	u           *store.User
	stepConstructor
}

// newDeleteRehearsal creates deleteRehearsal and registers all steps
func newDeleteRehearsal() *deleteRehearsal {
	c := &deleteRehearsal{}
	c.RegisterSteps(c.firstStep, c.secondStep)
	return c
}

func (c *deleteRehearsal) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

// firstStep shows list of all rehearsals for further selecting and deleting
func (c *deleteRehearsal) firstStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	rehearsals, err := store.Rehearsals()
	if err != nil {
		return nil, err
	}
	for _, r := range rehearsals {
		lText := fmt.Sprintf(lang.REHEARSAL_ITEM, r.ID, r.Time, r.PlaceName)
		replyMarkup.Buttons = append(replyMarkup.Buttons, lText)
	}
	replyMarkup.Text = lang.REHEARSAL_CHOSE_REHEARSAL
	return replyMarkup, nil
}

// secondStep deletes rehearsal from db
func (c *deleteRehearsal) secondStep(answer string) (*ReplyMarkup, error) {
	rID, err := parseID(answer)
	if err != nil {
		return nil, err
	}
	err = store.DeleteRehearsal(rID)
	if err != nil {
		return nil, err
	}
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
		Text:    lang.DELETE_REHEARSAL_COMPLETE,
	}
	return replyMarkup, nil
}

// asSoonAsPossible returns nearest time  according to start/end hour
// i.e. we have start hour = 10, end hour = 21
// current time = 2019-01-01 12:00:00
// asSoonAsPossible returns 2019-01-01 12:00:00
// if current time = 2019-01-01 09:00:00
// asSoonAsPossible returns 2019-01-01 10:00:00
// if current time = 2019-01-01 22:00:00
// asSoonAsPossible returns 2019-01-02 10:00:00
func asSoonAsPossible() time.Time {
	loc, _ := time.LoadLocation(Conf.Location)
	curr := local_time.Now().In(loc)
	y, m, d := curr.Date()
	currH := curr.Hour()
	if currH < Conf.RemindHourStart {
		return time.Date(y, m, d, Conf.RemindHourStart, 0, 0, 0, loc)
	} else if currH > Conf.RemindHourEnd {
		// next day at as early as possible
		return time.Date(y, m, d, Conf.RemindHourStart, 0, 0, 0, loc).AddDate(0, 0, 1)
	}
	return curr
}
