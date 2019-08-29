package scheduler

import (
	"github.com/alexkarlov/15x4bot/bot"
	"github.com/alexkarlov/15x4bot/store"
	"github.com/alexkarlov/simplelog"
)

// MessageToChannel sends message to channel
func MessageToChannel(t *store.Task, b *bot.Bot) {
	log.Info("got new message to channel:", t)
	if err := t.TakeTask(); err != nil {
		log.Errorf("failed to take task %d error:%s", t.ID, err)
		if err := t.ReleaseTask(); err != nil {
			log.Errorf("failed to release task %d error:%s", t.ID, err)
		}
		return
	}
	rc, err := t.LoadRemindChannel()
	if err != nil {
		log.Errorf("failed to load reminder rehearsal to channel of task %d error:%s", t.ID, err)
		if err := t.ErrorTask(); err != nil {
			log.Errorf("failed to release task %d error:%s", t.ID, err)
		}
		return
	}
	if err := b.SendMsgToChannel(rc.ChannelUsername, rc.Msg, rc.FileIDs); err != nil {
		log.Errorf("error while sending msg to %s. task %d error: %s", rc.ChannelUsername, t.ID, err)
		if err := t.ErrorTask(); err != nil {
			log.Errorf("failed to release task %d error:%s", t.ID, err)
		}
		return
	}
	if err := t.FinishTask(); err != nil {
		log.Errorf("failed to release task %d error:%s", t.ID, err)
	}
}
