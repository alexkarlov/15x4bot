package commands

import (
	"encoding/json"
	"fmt"
	"github.com/alexkarlov/15x4bot/lang"
	"github.com/alexkarlov/15x4bot/store"
	"strconv"
	"strings"
	"time"
)

const (
	UserRemindHour = 19
)

func nextDay(hour int) time.Time {
	curr := time.Now()
	y, m, d := curr.Date()
	loc, _ := time.LoadLocation(Conf.Location)
	rTime := time.Date(y, m, d, hour, 0, 0, 0, loc).AddDate(0, 0, 1)
	return rTime
}

type upsertLection struct {
	exists      bool
	ID          int
	u           *store.User
	name        string
	description string
	userID      int
	stepConstructor
}

func (c *upsertLection) IsAllow(u *store.User) bool {
	c.u = u
	return (u.Role == store.USER_ROLE_ADMIN || u.Role == store.USER_ROLE_LECTOR)
}

// newUpsertLection creates upsertLection and registers all steps
// it receives argument whether lection exists or no
func newUpsertLection(e bool) *upsertLection {
	c := &upsertLection{
		exists: e,
	}
	c.RegisterSteps(c.firstStep, c.secondStep, c.thirdStep, c.fourthStep, c.fifthStep)
	return c
}

// firstStep depends on the type of operation
// if we update existed lectures, at this step this command shows all available lectures
// if we create a new lecture and the current user is admin
// this command ask the name of speaker
// if we create a new lecture and the current user is NOT admin
// this command saves id (as owner of lecture) of the current user
// and asks lecture name
func (c *upsertLection) firstStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	if c.exists {
		return lections(c.u, false, false)
	}
	// if the user is a lector = add him as an lection owner and skip step
	if c.u.Role == store.USER_ROLE_LECTOR {
		c.userID = c.u.ID
		return c.SkipStep(answer)
	}
	users, err := store.Users([]store.UserRole{store.USER_ROLE_ADMIN, store.USER_ROLE_LECTOR})
	if err != nil {
		return nil, err
	}
	speakerText := ""
	for _, u := range users {
		speakerText += fmt.Sprintf(lang.UPSERT_LECTURE_STEP_SPEAKER_DETAILS, u.ID, u.Username, u.Name)
	}
	replyMarkup.Text = fmt.Sprintf(lang.UPSERT_LECTURE_STEP_SPEAKER, speakerText)
	return replyMarkup, nil
}

// secondStep saves the speaker id (if the current user is admin and we try to add new lecture)
// if the current user is NOT admin we skip this step
func (c *upsertLection) secondStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	replyMarkup.Text = lang.UPSERT_LECTURE_STEP_LECTURE_NAME
	if c.exists {
		lID, err := parseID(answer)
		if err != nil {
			return nil, err
		}
		c.ID = lID
		return replyMarkup, nil
	}
	userID, err := strconv.Atoi(answer)
	if err != nil {
		return nil, ErrWrongID
	}
	ok, err := store.DoesUserExist(userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		replyMarkup.Text = lang.WRONG_USER_ID
		return replyMarkup, nil
	}
	c.userID = userID
	return replyMarkup, nil
}

// thirdStep saves lection name and asks description of lecture
func (c *upsertLection) thirdStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	c.name = answer
	replyMarkup.Text = lang.UPSERT_LECTURE_STEP_LECTURE_DESCRIPTION
	replyMarkup.Buttons = append(replyMarkup.Buttons, lang.I_DONT_KNOW)
	return replyMarkup, nil
}

