package commands

import (
	"errors"
	"fmt"
	"github.com/alexkarlov/simplelog"
	"regexp"
	"strconv"
	"time"

	"github.com/alexkarlov/15x4bot/store"
)

const (
	END_PHRASE                    = "Кінець"
	TEMPLATE_LECTIONS_LIST        = "%d.%s.%s"
	TEMPLATE_INTRO_LECTIONS_LIST  = "Виберіть лекцію. Для закінчення натисніть " + END_PHRASE
	TEMPLATE_PLACES_LIST          = "%d. %s - %s\n"
	TEMPLATE_PLACES_LIST_BUTTONS  = "%d. %s\n"
	TEMPLATE_INTRO_PLACES_LIST    = "Де?\n%s"
	TEMPLATE_NEXT_EVENT           = "Де: %s, %s\nПочаток: %s\nКінець: %s"
	TEMPLATE_NEXT_EVENT_UNDEFINED = "Невідомо коли, запитайся пізніше"
)

type addEvent struct {
	step        int
	whenStart   time.Time
	whenEnd     time.Time
	where       int
	description string
	lections    []int
}

func (c *addEvent) IsAllow(u string) bool {
	//TODO: impove filter instead of read all records
	admins, err := store.Users([]store.UserRole{store.USER_ROLE_ADMIN})
	if err != nil {
		log.Error("error while reading admins ", err)
		return false
	}
	for _, admin := range admins {
		if admin.Username == u {
			return true
		}
	}
	return false
}

func (c *addEvent) NextStep(u *store.User, answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	switch c.step {
	case 0:
		replyMarkup.Text = "Коли починається? Дата та час в форматі 2018-12-31 19:00:00"
	case 1:
		t, err := time.Parse("2006-01-02 15:04:05", answer)
		if err != nil {
			replyMarkup.Text = "Невірний формат дати та часу. Наприклад, якщо івент буде 20-ого грудня о 19:00 то треба ввести: 2018-12-20 19:00:00. Спробуй ще!"
			return replyMarkup, nil
		}
		c.whenStart = t
		replyMarkup.Text = "Коли закінчується? Дата та час в форматі 2018-12-31 19:00:00"
	case 2:
		t, err := time.Parse("2006-01-02 15:04:05", answer)
		if err != nil {
			replyMarkup.Text = "Невірний формат дати та часу. Наприклад, якщо івент буде 20-ого грудня о 19:00 то треба ввести: 2018-12-20 19:00:00. Спробуй ще!"
			return replyMarkup, nil
		}
		c.whenEnd = t
		places, err := store.Places(store.PlaceTypes{store.PLACE_TYPE_FOR_EVENT, store.PLACE_TYPE_FOR_ALL})
		pText := ""
		replyMarkup.Buttons = nil
		for _, p := range places {
			pText += fmt.Sprintf(TEMPLATE_PLACES_LIST, p.ID, p.Name, p.Address)
			b := fmt.Sprintf(TEMPLATE_PLACES_LIST_BUTTONS, p.ID, p.Name)
			replyMarkup.Buttons = append(replyMarkup.Buttons, b)
		}
		replyMarkup.Text = fmt.Sprintf(TEMPLATE_INTRO_PLACES_LIST, pText)
	case 3:
		regexpLectionID := regexp.MustCompile(`^(\d+)?\.`)
		matches := regexpLectionID.FindStringSubmatch(answer)
		if len(matches) < 2 {
			return nil, ErrWrongPlace
		}
		where, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, errors.New("failed string to int converting")
		}
		c.where = where
		replyMarkup.Text = "Текст івенту"
	case 4:
		c.description = answer
		lections, err := store.Lections(true)
		if err != nil {
			return nil, err
		}
		for _, l := range lections {
			log.Info(l.Lector)
			lText := fmt.Sprintf(TEMPLATE_LECTIONS_LIST, l.ID, l.Name, l.Lector.Name)
			replyMarkup.Buttons = append(replyMarkup.Buttons, lText)
		}
		replyMarkup.Buttons = append(replyMarkup.Buttons, END_PHRASE)
		replyMarkup.Text = TEMPLATE_INTRO_LECTIONS_LIST
	case 5:
		if answer == END_PHRASE {
			_, err := store.AddEvent(c.whenStart, c.whenEnd, c.where, c.description, c.lections)
			if err != nil {
				return nil, err
			}
			replyMarkup.Text = "Івент створено"
			break
		}
		regexpLectionID := regexp.MustCompile(`^(\d+)?\.`)
		matches := regexpLectionID.FindStringSubmatch(answer)
		if len(matches) < 2 {
			return nil, ErrWrongPlace
		}
		lection, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, err
		}
		c.lections = append(c.lections, lection)
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
}

func (c *nextEvent) IsEnd() bool {
	return true
}

func (c *nextEvent) IsAllow(u string) bool {
	return true
}

func (c *nextEvent) NextStep(u *store.User, answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(u.Role),
	}
	e, err := store.NextEvent()
	if err != nil {
		if err == store.ErrUndefinedNextEvent {
			replyMarkup.Text = TEMPLATE_NEXT_EVENT_UNDEFINED
			return replyMarkup, nil
		}
		return nil, err
	}
	replyMarkup.Text = fmt.Sprintf(TEMPLATE_NEXT_EVENT, e.PlaceName, e.Address, e.StartTime.Format("2006-01-02 15:04:05"), e.EndTime.Format("2006-01-02 15:04:05"))
	return replyMarkup, nil
}
