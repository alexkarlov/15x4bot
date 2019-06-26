package commands

import (
	"math/rand"
)

type unknown struct {
}

func (c *unknown) IsEnd() bool {
	return true
}

func (c *unknown) IsAllow(u string) bool {
	return true
}

func (c *unknown) NextStep(answer string) (replyMsg string, err error) {
	msg := []string{"Будьте ж людьми, ребята. Ну все ж мы люди", "Шо?", "WTF?", "Ніпанятна", "А тепер подумай и нормально сформулюй питання", "Я не розумію тебе", "Какой самый известный самолет на тихоокеанском театре военных действий?!"}
	replyMsg = msg[rand.Intn(len(msg))]
	return
}
