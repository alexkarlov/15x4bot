package commands

import (
	"encoding/json"
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
	lectures    []int
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
		replyMarkup.Text = lang.WRONG_DATE_TIME
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
		replyMarkup.Text = lang.WRONG_DATE_TIME
		return replyMarkup, nil
	}
	pl, err := markupPlacesList()
	if err != nil {
		return nil, err
	}
	replyMarkup.Text = lang.PLACES_CHOSE_PLACE
	replyMarkup.Buttons = pl
	return replyMarkup, nil
}

func markupPlacesList() ([]string, error) {
	places, err := store.Places(store.PlaceTypes{store.PLACE_TYPE_FOR_EVENT, store.PLACE_TYPE_FOR_ALL})
	if err != nil {
		return nil, err
	}
	pl := make([]string, 0)
	for _, p := range places {
		b := fmt.Sprintf(lang.PLACES_LIST_BUTTONS, p.ID, p.Name)
		pl = append(pl, b)
	}
	return pl, nil
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
	lectures, err := store.Lectures(true)
	if err != nil {
		return nil, err
	}
	for _, l := range lectures {
		lText := fmt.Sprintf(lang.ADD_EVENT_LECTURES_LIST, l.ID, l.Name, l.Lector.Name)
		replyMarkup.Buttons = append(replyMarkup.Buttons, lText)
	}
	replyMarkup.Buttons = append(replyMarkup.Buttons, lang.ADD_EVENT_END_PHRASE)
	replyMarkup.Text = lang.ADD_EVENT_INTRO_LECTURES_LIST
	return replyMarkup, nil
}

// sixthStep saves the event in the db and sends final message
func (c *addEvent) sixthStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	if answer == lang.ADD_EVENT_END_PHRASE {
		_, err := store.AddEvent(c.whenStart, c.whenEnd, c.where, c.description, c.lectures)
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
	c.lectures = append(c.lectures, lID)
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
	return markupEventsList()
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

func markupEventsList() (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	events, err := store.Events()
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		eText := fmt.Sprintf(lang.DELETE_EVENT_ITEM, event.ID, event.StartTime.Format(Conf.TimeLayout))
		replyMarkup.Buttons = append(replyMarkup.Buttons, eText)
	}
	replyMarkup.Text = lang.EVENTS_CHOSE_EVENT
	return replyMarkup, nil
}

const (
	SEND_EVENT_TO_DESIGNERS typeSendEvent = 1
	SEND_EVENT_TO_CHAT      typeSendEvent = 2
	SEND_EVENT_TO_CHANNEL   typeSendEvent = 3
)

type typeSendEvent int

// sendEvent is a command for deleting events
type sendEvent struct {
	eventID int
	t       typeSendEvent
	u       *store.User
	stepConstructor
}

func (c *sendEvent) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

// newSendEvent creates sendEvent and registers all steps
func newSendEvent(t typeSendEvent) *sendEvent {
	c := &sendEvent{
		t: t,
	}
	c.RegisterSteps(c.firstStep, c.secondStep)
	return c
}

// firstStepDeleteEvent sends list of all events and asks a user to chose an event to send
func (c *sendEvent) firstStep(answer string) (*ReplyMarkup, error) {
	return markupEventsList()
}

// secondStepDeleteEvent send the event
func (c *sendEvent) secondStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	eID, err := parseID(answer)
	if err != nil {
		return nil, err
	}
	e, err := store.LoadEvent(eID)
	if err != nil {
		return nil, err
	}
	switch c.t {
	case SEND_EVENT_TO_DESIGNERS:
		for _, l := range e.Lectures {
			// check lectures descriptions
			if l.Description == "" {
				replyMarkup.Text = fmt.Sprintf(lang.ADD_EVENT_EMPTY_LECTURE_DESCRIPTION, l.ID, l.Name)
				return replyMarkup, nil
			}
			// check lectors profile pictures
			if l.Lector.PictureID == "" {
				replyMarkup.Text = fmt.Sprintf(lang.ADD_EVENT_EMPTY_PICTURE_SPEAKER, l.ID, l.Name)
				return replyMarkup, nil
			}
		}
		// create a task for send event to designers chat
		err = sendToDesignersChat(e)
	case SEND_EVENT_TO_CHANNEL:
		// check event FB field
		if e.FB == "" {
			replyMarkup.Text = lang.ADD_EVENT_EMPTY_FB
			return replyMarkup, nil
		}
		// check event poster
		if e.Poster == "" {
			replyMarkup.Text = lang.ADD_EVENT_EMPTY_POSTER
			return replyMarkup, nil
		}
		// create a task for send event to common channel
		err = sendToChannel(e)
	case SEND_EVENT_TO_CHAT:
		// check event FB field
		if e.FB == "" {
			replyMarkup.Text = lang.ADD_EVENT_EMPTY_FB
			return replyMarkup, nil
		}
		// create a task for send event to common chat
		err = sendToChat(e)
	}
	if err != nil {
		return nil, err
	}
	replyMarkup.Buttons = StandardMarkup(c.u.Role)
	replyMarkup.Text = lang.DONE
	return replyMarkup, nil
}

