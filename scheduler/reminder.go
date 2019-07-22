package scheduler

import (
	"fmt"
	"github.com/alexkarlov/15x4bot/bot"
	"github.com/alexkarlov/15x4bot/store"
	"github.com/alexkarlov/simplelog"
)

const (
	TEMPLATE_LECTION_DESCRIPTION_REMINDER = `Привіт, %username%!
	По можливості - напиши, будь ласка, опис до своєї лекції. Щоб я зрозумів тебе правильно, напиши в правильному форматі:
	task_%d:Опис твоєї лекції

	Дякую велетенське!
	`
)

// RemindLector sends message to the speaker about description of his lecture
func RemindLector(t *store.Task, b *bot.Bot) {
	log.Info("got new reminder lector:", t)
	err := t.TakeTask()
	if err != nil {
		log.Errorf("failed to take task %d error:%s", t.ID, err)
		return
	}
	l, err := t.LoadLection()
	if err != nil {
		log.Errorf("failed to load lection of task %d error:%s", t.ID, err)
		return
	}
	c, err := l.Lector.TGChat()
	if err != nil {
		log.Error("error while getting tg chat of the lector", err)
		return
	}
	msg := fmt.Sprintf(TEMPLATE_LECTION_DESCRIPTION_REMINDER, t.ID)
	b.SendText(c.TGChatID, msg)
	store.PostponeTask(t.ID, store.POSTPONE_PERIOD_ONE_DAY)
	log.Info(t)
}
