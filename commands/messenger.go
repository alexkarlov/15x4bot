package commands

import (
	"encoding/json"
	"github.com/alexkarlov/15x4bot/store"
	"github.com/alexkarlov/simplelog"
	"time"
)

const (
	TEMPLATE_MESSENGER_THANKS         = "Дякую! Я передав інформацію організаторам"
	TEMPLATE_MESSENGER_USERNAME_EMPTY = "Напиши, будь ласка, @alex_karlov ! Він розповість що робити далі)"
)

type messenger struct {
	u    *store.User
	role string
}

func (c *messenger) IsEnd() bool {
	return true
}

func (c *messenger) IsAllow(u *store.User) bool {
	c.u = u
	return true
}

func (c *messenger) NextStep(answer string) (*ReplyMarkup, error) {
	reply := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	if c.u.Username == "" {
		reply.Text = TEMPLATE_MESSENGER_USERNAME_EMPTY
		return reply, nil
	}
	loc, err := time.LoadLocation(MsgLocation)
	if err != nil {
		log.Error(err)
	}
	execTime := time.Now().In(loc)
	r := &store.Messenger{
		Username: c.u.Username,
		Role:     c.role,
	}
	details, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	store.AddTask(store.TASK_TYPE_MESSENGER, execTime, string(details))
	reply.Text = TEMPLATE_MESSENGER_THANKS
	return reply, nil
}
