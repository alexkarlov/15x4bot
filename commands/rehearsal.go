package commands

import (
	"encoding/json"
	"fmt"
	"github.com/alexkarlov/15x4bot/lang"
	"time"

	"github.com/alexkarlov/15x4bot/store"
)

type addRehearsal struct {
	step  int
	when  time.Time
	where int
	u     *store.User
}

func (c *addRehearsal) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

func (c *addRehearsal) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	var err error
	switch c.step {
	case 0:
		replyMarkup.Text = lang.ADD_REHEARSAL_WHEN
	case 1:
		t, err := time.Parse(Conf.TimeLayout, answer)
		if err != nil {
			replyMarkup.Text = lang.ADD_REHEARSAL_ERROR_DATE
			return replyMarkup, nil
		}
		c.when = t
		places, err := store.Places(store.PlaceTypes{store.PLACE_TYPE_FOR_REHEARSAL, store.PLACE_TYPE_FOR_ALL})
		replyMarkup.Buttons = nil
		for _, p := range places {
			b := fmt.Sprintf(lang.PLACES_LIST_BUTTONS, p.ID, p.Name)
			replyMarkup.Buttons = append(replyMarkup.Buttons, b)
		}
		replyMarkup.Text = lang.PLACES_CHOSE_PLACE
	case 2:
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
		id, err := store.AddRehearsal(c.when, c.where)
		// create a task for sending post in the internal channel
		if err != nil {
			return nil, err
		}
		err = addRehearsalReminder(id)
		if err != nil {
			replyMarkup.Text = lang.ADD_REHEARSAL_ERROR_REMINDER_MSG
			break
		}
		replyMarkup.Text = lang.ADD_REHEARSAL_SUCCESS_MSG
		replyMarkup.Buttons = StandardMarkup(c.u.Role)
	}
	c.step++
	return replyMarkup, nil
}

func addRehearsalReminder(ID int) error {
	// create chat reminder
	r, err := store.LoadRehearsal(ID)
	if err != nil {
		return err
	}
	wd := lang.Weekdays[r.Time.Weekday().String()]
	m := lang.Months[r.Time.Month().String()]
	msg := fmt.Sprintf(lang.REHEARSAL_MSG_TO_CHANNEL, r.PlaceName, wd, r.Time.Day(), m, r.Time.Format(Conf.TimeLayout), r.Address, r.MapUrl)
	rh := &store.RemindChannel{
		Msg:             msg,
		ChannelUsername: Conf.InternalChannelUsername,
	}
	details, err := json.Marshal(rh)
	if err != nil {
		return err
	}
	execTime := asSoonAsPossible()
	err = store.AddTask(store.TASK_TYPE_REMINDER_TG_CHANNEL, execTime, string(details))
	if err != nil {
		return err
	}
	// create channel reminder
	rh.ChannelUsername = Conf.InternalChannelUsername
	details, err = json.Marshal(rh)
	if err != nil {
		return err
	}
	return store.AddTask(store.TASK_TYPE_REMINDER_TG_CHANNEL, execTime, string(details))
}

func (c *addRehearsal) IsEnd() bool {
	return c.step == 3
}

type nextRep struct {
	u *store.User
}

func (c *nextRep) IsEnd() bool {
	return true
}

func (c *nextRep) IsAllow(u *store.User) bool {
	c.u = u
	return true
}

func (c *nextRep) NextStep(answer string) (*ReplyMarkup, error) {
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
	replyMarkup.Text = fmt.Sprintf(lang.NEXT_REHEARSAL, r.PlaceName, r.Address, r.Time.Format("2006-01-02 15:04:05"), r.MapUrl)
	return replyMarkup, nil
}

type deleteRehearsal struct {
	step        int
	rehearsalID int
	u           *store.User
}

func (c *deleteRehearsal) IsEnd() bool {
	return c.step == 2
}

func (c *deleteRehearsal) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN
}

func (c *deleteRehearsal) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	switch c.step {
	case 0:
		rehearsals, err := store.Rehearsals()
		if err != nil {
			return nil, err
		}
		for _, r := range rehearsals {
			lText := fmt.Sprintf(lang.REHEARSAL_ITEM, r.ID, r.Time, r.PlaceName)
			replyMarkup.Buttons = append(replyMarkup.Buttons, lText)
		}
		replyMarkup.Text = lang.REHEARSAL_CHOSE_REHEARSAL
	case 1:
		rID, err := parseID(answer)
		if err != nil {
			return nil, err
		}
		err = store.DeleteRehearsal(rID)
		if err != nil {
			return nil, err
		}
		replyMarkup.Buttons = StandardMarkup(c.u.Role)
		replyMarkup.Text = lang.DELETE_REHEARSAL_COMPLETE
	}
	c.step++
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
	curr := time.Now().In(loc)
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
