package commands

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/alexkarlov/15x4bot/store"
)

var (
	ErrWrongPlace = errors.New("wrong place id: failed to conver from string to int")
)

const (
	TEMPLATE_NEXT_REHEARSAL           = "Де: %s, %s\nКоли: %s\nМапа:%s"
	TEMPLATE_ADD_REHEARSAL_WHEN       = "Коли? Дата та час в форматі 2018-12-31 19:00:00"
	TEMPLATE_ADD_REHEARSAL_ERROR_DATE = "Невірний формат дати та часу. Наприклад, якщо репетиція буде 20-ого грудня о 19:00 то треба ввести: 2018-12-20 19:00:00. Спробуй ще!"
	TEMPLATE_ADDREHEARSAL_SUCCESS_MSG = "Репетиція створена"
	TEMPLATE_NEXT_REHEARSAL_UNDEFINED = "Невідомо коли, запитайся пізніше"

	TEMPLATE_ADD_REHEARSAL_ERROR_WRONG_PLACE = "Неправильне місце"
)

type addRehearsal struct {
	step  int
	when  time.Time
	where int
}

func (c *addRehearsal) IsAllow(u string) bool {
	//TODO: move it to db
	admins := []string{"zedman95", "alex_karlov"}
	for _, admin := range admins {
		if admin == u {
			return true
		}
	}
	return false
}

func (c *addRehearsal) NextStep(u *store.User, answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	var err error
	switch c.step {
	case 0:
		replyMarkup.Text = TEMPLATE_ADD_REHEARSAL_WHEN
	case 1:
		t, err := time.Parse("2006-01-02 15:04:05", answer)
		if err != nil {
			replyMarkup.Text = TEMPLATE_ADD_REHEARSAL_ERROR_DATE
			return replyMarkup, nil
		}
		c.when = t
		places, err := store.Places(store.PlaceTypes{store.PLACE_TYPE_FOR_REHEARSAL, store.PLACE_TYPE_FOR_ALL})
		pText := ""
		replyMarkup.Buttons = nil
		for _, p := range places {
			pText += fmt.Sprintf(TEMPLATE_PLACES_LIST, p.ID, p.Name, p.Address)
			b := fmt.Sprintf(TEMPLATE_PLACES_LIST_BUTTONS, p.ID, p.Name)
			replyMarkup.Buttons = append(replyMarkup.Buttons, b)
		}
		replyMarkup.Text = fmt.Sprintf(TEMPLATE_INTRO_PLACES_LIST, pText)
	case 2:
		regexpLectionID := regexp.MustCompile(`^(\d+)?\.`)
		matches := regexpLectionID.FindStringSubmatch(answer)
		if len(matches) < 2 {
			return nil, ErrWrongPlace
		}
		c.where, err = strconv.Atoi(matches[1])
		if err != nil {
			return nil, ErrWrongPlace
		}
		err = store.AddRehearsal(c.when, c.where)
		if err != nil {
			return nil, err
		}
		replyMarkup.Text = TEMPLATE_ADDREHEARSAL_SUCCESS_MSG
		replyMarkup.Buttons = StandardMarkup(u.Role)
	}
	c.step++
	return replyMarkup, nil
}

func (c *addRehearsal) IsEnd() bool {
	return c.step == 3
}

type nextRep struct {
}

func (c *nextRep) IsEnd() bool {
	return true
}

func (c *nextRep) IsAllow(u string) bool {
	return true
}

func (c *nextRep) NextStep(u *store.User, answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(u.Role),
	}
	r, err := store.NextRehearsal()
	if err != nil {
		if err == store.ErrUndefinedNextRehearsal {
			replyMarkup.Text = TEMPLATE_NEXT_REHEARSAL_UNDEFINED
			return replyMarkup, nil
		}
		return nil, err
	}
	replyMarkup.Text = fmt.Sprintf(TEMPLATE_NEXT_REHEARSAL, r.PlaceName, r.Address, r.Time.Format("2006-01-02 15:04:05"), r.MapUrl)
	return replyMarkup, nil
}