func sendToDesignersChat(e *store.Event) error {
	// if we don't have designers chat
	if Conf.DesignerChatID == "" {
		return nil
	}
	msg := fmt.Sprintf(lang.ADD_EVENT_SEND_EVENT_TO_DESIGNERS_CHAT, e.StartTime, e.EndTime, e.PlaceName)
	msgLectures := lang.LECTURE_LIST_ITEM
	rh := &store.RemindChannel{
		ChannelUsername: Conf.DesignerChatID,
	}
	for _, l := range e.Lectures {
		msg += "\n" + fmt.Sprintf(msgLectures, l.ID, l.Name, l.Description, l.Lector.Username, l.Lector.Name)
		rh.FileIDs = append(rh.FileIDs, l.Lector.PictureID)
	}
	rh.Msg = msg
	details, err := json.Marshal(rh)
	if err != nil {
		return err
	}
	execTime := asSoonAsPossible()
	return store.AddTask(store.TASK_TYPE_REMINDER_TG_CHANNEL, execTime, string(details))
}

func sendToChannel(e *store.Event) error {
	rh := &store.RemindChannel{
		Msg:             fmt.Sprintf(lang.ADD_EVENT_SEND_EVENT_TO_CHANNEL, e.StartTime.Format(Conf.TimeLayout), e.EndTime.Format("15:04"), e.PlaceName),
		ChannelUsername: Conf.MainChannelUsername,
		FileIDs:         []string{e.Poster},
	}
	details, err := json.Marshal(rh)
	if err != nil {
		return err
	}
	execTime := asSoonAsPossible()
	return store.AddTask(store.TASK_TYPE_REMINDER_TG_CHANNEL, execTime, string(details))
}

func sendToChat(e *store.Event) error {
	// if we don't have org chat - skip it
	if Conf.OrgChatID == "" {
		return nil
	}
	// load manual photo id
	a, err := store.LoadArticle("event_manual_photo_id")
	if err != nil {
		return err
	}
	rh := &store.RemindChannel{
		Msg:             fmt.Sprintf(lang.ADD_EVENT_SEND_EVENT_TO_COMMON_CHAT, e.StartTime.Format(Conf.TimeLayout), e.EndTime.Format("15:04"), e.PlaceName, e.FB),
		ChannelUsername: Conf.OrgChatID,
		FileIDs:         []string{a.Text},
	}
	details, err := json.Marshal(rh)
	if err != nil {
		return err
	}
	execTime := asSoonAsPossible()
	return store.AddTask(store.TASK_TYPE_REMINDER_TG_CHANNEL, execTime, string(details))
}

// updateEvent is a command for updating events
type updateEvent struct {
	field string
	e     *store.Event
	stepConstructor
}

// newUpdateEvent creates updateEvent and registers all steps
func newUpdateEvent() *updateEvent {
	c := &updateEvent{}
	c.RegisterSteps(c.firstStep, c.secondStep, c.thirdStep, c.fourthStep)
	return c
}

func (c *updateEvent) IsAllow(u *store.User) bool {
	return u.Role == store.USER_ROLE_ADMIN
}

func (c *updateEvent) firstStep(answer string) (*ReplyMarkup, error) {
	return markupEventsList()
}

func (c *updateEvent) secondStep(answer string) (*ReplyMarkup, error) {
	eID, err := parseID(answer)
	if err != nil {
		return nil, err
	}
	e, err := store.LoadEvent(eID)
	if err != nil {
		return nil, err
	}
	c.e = e
	replyMarkup := eventMarkupButtons()
	return replyMarkup, nil
}

