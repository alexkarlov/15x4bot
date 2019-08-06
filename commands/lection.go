package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexkarlov/15x4bot/store"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ErrWrongLection = errors.New("wrong lection id: failed to convert from string to int")
	ErrWrongUserID  = errors.New("wrong user id: failed to convert from string to int")
)

const (
	UserRemindHour = 19

	TEMPLATE_CREATE_EVENT_STEP_SPEAKER_DETAILS     = "%d - %s, %s\n"
	TEMPLATE_CREATE_EVENT_STEP_SPEAKER             = "Хто лектор?\n%s"
	TEMPLATE_CREATE_EVENT_STEP_LECTION_NAME        = "Назва лекції"
	TEMPLATE_CREATE_EVENT_STEP_LECTION_DESCRIPTION = "Опис лекції"
	TEMPLATE_CREATE_EVENT_SUCCESS_MSG              = "Лекцію створено"
	TEMPLATE_LECTION_NAME                          = "Лекція %d: %s"
	TEMPLATE_ADD_LECTION_DESCIRPTION_CHOSE_LECTION = "Оберіть лекцію"
	TEMPLATE_ADD_LECTION_DESCRIPTION_COMPLETE      = "Опис лекції створено"

	TEMPLATE_ADD_LECTION_DESCIRPTION_ERROR_NOT_YOUR = "Це не твоя лекція!"
	TEMPLATE_LECTION_ERROR_WRONG_ID                 = "Невірно вибрана лекція"
	TEMPLATE_WRONG_USER_ID                          = "Невідомий користувач"
	TEMPLATE_LECTION_LIST_ITEM                      = "Лекція %d: %s\nЛектор: @%s,  %s"
	TEMPLATE_LECTION_LIST_EMPTY                     = "Поки лекцій немає"
	TEMPLATE_DELETE_LECTION_COMPLETE                = "Лекцію успішно видалено"
)

func nextDay(hour int) time.Time {
	curr := time.Now()
	y, m, d := curr.Date()
	loc, _ := time.LoadLocation("UTC")
	rTime := time.Date(y, m, d, hour, 0, 0, 0, loc).AddDate(0, 0, 1)
	return rTime
}

type addLection struct {
	u           *store.User
	step        int
	name        string
	description string
	userID      int
}

func (c *addLection) IsAllow(u *store.User) bool {
	c.u = u
	return (u.Role == store.USER_ROLE_ADMIN || u.Role == store.USER_ROLE_LECTOR)
}

func (c *addLection) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	switch c.step {
	case 0:
		// if the user is a lector = add him as an lection owner and skip the next step
		if c.u.Role == store.USER_ROLE_LECTOR {
			c.userID = c.u.ID
			replyMarkup.Text = TEMPLATE_CREATE_EVENT_STEP_LECTION_NAME
			c.step++
			break
		}
		users, err := store.Users([]store.UserRole{store.USER_ROLE_ADMIN, store.USER_ROLE_LECTOR})
		if err != nil {
			return nil, err
		}
		speakerText := ""
		for _, u := range users {
			speakerText += fmt.Sprintf(TEMPLATE_CREATE_EVENT_STEP_SPEAKER_DETAILS, u.ID, u.Username, u.Name)
		}
		replyMarkup.Text = fmt.Sprintf(TEMPLATE_CREATE_EVENT_STEP_SPEAKER, speakerText)
	case 1:
		userID, err := strconv.Atoi(answer)
		if err != nil {
			return nil, ErrWrongUserID
		}
		ok, err := store.DoesUserExist(userID)
		if err != nil {
			return nil, err
		}
		if !ok {
			replyMarkup.Text = TEMPLATE_WRONG_USER_ID
			return replyMarkup, nil
		}
		c.userID = userID
		replyMarkup.Text = TEMPLATE_CREATE_EVENT_STEP_LECTION_NAME
	case 2:
		c.name = answer
		replyMarkup.Text = TEMPLATE_CREATE_EVENT_STEP_LECTION_DESCRIPTION
		replyMarkup.Buttons = append(replyMarkup.Buttons, TEMPLATE_I_DONT_KNOW)
	case 3:
		if answer != TEMPLATE_I_DONT_KNOW {
			c.description = answer
		}
		lectionID, err := store.AddLection(c.name, c.description, c.userID)
		if err != nil {
			return nil, err
		}
		if answer == TEMPLATE_I_DONT_KNOW {
			execTime := nextDay(UserRemindHour)
			r := &store.RemindLection{
				ID: lectionID,
			}
			details, err := json.Marshal(r)
			if err != nil {
				return nil, err
			}
			store.AddTask(store.TASK_TYPE_REMINDER_LECTOR, execTime, string(details))
		}
		replyMarkup.Text = TEMPLATE_CREATE_EVENT_SUCCESS_MSG
	}
	c.step++
	return replyMarkup, nil
}

func (c *addLection) IsEnd() bool {
	return c.step == 4
}

type addDescriptionLection struct {
	u         *store.User
	lectionID int
	step      int
}

