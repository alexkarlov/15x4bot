package commands

import (
	"github.com/alexkarlov/15x4bot/store"
	"math/rand"
)

var unknownMsgs = []string{"Вибач, я не розумію тебе", "Ніпанятна", "Шта?"}

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
	text := unknownMsgs[rand.Intn(len(unknownMsgs))]
	replyMsg := &ReplyMarkup{
		Text:    text,
		Buttons: StandardMarkup(c.u.Role),
	}
	return replyMsg, nil
}
