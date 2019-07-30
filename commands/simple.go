package commands

import (
	"github.com/alexkarlov/15x4bot/store"
)

type simple struct {
	action string
	u      *store.User
}

func (c *simple) IsEnd() bool {
	return true
}

func (c *simple) IsAllow(u *store.User) bool {
	c.u = u
	return true
}

func (c *simple) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(c.u.Role),
	}
	replyMsg, err := store.ActionMsg(c.action)
	replyMarkup.Text = replyMsg
	return replyMarkup, err
}
