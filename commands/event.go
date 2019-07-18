package commands

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/alexkarlov/15x4bot/store"
)

const END_PHRASE = "end"

type addEvent struct {
	step        int
	whenStart   time.Time
	whenEnd     time.Time
	where       int
	description string
	lections    []int
}

func (c *addEvent) IsAllow(u string) bool {
	//TODO: move it to db
	admins := []string{"zedman95", "alex_karlov"}
	for _, admin := range admins {
		if admin == u {
			return true
		}
	}
	return false
}

func (c *addEvent) NextStep(answer string) (string, error) {
	replyMsg := ""
	switch c.step {
	case 0:
		replyMsg = "Коли починається? Дата та час в форматі 2018-12-31 19:00:00"
	case 1:
		t, err := time.Parse("2006-01-02 15:04:05", answer)
		if err != nil {
			replyMsg = "Невірний формат дати та часу. Наприклад, якщо івент буде 20-ого грудня о 19:00 то треба ввести: 2018-12-20 19:00:00. Спробуй ще!"
			return replyMsg, nil
		}
		c.whenStart = t
		replyMsg = "Коли закінчується? Дата та час в форматі 2018-12-31 19:00:00"
	case 2:
		t, err := time.Parse("2006-01-02 15:04:05", answer)
		if err != nil {
			replyMsg = "Невірний формат дати та часу. Наприклад, якщо івент буде 20-ого грудня о 19:00 то треба ввести: 2018-12-20 19:00:00. Спробуй ще!"
			return replyMsg, nil
		}
		c.whenEnd = t
		places, err := store.GetPlaces(store.PlaceTypes{store.PLACE_TYPE_FOR_EVENT, store.PLACE_TYPE_FOR_ALL})
		if err != nil {
			return "", err
		}
		replyMsg = strings.Join([]string{"Де?", strings.Join(places, "\n")}, "\n")

	case 3:
		where, err := strconv.Atoi(answer)
		if err != nil {
			return "", errors.New("failed string to int converting")
		}
		c.where = where
		replyMsg = "Текст івенту"
	case 4:
		c.description = answer
		lections, err := store.GetLections(true)
		if err != nil {
			return "", err
		}
		replyMsg = strings.Join([]string{"Виберіть лекцію. Для закінчення напишіть \"" + END_PHRASE + "\"", strings.Join(lections, "\n")}, "\n")
	case 5:
		if answer == END_PHRASE {
			_, err := store.AddEvent(c.whenStart, c.whenEnd, c.where, c.description, c.lections)
			if err != nil {
				return "", err
			}
			replyMsg = "Івент створено"
			break
		}
		// TODO: process answer as lections in one answer
		lection, err := strconv.Atoi(answer)
		if err != nil {
			return "", err
		}
		c.lections = append(c.lections, lection)
		// desrese step counter for returning on the next iteration to the same step
		c.step--
	default:
		return "", errors.New("next step for command addRepetition was called in a wrong way")
	}
	c.step++
	return replyMsg, nil
}

func (c *addEvent) IsEnd() bool {
	return c.step == 6
}

type nextEvent struct {
}

func (c *nextEvent) IsEnd() bool {
	return true
}

func (c *nextEvent) IsAllow(u string) bool {
	return true
}

func (c *nextEvent) NextStep(answer string) (string, error) {
	e, err := store.GetNextEvent()
	if err != nil {
		if err == store.ErrUndefinedNextEvent {
			return "Невідомо коли, запитайся пізніше", nil
		}
		return "", err
	}
	replyMsg := strings.Join([]string{"Де: ", e.PlaceName, ", ", e.Address, "\n", "Початок: ", e.StartTime.Format("2006-01-02 15:04:05"), "\n", "Кінець: ", e.EndTime.Format("2006-01-02 15:04:05")}, "")
	return replyMsg, nil
}
