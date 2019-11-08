package scheduler

import (
	"fmt"
	"github.com/alexkarlov/15x4bot/bot"
	"github.com/alexkarlov/15x4bot/lang"
	"github.com/alexkarlov/15x4bot/store"
	"github.com/alexkarlov/simplelog"
)

// RemindLector sends message to the speaker about description of his lecture
func RemindLector(t *store.Task, b *bot.Bot) {
	log.Info("got new reminder lector:", t)
	if err := t.TakeTask(); err != nil {
		log.Errorf("failed to take task %d error:%s", t.ID, err)
		if err := t.ReleaseTask(); err != nil {
			log.Errorf("failed to release task %d error:%s", t.ID, err)
		}
		return
	}
	r, err := t.LoadReminderLecture()
	if err != nil {
		log.Errorf("failed to load reminder lecture of task %d error:%s", t.ID, err)
		if err := t.ReleaseTask(); err != nil {
			log.Errorf("failed to release task %d error:%s", t.ID, err)
		}
		return
	}
	l, err := r.LoadLecture()
	if err != nil {
		log.Errorf("failed to load lecture of task %d error:%s", t.ID, err)
		if err := t.ReleaseTask(); err != nil {
			log.Errorf("failed to release task %d error:%s", t.ID, err)
		}
		return
	}
	if l.Description != "" {
		if err = t.FinishTask(); err != nil {
			log.Errorf("failed to finish task %d error:%s", t.ID, err)
		}
		return
	}
	// skip unregistered users (bot hasn't spoken with them yet)
	if l.Lector.TGUserID == 0 {
		log.Info(fmt.Sprintf("reminder skip the user %d since it doesn't have tg id", l.Lector.ID))
		return
	}
	b.SendText(int64(l.Lector.TGUserID), lang.TEMPLATE_LECTURE_DESCRIPTION_REMINDER)
	// Udate task with new execution time and attempts
	if err = r.PostponeTask(t.ID); err != nil {
		log.Errorf("failed to postpone task %d error:%s", t.ID, err)
		return
	}
	log.Info(t)
}
