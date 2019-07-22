package commands

import (
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

func (c *unknown) NextStep(answer string) (string, error) {
	replyMsg := unknownMsgs[rand.Intn(len(unknownMsgs))]
	return replyMsg, nil
}
