package commands

import (
	"fmt"
	"github.com/alexkarlov/15x4bot/lang"
	"time"

	"github.com/alexkarlov/15x4bot/store"
)

type addEvent struct {
	step        int
	whenStart   time.Time
	whenEnd     time.Time
	where       int
	description string
	lections    []int
	u           *store.User
}

func (c *addEvent) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

func (c *addEvent) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	var err error
	switch c.step {
	case 0:
		replyMarkup.Text = lang.ADD_EVENT_WHEN_START
	case 1:
		t, err := time.Parse(Conf.TimeLayout, answer)
		if err != nil {
			replyMarkup.Text = lang.ADD_EVENT_WRONG_DATE
			return replyMarkup, nil
		}
		c.whenStart = t
		replyMarkup.Text = lang.ADD_EVENT_WHEN_END
	case 2:
		c.whenEnd, err = time.Parse(Conf.TimeLayout, answer)
		if err != nil {
			replyMarkup.Text = lang.ADD_EVENT_WRONG_DATE
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
	case 3:
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
		replyMarkup.Text = "Текст івенту"
	case 4:
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
	case 5:
		if answer == lang.ADD_EVENT_END_PHRASE {
			_, err := store.AddEvent(c.whenStart, c.whenEnd, c.where, c.description, c.lections)
			if err != nil {
				return nil, err
			}
			replyMarkup.Text = "Івент створено"
			replyMarkup.Buttons = StandardMarkup(c.u.Role)
			break
		}
		lID, err := parseID(answer)
		if err != nil {
			return nil, err
		}
		c.lections = append(c.lections, lID)
		// desrese step counter for returning on the next iteration to the same step
		c.step--
	}
	c.step++
	return replyMarkup, nil
}

func (c *addEvent) IsEnd() bool {
	return c.step == 6
}

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

type deleteEvent struct {
	step    int
	eventID int
	u       *store.User
}

func (c *deleteEvent) IsEnd() bool {
	return c.step == 2
}

func (c *deleteEvent) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

func (c *deleteEvent) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	switch c.step {
	case 0:
		events, err := store.Events()
		if err != nil {
			return nil, err
		}
		for _, event := range events {
			eText := fmt.Sprintf(lang.DELETE_EVENT_ITEM, event.ID, event.StartTime)
			replyMarkup.Buttons = append(replyMarkup.Buttons, eText)
		}
		replyMarkup.Text = lang.EVENTS_CHOSE_EVENT
	case 1:
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
	}
	c.step++
	return replyMarkup, nil
}
