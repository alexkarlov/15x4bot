package scheduler

import (
	"fmt"
	"github.com/alexkarlov/15x4bot/bot"
	"github.com/alexkarlov/15x4bot/store"
	"github.com/alexkarlov/simplelog"
)

const (
	TEMPLATE_MSG_TO_ADMIN_NEW_VOLUNTEER_COME = "Йу-ху! Новий %s хоче до нас долучитися - %s"
)

// MessageToAdmin sends message to admin(s) when new speaker/volunteer comes
func MessageToAdmin(t *store.Task, b *bot.Bot) {
	log.Info("got new message to admin:", t)
	if err := t.TakeTask(); err != nil {
		log.Errorf("failed to take task %d error:%s", t.ID, err)
		if err := t.ReleaseTask(); err != nil {
			log.Errorf("failed to release task %d error:%s", t.ID, err)
		}
		return
	}
	admins, err := store.Users([]store.UserRole{store.USER_ROLE_ADMIN})
	if err != nil {
		log.Error("error while getting admins", err)
		if err := t.ReleaseTask(); err != nil {
			log.Errorf("failed to release task %d error:%s", t.ID, err)
		}
		return
	}
	m, err := t.LoadMessenger()
	if err != nil {
		log.Errorf("failed to load reminder lection of task %d error:%s", t.ID, err)
		if err := t.ReleaseTask(); err != nil {
			log.Errorf("failed to release task %d error:%s", t.ID, err)
		}
		return
	}
	msg := fmt.Sprintf(TEMPLATE_MSG_TO_ADMIN_NEW_VOLUNTEER_COME, m.Role, m.Username)
	for _, a := range admins {
		chat, err := a.TGChat()
		if err != nil {
			log.Errorf("error while getting a chat id. task %d error: %s", t.ID, err)
		}
		if err := b.SendText(chat.TGChatID, msg); err != nil {
			log.Errorf("error while sending msg to %s. task %d error: %s", a.Username, t.ID, err)
		}
	}
	if err := t.ReleaseTask(); err != nil {
		log.Errorf("failed to release task %d error:%s", t.ID, err)
	}
}
