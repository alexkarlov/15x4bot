package scheduler

import (
	"time"

	"github.com/alexkarlov/simplelog"
)

// Run starts checking background tasks from db every minute
func Run() {
	for {
		log.Info("checking for scheduler task...")
		// check tasks in db
		log.Info("no tasks so far")
		time.Sleep(time.Minute * 1)
	}
}
