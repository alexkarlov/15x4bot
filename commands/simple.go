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

func (c *simple) NextStep(answer string) (replyMsg string, err error) {
	replyMsg, err = store.GetActionMsg(c.action)

	return
}