func (c *addDescriptionLection) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	switch c.step {
	case 0:
		lections, err := store.Lections(true)
		if err != nil {
			return nil, err
		}
		var l []string
		for _, lection := range lections {
			if (c.u.Role != store.USER_ROLE_ADMIN && c.u.ID != lection.Lector.ID) || lection.Description != "" {
				// skip lections which doesn't belong to user (if he isn't admin) or it has description
				continue
			}
			l = append(l, fmt.Sprintf(TEMPLATE_LECTION_NAME, lection.ID, lection.Name))
		}
		// if there are no appropriate lections - send special response
		if len(l) == 0 {
			replyMarkup.Text = TEMPLATE_LECTION_LIST_EMPTY
			// TODO: OMG, remove that shit
			c.step = 3
			return replyMarkup, nil
		}
		replyMarkup.Buttons = append(replyMarkup.Buttons, l...)
		replyMarkup.Text = TEMPLATE_ADD_LECTION_DESCIRPTION_CHOSE_LECTION
	case 1:
		regexpLectionID := regexp.MustCompile(`^Лекція (\d+)?\:`)
		matches := regexpLectionID.FindStringSubmatch(answer)
		if len(matches) < 2 {
			replyMarkup.Text = TEMPLATE_LECTION_ERROR_WRONG_ID
			return replyMarkup, nil
		}
		lID, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, ErrWrongLection
		}
		l, err := store.LoadLection(lID)
		if err != nil {
			return nil, err
		}
		if c.u.Role != store.USER_ROLE_ADMIN && l.Lector.Username != c.u.Username {
			replyMarkup.Text = TEMPLATE_ADD_LECTION_DESCIRPTION_ERROR_NOT_YOUR
			return replyMarkup, nil
		}
		c.lectionID = lID
		replyMarkup.Text = TEMPLATE_CREATE_EVENT_STEP_LECTION_DESCRIPTION
	case 2:
		err := store.AddLectionDescription(c.lectionID, answer)
		if err != nil {
			return nil, err
		}
		replyMarkup.Text = TEMPLATE_ADD_LECTION_DESCRIPTION_COMPLETE
	}
	c.step++
	return replyMarkup, nil
}

func (c *addDescriptionLection) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN || u.Role == store.USER_ROLE_LECTOR
}

func (c *addDescriptionLection) IsEnd() bool {
	return c.step == 3
}

type lectionsList struct {
	u                  *store.User
	withoutDescription bool
}

func (c *lectionsList) IsEnd() bool {
	return true
}

func (c *lectionsList) IsAllow(u *store.User) bool {
	c.u = u
	return (u.Role == store.USER_ROLE_ADMIN || u.Role == store.USER_ROLE_LECTOR)
}

func (c *lectionsList) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	list, err := store.Lections(true)
	if err != nil {
		return nil, err
	}
	var l []string
	for _, lection := range list {
		if c.withoutDescription && lection.Description != "" {
			// if we want to see only lections without descriptions and the current lection does have description - skip it
			continue
		}
		if c.u.Role != store.USER_ROLE_ADMIN && c.u.ID != lection.Lector.ID {
			// skip lections which doesn't belong to user (if he isn't admin)
			continue
		}
		l = append(l, fmt.Sprintf(TEMPLATE_LECTION_LIST_ITEM, lection.ID, lection.Name, lection.Lector.Username, lection.Lector.Name))
	}
	// if there are no appropriate lections - send special response
	if len(l) == 0 {
		replyMarkup.Text = TEMPLATE_LECTION_LIST_EMPTY
		return replyMarkup, nil
	}
	replyMarkup.Text = strings.Join(l, "\n\n")
	return replyMarkup, nil
}

type deleteLection struct {
	step      int
	lectionID int
	u         *store.User
}

func (c *deleteLection) IsEnd() bool {
	return c.step == 2
}

func (c *deleteLection) IsAllow(u *store.User) bool {
	c.u = u
	return u.Role == store.USER_ROLE_ADMIN || u.Role == store.USER_ROLE_LECTOR
}

func (c *deleteLection) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: MainMarkup,
	}
	switch c.step {
	case 0:
		lections, err := store.Lections(false)
		if err != nil {
			return nil, err
		}

		var l []string
		for _, lection := range lections {
			if c.u.Role != store.USER_ROLE_ADMIN && c.u.ID != lection.Lector.ID {
				// skip lections which doesn't belong to user (if he isn't admin)
				continue
			}
			l = append(l, fmt.Sprintf(TEMPLATE_LECTION_NAME, lection.ID, lection.Name))
		}
		// if there are no appropriate lections - send special response
		if len(l) == 0 {
			replyMarkup.Text = TEMPLATE_LECTION_LIST_EMPTY
			// TODO: OMG, remove that shit
			c.step = 2
			return replyMarkup, nil
		}
		replyMarkup.Buttons = append(replyMarkup.Buttons, l...)
		replyMarkup.Text = TEMPLATE_ADD_LECTION_DESCIRPTION_CHOSE_LECTION
	case 1:
		regexpLectionID := regexp.MustCompile(`^Лекція (\d+)?\:`)
		matches := regexpLectionID.FindStringSubmatch(answer)
		if len(matches) < 2 {
			replyMarkup.Text = TEMPLATE_LECTION_ERROR_WRONG_ID
			return replyMarkup, nil
		}
		lID, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, ErrWrongLection
		}
		l, err := store.LoadLection(lID)
		if err != nil {
			return nil, err
		}
		if c.u.Role != store.USER_ROLE_ADMIN && c.u.ID != l.Lector.ID {
			replyMarkup.Text = TEMPLATE_ADD_LECTION_DESCIRPTION_ERROR_NOT_YOUR
			return replyMarkup, nil
		}
		err = store.DeleteLection(lID)
		if err != nil {
			return nil, err
		}
		replyMarkup.Buttons = StandardMarkup(c.u.Role)
		replyMarkup.Text = TEMPLATE_DELETE_LECTION_COMPLETE
	}
	c.step++
	return replyMarkup, nil
}
