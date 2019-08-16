package commands

import (
	"fmt"
	"github.com/alexkarlov/15x4bot/lang"
	"time"

	"github.com/alexkarlov/15x4bot/store"
)

type addEvent struct {
	whenStart   time.Time
	whenEnd     time.Time
	where       int
	description string
	lections    []int
	u           *store.User
	stepConstructor
}

// newAddEvent creates addEvent and registers all steps
func newAddEvent() *addEvent {
	c := &addEvent{}
	c.RegisterSteps(c.firstStep, c.secondStep, c.thirdStep, c.fourthStep, c.fifthStep, c.sixthStep)
	return c
}

func (c *addEvent) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

// firstStep saves start date and asks end date
func (c *addEvent) firstStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
		Text:    lang.ADD_EVENT_WHEN_START,
	}
	return replyMarkup, nil
}

// secondStep saves start date and asks end date
func (c *addEvent) secondStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	t, err := time.Parse(Conf.TimeLayout, answer)
	if err != nil {
		replyMarkup.Text = lang.ADD_EVENT_WRONG_DATETIME
		return replyMarkup, nil
	}
	c.whenStart = t
	replyMarkup.Text = lang.ADD_EVENT_WHEN_END
	return replyMarkup, nil
}

// thirdStep saves end date and asks place for the event
func (c *addEvent) thirdStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	var err error
	c.whenEnd, err = time.Parse(Conf.TimeLayout, answer)
	if err != nil {
		replyMarkup.Text = lang.ADD_EVENT_WRONG_DATETIME
		return replyMarkup, nil
	}
	places, err := store.Places(store.PlaceTypes{store.PLACE_TYPE_FOR_EVENT, store.PLACE_TYPE_FOR_ALL})
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

// fourthStep saves place for the event and asks for text of the event
func (c *addEvent) fourthStep(answer string) (*ReplyMarkup, error) {
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
	replyMarkup.Text = lang.ADD_EVENT_TEXT_EVENT
	return replyMarkup, nil
}

// fifthStep saves text of the event and asks for lectures of the event
func (c *addEvent) fifthStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	var err error
	c.description = answer
	lections, err := store.Lections(true)
	if err != nil {
		return nil, err
	}
	for _, l := range lections {
		lText := fmt.Sprintf(lang.ADD_EVENT_LECTIONS_LIST, l.ID, l.Name, l.Lector.Name)
		replyMarkup.Buttons = append(replyMarkup.Buttons, lText)
	}
	replyMarkup.Buttons = append(replyMarkup.Buttons, lang.ADD_EVENT_END_PHRASE)
	replyMarkup.Text = lang.ADD_EVENT_INTRO_LECTIONS_LIST
	return replyMarkup, nil
}

// sixthStep saves the event in the db and sends final message
func (c *addEvent) sixthStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	if answer == lang.ADD_EVENT_END_PHRASE {
		_, err := store.AddEvent(c.whenStart, c.whenEnd, c.where, c.description, c.lections)
		if err != nil {
			return nil, err
		}
		replyMarkup.Text = lang.ADD_EVENT_SUCCESS_MSG
		replyMarkup.Buttons = StandardMarkup(c.u.Role)
		return replyMarkup, nil
	}
	lID, err := parseID(answer)
	if err != nil {
		return nil, err
	}
	c.lections = append(c.lections, lID)
	// desrese step counter for returning on the next iteration to the same step
	c.RepeatStep()
	return replyMarkup, nil
}

// nextEvent is a command which just selects next event and sends it (or sends text that next event is undefined)
type nextEvent struct {
	u *store.User
}

func (c *nextEvent) IsEnd() bool {
	return true
}

func (c *nextEvent) IsAllow(u *store.User) bool {
	c.u = u
	return true
}

func (c *nextEvent) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	e, err := store.NextEvent()
	if err != nil {
		if err == store.ErrUndefinedNextEvent {
			replyMarkup.Text = lang.NEXT_EVENT_UNDEFINED
			return replyMarkup, nil
		}
		return nil, err
	}
	replyMarkup.Text = fmt.Sprintf(lang.NEXT_EVENT, e.PlaceName, e.Address, e.StartTime.Format(Conf.TimeLayout), e.EndTime.Format(Conf.TimeLayout))
	return replyMarkup, nil
}

// eventsList is a command for sends list of all events
type eventsList struct {
	u *store.User
}

func (c *eventsList) IsEnd() bool {
	return true
}

func (c *eventsList) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

func (c *eventsList) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	list, err := store.Events()
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		replyMarkup.Text = lang.EVENTS_LIST_EMPTY
		return replyMarkup, nil
	}
	for _, event := range list {
		replyMarkup.Text += fmt.Sprintf(lang.EVENTS_LIST_ITEM, event.ID, event.StartTime, event.EndTime, event.PlaceName, event.Address)
	}
	return replyMarkup, nil
}

// deleteEvent is a command for deleting events
type deleteEvent struct {
	eventID int
	u       *store.User
	stepConstructor
}

func (c *deleteEvent) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

// newDeleteEvent creates deleteEvent and registers all steps
func newDeleteEvent() *deleteEvent {
	c := &deleteEvent{}
	c.RegisterSteps(c.firstStep, c.secondStep)
	return c
}

// firstStepDeleteEvent sends list of all events and asks a user to chose an event to delete
func (c *deleteEvent) firstStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	events, err := store.Events()
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		eText := fmt.Sprintf(lang.DELETE_EVENT_ITEM, event.ID, event.StartTime)
		replyMarkup.Buttons = append(replyMarkup.Buttons, eText)
	}
	replyMarkup.Text = lang.EVENTS_CHOSE_EVENT
	return replyMarkup, nil
}

// secondStepDeleteEvent deletes the event by user's answers
func (c *deleteEvent) secondStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	eID, err := parseID(answer)
	if err != nil {
		return nil, err
	}
	err = store.DeleteEvent(eID)
	if err != nil {
		return nil, err
	}
	replyMarkup.Buttons = StandardMarkup(c.u.Role)
	replyMarkup.Text = lang.DELETE_EVENT_COMPLETE
	return replyMarkup, nil
}
