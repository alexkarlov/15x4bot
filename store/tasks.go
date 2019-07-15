package store

import (
	"time"
)

const (
	// values for task type
	TASK_TYPE_REMINDER_LECTOR       TaskType = "reminder_lector"
	TASK_TYPE_REMINDER_DESIGNER     TaskType = "reminder_designer"
	TASK_TYPE_REMINDER_GRAMMAR_NAZI TaskType = "reminder_grammar_nazi"
	TASK_TYPE_REMINDER_FB_EVENT     TaskType = "reminder_fb_event"
	TASK_TYPE_REMINDER_TG_CHAT      TaskType = "post_tg_chat"
	TASK_TYPE_REMINDER_TG_CHANNEL   TaskType = "post_tg_channel"

	// values for status
	TASK_STATUS_NEW         StatusType = 1
	TASK_STATUS_DONE        StatusType = 2
	TASK_STATUS_IN_PROGRESS StatusType = 3
	TASK_STATUS_ERROR       StatusType = 4
)

type TaskType string

type StatusType int

type Task struct {
	ID            int
	Type          TaskType
	ExecutionTime time.Time
	Status        StatusType
}

func AddTask(t TaskType, excTime time.Time) error {
	return nil
}
