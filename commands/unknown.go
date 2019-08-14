package commands

import (
	"github.com/alexkarlov/15x4bot/lang"
	"github.com/alexkarlov/15x4bot/store"
	"math/rand"
)

type unknown struct {
	u *store.User
}

func (c *unknown) IsEnd() bool {
	return true
}

func (c *unknown) IsAllow(u *store.User) bool {
	c.u = u
	return true
}

func (c *unknown) NextStep(answer string) (*ReplyMarkup, error) {
	text := lang.UnknownMsgs[rand.Intn(len(lang.UnknownMsgs))]
	replyMsg := &ReplyMarkup{
		Text:    text,
		Buttons: StandardMarkup(c.u.Role),
	}
	return replyMsg, nil
}
