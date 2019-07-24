package commands

import (
	"github.com/alexkarlov/15x4bot/store"
	"github.com/alexkarlov/simplelog"
)

type simple struct {
	action string
	user   *User
}

func (c *simple) IsEnd() bool {
	return true
}

func (c *simple) IsAllow(u string) bool {
	var err error
	user, err := store.LoadUser(u)
	c.user = &User{
		User: user,
	}
	if err != nil {
		log.Error("error while getting user", err)
	}
	return true
}

func (c *simple) NextStep(answer string) (*ReplyMarkup, error) {
	replyMarkup := &ReplyMarkup{}
	replyMsg, err := store.ActionMsg(c.action)
	replyMarkup.Text = replyMsg
	replyMarkup.Buttons = c.user.Markup()
	return replyMarkup, err
}