// fourthStep saves lecture into db and (if needed) creates a new task reminder
func (c *upsertLection) fourthStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	if answer != lang.I_DONT_KNOW {
		c.description = answer
	}
	var lectionID int
	var err error
	if c.exists {
		lectionID = c.ID
		err = store.UpdateLection(c.ID, c.name, c.description)
	} else {
		lectionID, err = store.AddLection(c.name, c.description, c.userID)
		c.ID = lectionID
	}
	if err != nil {
		return nil, err
	}
	// if user didn't provide lecture description
	// create a new reminder task
	if c.description == "" {
		execTime := nextDay(UserRemindHour)
		r := &store.RemindLection{
			ID: lectionID,
		}
		details, err := json.Marshal(r)
		if err != nil {
			return nil, err
		}
		store.AddTask(store.TASK_TYPE_REMINDER_LECTOR, execTime, string(details))
		replyMarkup.Text += "\n" + fmt.Sprintf(lang.UPSERT_LECTURE_I_WILL_REMIND, execTime.Format(Conf.TimeLayout))
	} else {
		replyMarkup.Text += "\n\n" + lang.UPSERT_LECTURE_SEND_TO_GRAMMAR_NAZI
		replyMarkup.Buttons = MessageButtons{lang.MARKUP_BUTTON_NO, lang.MARKUP_BUTTON_YES}
	}
	return replyMarkup, nil
}

// fourthStep saves lecture into db and (if needed) creates a new task reminder
func (c *upsertLection) fifthStep(answer string) (*ReplyMarkup, error) {
	replyMarkup, err := stepSendToGrammarChat(c.ID, answer, c.u)
	if c.exists {
		replyMarkup.Text = lang.UPSERT_LECTURE_SUCCESS_UPDATE_MSG
	} else {
		replyMarkup.Text = lang.UPSERT_LECTURE_SUCCESS_CREATE_MSG
	}
	return replyMarkup, err
}

type addDescriptionLection struct {
	u         *store.User
	lectionID int
	stepConstructor
}

// newAddDescriptionLection creates addDescriptionLection and registers all steps
func newAddDescriptionLection() *addDescriptionLection {
	c := &addDescriptionLection{}
	c.RegisterSteps(c.firstStep, c.secondStep, c.thirdStep, c.fourthStep)
	return c
}

func (c *addDescriptionLection) firstStep(answer string) (*ReplyMarkup, error) {
	return lections(c.u, true, true)
}

// secondStep asks lecture description
func (c *addDescriptionLection) secondStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	lID, err := parseID(answer)
	if err != nil {
		return nil, err
	}
	l, err := store.LoadLection(lID)
	if err != nil {
		return nil, err
	}
	if c.u.Role != store.USER_ROLE_ADMIN && l.Lector.Username != c.u.Username {
		replyMarkup.Text = lang.LECTURES_ERROR_NOT_YOUR
		return replyMarkup, nil
	}
	c.lectionID = lID
	replyMarkup.Text = lang.UPSERT_LECTURE_STEP_LECTURE_DESCRIPTION
	return replyMarkup, nil
}

// thirdStep saves description in the db and sends it to the grammar nazi chat
func (c *addDescriptionLection) thirdStep(answer string) (*ReplyMarkup, error) {
	err := store.AddLectionDescription(c.lectionID, answer)
	if err != nil {
		return nil, err
	}
	replyMarkup := &ReplyMarkup{
		Buttons: MessageButtons{lang.MARKUP_BUTTON_YES, lang.MARKUP_BUTTON_NO},
		Text:    lang.UPSERT_LECTURE_SEND_TO_GRAMMAR_NAZI,
	}
	return replyMarkup, nil
}

func (c *addDescriptionLection) fourthStep(answer string) (*ReplyMarkup, error) {
	replyMarkup, err := stepSendToGrammarChat(c.lectionID, answer, c.u)
	replyMarkup.Text = lang.ADD_LECTURE_DESCRIPTION_COMPLETE
	return replyMarkup, err
}

func stepSendToGrammarChat(lID int, answer string, u *store.User) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(u.Role),
	}
	if answer == lang.MARKUP_BUTTON_NO {
		return replyMarkup, nil
	}
	// send to the grammar-nazi chat
	err := sendTextToGrammarNazi(lID)
	if err != nil {
		return nil, err
	}
	replyMarkup.Text = lang.ADD_LECTURE_DESCRIPTION_REMINDER_MSG_OK
	return replyMarkup, nil
}

