package commands

import (
	"reminder/db"

	"gopkg.in/telegram-bot-api.v4"
)

type simple struct {
	action string
}

func (c *simple) IsEnd() bool {
	return true
}

func (c *simple) IsAllow(u *tgbotapi.User) bool {
	return true
}

func (c *simple) NextStep(answer string) (replyMsg string, err error) {
	replyMsg, err = db.GetActionMsg(c.action)

	return
}
