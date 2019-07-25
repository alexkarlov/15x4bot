package commands

import (
	"github.com/alexkarlov/15x4bot/store"
	"math/rand"
)

var unknownMsgs = []string{"Будьте ж людьми, ребята. Ну все ж мы люди", "Шо?", "WTF?", "Ніпанятна", "А тепер подумай и нормально сформулюй питання", "Я не розумію тебе", "Какой самый известный самолет на тихоокеанском театре военных действий?!"}

type unknown struct {
}

func (c *unknown) IsEnd() bool {
	return true
}

func (c *unknown) IsAllow(u string) bool {
	return true
}

func (c *unknown) NextStep(u *store.User, answer string) (*ReplyMarkup, error) {
	text := unknownMsgs[rand.Intn(len(unknownMsgs))]
	replyMsg := &ReplyMarkup{
		Text:    text,
		Buttons: StandardMarkup(u.Role),
	}
	return replyMsg, nil
}
