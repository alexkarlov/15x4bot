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

func (c *simple) NextStep(answer string) (string, error) {
	replyMsg, err := store.ActionMsg(c.action)
	return replyMsg, err
}
