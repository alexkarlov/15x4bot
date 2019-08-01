package scheduler

import (
	"github.com/alexkarlov/15x4bot/bot"
	"github.com/alexkarlov/15x4bot/store"
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
			case store.TASK_TYPE_MESSENGER:
				// TODO: add task manager
				go MessageToAdmin(t, b)
			}
		}
		time.Sleep(time.Minute * 1)
	}
}
