package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexkarlov/15x4bot/store"
	"strconv"
	"time"
)

const (
	TEMPLATE_CREATE_EVENT_STEP_SPEAKER_DETAILS      = "%d - %s, %s\n"
	TEMPLATE_CREATE_EVENT_STEP_SPEAKER              = "Хто лектор?\n%s"
	TEMPLATE_CREATE_EVENT_STEP_LECTION_NAME         = "Назва лекції"
	TEMPLATE_CREATE_EVENT_STEP_LECTION_DESCRIPTION  = "Опис лекції"
	TEMPLATE_CREATE_EVENT_SUCCESS_MSG               = "Лекцію створено"
	TEMPLATE_ADD_LECTION_DESCIRPTION_ERROR_NOT_YOUR = "Це не твоя лекція!"
)

type addLection struct {
	step        int
	name        string
	description string
	user_id     int
}

type addDescriptionLection struct {
	description string
	username    string
	taskID      int
	step        int
}

func (c *addLection) IsAllow(u string) bool {
	//TODO: move it to db
	admins := []string{"zedman95", "alex_karlov"}
	for _, admin := range admins {
		if admin == u {
			return true
		}
	}
	return false
}

func (c *addLection) NextStep(u *store.User, answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{}
	switch c.step {
	case 0:
		users, err := store.Users([]store.UserRole{store.USER_ROLE_ADMIN, store.USER_ROLE_LECTOR})
		if err != nil {
			return nil, err
		}
		speakerText := ""
		for _, u := range users {
			speakerText += fmt.Sprintf(TEMPLATE_CREATE_EVENT_STEP_SPEAKER_DETAILS, u.ID, u.Username, u.Name)
			replyMarkup.Buttons = append(replyMarkup.Buttons, strconv.Itoa(u.ID))
		}
		replyMarkup.Text = fmt.Sprintf(TEMPLATE_CREATE_EVENT_STEP_SPEAKER, speakerText)
	case 1:
		// TODO: validate it
		userID, err := strconv.Atoi(answer)
		if err != nil {
			return nil, errors.New("failed string to int converting")
		}
		c.user_id = userID
		replyMarkup.Text = TEMPLATE_CREATE_EVENT_STEP_LECTION_NAME
	case 2:
		c.name = answer
		replyMarkup.Text = TEMPLATE_CREATE_EVENT_STEP_LECTION_DESCRIPTION
	case 3:
		c.description = answer
		lectionID, err := store.AddLection(c.name, c.description, c.user_id)
		if err != nil {
			return nil, err
		}
		if c.description == "-" {
			execTime := lectionRemindTime()
			l := &store.Lection{
				ID: lectionID,
			}
			details, err := json.Marshal(l)
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

func lectionRemindTime() time.Time {
	curr := time.Now()
	y, m, d := curr.Date()
	loc, _ := time.LoadLocation("UTC")
	rTime := time.Date(y, m, d, 19, 0, 0, 0, loc).AddDate(0, 0, 1)
	return rTime
}

func (c *addDescriptionLection) NextStep(u *store.User, answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{}
	t, err := store.LoadTask(c.taskID)
	if err != nil {
		return nil, err
	}
	l, err := t.LoadLection()
	if err != nil {
		return nil, err
	}
	if l.Lector.Username != c.username {
		replyMarkup.Text = TEMPLATE_ADD_LECTION_DESCIRPTION_ERROR_NOT_YOUR
		return replyMarkup, nil
	}
	err = l.AddDescriptionLection(c.description)
	if err != nil {
		return nil, err
	}
	err = store.FinishTask(c.taskID)
	if err != nil {
		return nil, err
	}
	replyMarkup.Text = "Опис лекції створено"
	return replyMarkup, nil
}

func (c *addDescriptionLection) IsAllow(u string) bool {
	c.username = u
	//TODO: check whether it's a lector or admin
	return true
}

func (c *addDescriptionLection) IsEnd() bool {
	return c.step == 1
}
