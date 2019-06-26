package commands

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/alexkarlov/15x4bot/store"
	"gopkg.in/telegram-bot-api.v4"
)

type addEvent struct {
	step  int
	when  time.Time
	where int
}

func (c *addEvent) IsAllow(u *tgbotapi.User) bool {
	return true
}

func (c *addEvent) NextStep(answer string) (replyMsg string, err error) {
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
		places, err := store.GetPlaces()
		if err != nil {
			return "", err
		}
		replyMsg = strings.Join([]string{"Де?", strings.Join(places, "\n")}, "\n")

	case 2:
		c.where, err = strconv.Atoi(answer)
		if err != nil {
			err = errors.New("failed string to int converting")
			return
		}
		store.AddRepetition(c.when, c.where)
		replyMsg = "Репетиція створена"
	default:
		err = errors.New("next step for command addRepetition was called in a wrong way")
	}
	c.step++
	return
}

type nextEvent struct {
}

func (c *nextEvent) IsEnd() bool {
	return true
}

func (c *nextEvent) IsAllow(u string) bool {
	return true
}

func (c *nextEvent) NextStep(answer string) (replyMsg string, err error) {
	replyMsg = "Поки невідомо, запитай пізніше"

	return
}