func (c *updateEvent) thirdStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	switch answer {
	case lang.MARKUP_BUTTON_EVENT_CHANGE_START_DATE:
		c.field = "starttime"
		replyMarkup.Text = lang.CURRENT_VALUE + "\n" + c.e.StartTime.Format(Conf.TimeLayout) + "\n" + lang.EVENT_WHAT_IS_START_TIME
	case lang.MARKUP_BUTTON_EVENT_CHANGE_END_DATE:
		c.field = "endtime"
		replyMarkup.Text = lang.CURRENT_VALUE + "\n" + c.e.EndTime.Format(Conf.TimeLayout) + "\n" + lang.EVENT_WHAT_IS_END_TIME
	case lang.MARKUP_BUTTON_EVENT_CHANGE_DESCRIPTION:
		c.field = "description"
		replyMarkup.Text = lang.CURRENT_VALUE + "\n" + c.e.Description + "\n" + lang.EVENT_WHAT_IS_DESCRIPTION
	case lang.MARKUP_BUTTON_EVENT_CHANGE_PLACE:
		c.field = "place"
		replyMarkup.Text = lang.CURRENT_VALUE + "\n" + c.e.PlaceName + "\n" + lang.EVENT_WHAT_IS_PLACE
		pl, err := markupPlacesList()
		if err != nil {
			return nil, err
		}
		replyMarkup.Buttons = pl
	case lang.MARKUP_BUTTON_EVENT_CHANGE_FB:
		c.field = "fb"
		replyMarkup.Text = lang.CURRENT_VALUE + "\n" + c.e.FB + "\n" + lang.EVENT_WHAT_IS_FB
	case lang.MARKUP_BUTTON_EVENT_CHANGE_POSTER:
		c.field = "poster"
		replyMarkup.FileID = c.e.Poster
		replyMarkup.Text = lang.EVENT_WHAT_IS_POSTER
	}
	return replyMarkup, nil
}

func (c *updateEvent) fourthStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(store.USER_ROLE_ADMIN),
	}
	var err error
	switch c.field {
	case "starttime":
		sd, err := time.Parse(Conf.TimeLayout, answer)
		if err != nil {
			replyMarkup.Text = lang.WRONG_DATE_TIME
			c.RepeatStep()
			return replyMarkup, nil
		}
		replyMarkup.Text = lang.EVENT_START_TIME_SUCCESSFULY_UPDATED
		err = c.e.UpdateStartTime(sd)
	case "endtime":
		sd, err := time.Parse(Conf.TimeLayout, answer)
		if err != nil {
			replyMarkup.Text = lang.WRONG_DATE_TIME
			c.RepeatStep()
			return replyMarkup, nil
		}
		replyMarkup.Text = lang.EVENT_END_TIME_SUCCESSFULY_UPDATED
		err = c.e.UpdateEndTime(sd)
	case "description":
		replyMarkup.Text = lang.EVENT_DESCRIPTION_SUCCESSFULY_UPDATED
		err = c.e.UpdateDescription(answer)
	case "place":
		place, err := parseID(answer)
		if err != nil {
			return nil, err
		}
		ok, err := store.DoesPlaceExist(place)
		if err != nil {
			return nil, err
		}
		if !ok {
			replyMarkup.Text = lang.WRONG_PLACE_ID
			return replyMarkup, nil
		}
		replyMarkup.Text = lang.EVENT_PLACE_SUCCESSFULY_UPDATED
		err = c.e.UpdatePlace(place)
	case "fb":
		replyMarkup.Text = lang.EVENT_FB_SUCCESSFULY_UPDATED
		err = c.e.UpdateFB(answer)
	case "poster":
		replyMarkup.Text = lang.EVENT_POSTER_SUCCESSFULY_UPDATED
		err = c.e.UpdatePoster(answer)
	}
	return replyMarkup, err
}

func eventMarkupButtons() *ReplyMarkup {
	replyMarkup := &ReplyMarkup{
		Buttons: MessageButtons{lang.MARKUP_BUTTON_EVENT_CHANGE_START_DATE,
			lang.MARKUP_BUTTON_EVENT_CHANGE_END_DATE,
			lang.MARKUP_BUTTON_EVENT_CHANGE_DESCRIPTION,
			lang.MARKUP_BUTTON_EVENT_CHANGE_PLACE,
			lang.MARKUP_BUTTON_EVENT_CHANGE_FB,
			lang.MARKUP_BUTTON_EVENT_CHANGE_POSTER,
			lang.MARKUP_BUTTON_MAIN_MENU},
		Text: lang.CHOSE_MENU,
	}
	return replyMarkup
}
