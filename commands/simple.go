package commands

import (
	"github.com/alexkarlov/15x4bot/store"
)

type simple struct {
	action string
}

func (c *simple) IsEnd() bool {
	return true
}

func (c *simple) IsAllow(u string) bool {
	return true
}

func (c *simple) NextStep(u *store.User, answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{
		Buttons: StandardMarkup(u.Role),
	}
	replyMsg, err := store.ActionMsg(c.action)
	replyMarkup.Text = replyMsg
	return replyMarkup, err
}
