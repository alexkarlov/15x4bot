package commands

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/alexkarlov/15x4bot/store"
	"github.com/alexkarlov/15x4bot/utils"
	"gopkg.in/telegram-bot-api.v4"
)

type addRepetition struct {
	step  int
	when  time.Time
	where int
}

func (c *addRepetition) IsAllow(u *tgbotapi.User) bool {
	//TODO: move it to db
	t := []string{"zedman95", "alex_karlov"}
	return utils.Contains(t, u.UserName)
}

func (c *addRepetition) NextStep(answer string) (replyMsg string, err error) {
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
			err = errors.New("Failed string to int converting")
			return
		}
		store.AddRepetition(c.when, c.where)
		replyMsg = "Репетиція створена"
	default:
		err = errors.New("Next step for command addRepetition was called in a wrong way")
	}
	c.step++
	return
}

func (c *addRepetition) IsEnd() bool {
	return c.step == 3
}

type nextRep struct {
}

func (c *nextRep) IsEnd() bool {
	return true
}

func (c *nextRep) IsAllow(u *tgbotapi.User) bool {
	return true
}

func (c *nextRep) NextStep(answer string) (replyMsg string, err error) {
	r, err := store.GetNextRepetition()
	if err != nil {
		return "", err
	}
	replyMsg = strings.Join([]string{"Де: ", r.PlaceName, ", ", r.Address, "\n", "Коли: ", r.Time.Format("2006-01-02 15:04:05"), "\n", "Мапа: ", r.MapUrl}, "")

	return
}
