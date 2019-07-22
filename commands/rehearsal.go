package commands

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/alexkarlov/15x4bot/store"
)

var (
	ErrWrongPlace = errors.New("wrong place id: failed to conver from string to int")
)

const (
	TEMPLATE_NEXT_EVENT = "Де: %s, %s\nКоли: %s\nМапа:%s"
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

func (c *addRehearsal) NextStep(answer string) (replyMsg string, err error) {
	switch c.step {
	case 0:
		replyMsg = "Коли? Дата та час в форматі 2018-12-31 19:00:00"
	case 1:
		t, err := time.Parse("2006-01-02 15:04:05", answer)
		if err != nil {
			replyMsg = "Невірний формат дати та часу. Наприклад, якщо репетиція буде 20-ого грудня о 19:00 то треба ввести: 2018-12-20 19:00:00. Спробуй ще!"
			return replyMsg, nil
		}
		c.when = t
		places, err := store.Places(store.PlaceTypes{store.PLACE_TYPE_FOR_REHEARSALS, store.PLACE_TYPE_FOR_ALL})
		pText := ""
		for _, p := range places {
			pText = fmt.Sprintf(TEMPLATE_PLACES_LIST, p.ID, p.Name, p.Address)
		}
		replyMsg = fmt.Sprintf(TEMPLATE_INTRO_PLACES_LIST, pText)
	case 2:
		c.where, err = strconv.Atoi(answer)
		if err != nil {
			return "", ErrWrongPlace
		}
		err = store.AddRehearsal(c.when, c.where)
		if err != nil {
			return "", err
		}
		replyMsg = "Репетиція створена"
	}
	c.step++
	return
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

func (c *nextRep) NextStep(answer string) (string, error) {
	r, err := store.NextRehearsal()
	if err != nil {
		if err == store.ErrUndefinedNextRehearsal {
			return "Невідомо коли, запитайся пізніше", nil
		}
		return "", err
	}
	replyMsg := fmt.Sprintf(TEMPLATE_NEXT_EVENT, r.PlaceName, r.Address, r.Time.Format("2006-01-02 15:04:05"), r.MapUrl)
	return replyMsg, nil
}
