package commands

import (
	"encoding/json"
	"errors"
	"github.com/alexkarlov/15x4bot/store"
	"strconv"
	"strings"
	"time"
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

func (c *addLection) NextStep(answer string) (string, error) {
	replyMsg := ""
	switch c.step {
	case 0:
		users, err := store.GetUsers([]store.UserRole{store.USER_ROLE_ADMIN, store.USER_ROLE_LECTOR})
		if err != nil {
			return "", err
		}
		replyMsg = strings.Join([]string{"Хто лектор?", strings.Join(users, "\n")}, "\n")
	case 1:
		// TODO: validate it
		userID, err := strconv.Atoi(answer)
		if err != nil {
			return "", errors.New("failed string to int converting")
		}
		c.user_id = userID
		replyMsg = "Назва лекції"
	case 2:
		c.name = answer
		replyMsg = "Опис лекції"
	case 3:
		c.description = answer
		lectionID, err := store.AddLection(c.name, c.description, c.user_id)
		if err != nil {
			return "", err
		}
		if c.description == "-" {
			execTime := lectionRemindTime()
			l := &store.Lection{
				ID: lectionID,
			}
			details, err := json.Marshal(l)
			if err != nil {
				return "", err
			}
			store.AddTask(store.TASK_TYPE_REMINDER_LECTOR, execTime, string(details))
		}
		replyMsg = "Лекцію створено"
	}
	c.step++
	return replyMsg, nil
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

func (c *addDescriptionLection) NextStep(answer string) (string, error) {
	t, err := store.LoadTask(c.taskID)
	if err != nil {
		return "", err
	}
	l := &store.Lection{}
	// TODO: implement json.Unmarshal
	err = json.Unmarshal([]byte(t.Details), l)
	replyMsg := ""
	if err != nil {
		replyMsg = "wrong task"
		return replyMsg, err
	}
	// TODO: refactor it - add methods for Lection instead of function
	if !l.OwnedBy(c.username) {
		replyMsg = "this is not your lection"
		return replyMsg, nil
	}
	err = l.AddDescriptionLection(c.description)
	if err != nil {
		return "", err
	}
	err = store.FinishTask(c.taskID)
	if err != nil {
		return "", err
	}
	replyMsg = "Опис лекції створено"
	return replyMsg, nil
}

func (c *addDescriptionLection) IsAllow(u string) bool {
	c.username = u
	//TODO: check whether it's a lector or admin
	return true
}

func (c *addDescriptionLection) IsEnd() bool {
	return c.step == 1
}
