package commands

import (
	"encoding/json"
	"github.com/alexkarlov/15x4bot/lang"
	"github.com/alexkarlov/15x4bot/store"
	"github.com/alexkarlov/simplelog"
	"time"
)

// messenger fires when a new user wants to become a speaker or volunteer
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
		reply.Text = lang.MESSENGER_USERNAME_EMPTY
		return reply, nil
	}
	loc, err := time.LoadLocation(Conf.Location)
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
	reply.Text = lang.MESSENGER_THANKS
	return reply, nil
}
