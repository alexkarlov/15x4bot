package scheduler

import (
	"encoding/json"
	"github.com/alexkarlov/15x4bot/bot"
	"github.com/alexkarlov/15x4bot/store"
	"strconv"
	"time"

	"github.com/alexkarlov/simplelog"
)

// Run starts checking background tasks from db every minute
func Run(b *bot.Bot) {
	for {
		log.Info("checking for scheduler task...")
		// check tasks in db
		tasks, err := store.GetTasks()
		if err != nil {
			log.Error("error while getting tasks:", err)
			// TODO: refactor it
			time.Sleep(time.Minute * 1)
			continue
		}
		for _, t := range tasks {
			switch t.Type {
			case store.TASK_TYPE_REMINDER_LECTOR:
				// TODO: add task manager
				go RemindLector(t, b)
			}
		}
		time.Sleep(time.Minute * 1)
	}
}

func RemindLector(t *store.Task, b *bot.Bot) {
	log.Info("got new reminder lector:", t)
	err := store.TakeTask(t.ID)
	if err != nil {
		log.Errorf("failed to take task %d error:%s", t.ID, err)
	}
	l := &store.Lection{}
	err = json.Unmarshal([]byte(t.Details), l)
	if err != nil {
		log.Error("error while reminding lector", err)
		return
	}
	u, err := l.Lector()
	if err != nil {
		// TODO: refactor it
		log.Error("error while getting lector of the lection", err)
		return
	}
	c, err := u.TGChat()
	if err != nil {
		// TODO: refactor it
		log.Error("error while getting tg chat of the lector", err)
		return
	}
	msg := "hello from reminder:" + strconv.Itoa(t.ID)
	b.SendText(c.TGChatID, msg)
	store.PostponeTask(t.ID, store.POSTPONE_PERIOD_ONE_DAY)
	log.Info(t)
}
