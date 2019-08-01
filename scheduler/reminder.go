package scheduler

import (
	"github.com/alexkarlov/15x4bot/bot"
	"github.com/alexkarlov/15x4bot/store"
	"github.com/alexkarlov/simplelog"
)

const (
	TEMPLATE_LECTION_DESCRIPTION_REMINDER = `Привіт!
	По можливості - напиши, будь ласка, опис до своєї лекції (два-три речення про що буде лекція). В головному меню є пункт "Додати опис до лекції". Якщо будуть питання - звертайся до @alex_karlov
	Дякую!
	`
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
	r, err := t.LoadReminderLection()
	if err != nil {
		log.Errorf("failed to load reminder lection of task %d error:%s", t.ID, err)
		if err := t.ReleaseTask(); err != nil {
			log.Errorf("failed to release task %d error:%s", t.ID, err)
		}
		return
	}
	l, err := r.LoadLection()
	if err != nil {
		log.Errorf("failed to load lection of task %d error:%s", t.ID, err)
		if err := t.ReleaseTask(); err != nil {
			log.Errorf("failed to release task %d error:%s", t.ID, err)
		}
		return
	}
	if l.Description != "" {
		if err = store.FinishTask(t.ID); err != nil {
			log.Errorf("failed to finish task %d error:%s", t.ID, err)
		}
		return
	}
	c, err := l.Lector.TGChat()
	if err != nil {
		log.Error("error while getting tg chat of the lector", err)
		if err := t.ReleaseTask(); err != nil {
			log.Errorf("failed to release task %d error:%s", t.ID, err)
		}
		return
	}
	b.SendText(c.TGChatID, TEMPLATE_LECTION_DESCRIPTION_REMINDER)
	// Udate task with new execution time and attempts
	if err = r.PostponeTask(t.ID); err != nil {
		log.Errorf("failed to postpone task %d error:%s", t.ID, err)
		return
	}
	log.Info(t)
}