func sendTextToGrammarNazi(ID int) error {
	l, err := store.LoadLection(ID)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf(lang.ADD_LECTURE_DESCRIPTION_MSG_TO_GRAMMAR_NAZI, l.Name, l.Description)
	rh := &store.RemindChannel{
		Msg:             msg,
		ChannelUsername: Conf.GrammarNaziChatID,
	}
	details, err := json.Marshal(rh)
	if err != nil {
		return err
	}
	execTime := asSoonAsPossible()
	return store.AddTask(store.TASK_TYPE_REMINDER_TG_CHANNEL, execTime, string(details))
}

// lections adds buttons for further manipulation with particular lecture
func lections(u *store.User, onlyNew bool, withoutDescription bool) (*ReplyMarkup, error) {
	reply := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	lections, err := store.Lections(onlyNew)
	if err != nil {
		return nil, err
	}
	var l []string
	for _, lection := range lections {
		if (u.Role != store.USER_ROLE_ADMIN && u.ID != lection.Lector.ID) || (withoutDescription && lection.Description != "") {
			// skip lections which doesn't belong to user (if he isn't admin) or it has description
			continue
		}
		l = append(l, fmt.Sprintf(lang.UPSERT_LECTURE_ITEM, lection.ID, lection.Name))
	}
	// if there are no appropriate lections - send special response
	if len(l) == 0 {
		reply.Text = lang.LECTURE_LIST_EMPTY
		return reply, nil
	}
	reply.Buttons = append(reply.Buttons, l...)
	reply.Text = lang.ADD_LECTURE_DESCIRPTION_CHOSE_LECTURE
	return reply, nil
}

func (c *addDescriptionLection) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN || u.Role == store.USER_ROLE_LECTOR
}

type lectionsList struct {
	u                  *store.User
	withoutDescription bool
}

func (c *lectionsList) IsEnd() bool {
	return true
}

func (c *lectionsList) IsAllow(u *store.User) bool {
	c.u = u
	return (u.Role == store.USER_ROLE_ADMIN || u.Role == store.USER_ROLE_LECTOR)
}

func (c *lectionsList) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	list, err := store.Lections(true)
	if err != nil {
		return nil, err
	}
	var l []string
	for _, lection := range list {
		if c.withoutDescription && lection.Description != "" {
			// if we want to see only lections without descriptions and the current lection does have description - skip it
			continue
		}
		if c.u.Role != store.USER_ROLE_ADMIN && c.u.ID != lection.Lector.ID {
			// skip lections which doesn't belong to user (if he isn't admin)
			continue
		}
		l = append(l, fmt.Sprintf(lang.LECTURE_LIST_ITEM, lection.ID, lection.Name, lection.Description, lection.Lector.Username, lection.Lector.Name))
	}
	// if there are no appropriate lections - send special response
	if len(l) == 0 {
		replyMarkup.Text = lang.LECTURE_LIST_EMPTY
		return replyMarkup, nil
	}
	replyMarkup.Text = strings.Join(l, "\n\n")
	return replyMarkup, nil
}

type deleteLection struct {
	lectionID int
	u         *store.User
	stepConstructor
}

// newDeleteLection creates deleteLection and registers all steps
func newDeleteLection() *deleteLection {
	c := &deleteLection{}
	c.RegisterSteps(c.firstStep, c.secondStep)
	return c
}

func (c *deleteLection) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN || u.Role == store.USER_ROLE_LECTOR
}

func (c *deleteLection) firstStep(answer string) (*ReplyMarkup, error) {
	return lections(c.u, false, false)
}

// secondStep deletes particular lecture
func (c *deleteLection) secondStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	lID, err := parseID(answer)
	if err != nil {
		return nil, err
	}
	l, err := store.LoadLection(lID)
	if err != nil {
		return nil, err
	}
	if c.u.Role != store.USER_ROLE_ADMIN && c.u.ID != l.Lector.ID {
		replyMarkup.Text = lang.LECTURES_ERROR_NOT_YOUR
		return replyMarkup, nil
	}
	err = store.DeleteLection(lID)
	if err != nil {
		return nil, err
	}
	replyMarkup.Buttons = StandardMarkup(c.u.Role)
	replyMarkup.Text = lang.DELETE_LECTURE_COMPLETE
	return replyMarkup, nil
}
