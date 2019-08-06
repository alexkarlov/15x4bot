package scheduler

import (
	"fmt"
	"github.com/alexkarlov/15x4bot/bot"
	"github.com/alexkarlov/15x4bot/store"
	"github.com/alexkarlov/simplelog"
)

var (
	weekdays = map[string]string{
		"Monday":    "Понеділок",
		"Tuesday":   "Вівторок",
		"Wednesday": "Середа",
		"Thursday":  "Четвер",
		"Friday":    "П'ятниця",
		"Saturday":  "Субота",
		"Sunday":    "Неділя",
	}
	months = map[string]string{
		"January":   "Січень",
		"February":  "Лютий",
		"March":     "Березень",
		"April":     "Квітень",
		"May":       "Травень",
		"June":      "Червень",
		"July":      "Липень",
		"August":    "Серпень",
		"September": "Вересень",
		"October":   "Жовтень",
		"November":  "Листопад",
		"December":  "Грудень",
	}
)

const (
	TimeLayout                        = "15:04"
	TEMPLATE_REHEARSAL_MSG_TO_CHANNEL = "Привіт! Нова репетиція\nДе: %s\nКоли: %s, %d %s, %s\nАдреса: %s\nМапа: %s\n"
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
	rh, err := t.LoadRemindRehearsalChannel()
	if err != nil {
		log.Errorf("failed to load reminder rehearsal to channel of task %d error:%s", t.ID, err)
		if err := t.ReleaseTask(); err != nil {
			log.Errorf("failed to release task %d error:%s", t.ID, err)
		}
		return
	}
	r, err := store.LoadRehearsal(rh.RehearsalID)
	if err != nil {
		log.Errorf("failed to load rehearsal of task %d error:%s", t.ID, err)
		if err := t.ReleaseTask(); err != nil {
			log.Errorf("failed to release task %d error:%s", t.ID, err)
		}
		return
	}
	wd := weekdays[r.Time.Weekday().String()]
	m := months[r.Time.Month().String()]
	msg := fmt.Sprintf(TEMPLATE_REHEARSAL_MSG_TO_CHANNEL, r.PlaceName, wd, r.Time.Day(), m, r.Time.Format(TimeLayout), r.Address, r.MapUrl)
	if err := b.SendTextToChannel(rh.ChannelUsername, msg); err != nil {
		log.Errorf("error while sending msg to %s. task %d error: %s", rh.ChannelUsername, t.ID, err)
		if err := t.ReleaseTask(); err != nil {
			log.Errorf("failed to release task %d error:%s", t.ID, err)
		}
		return
	}
	if err := t.FinishTask(); err != nil {
		log.Errorf("failed to release task %d error:%s", t.ID, err)
	}
}
