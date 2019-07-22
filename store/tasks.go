package store

import (
	"encoding/json"
	"errors"
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

	// Execution time range filter
	EXECUTION_TIME_FILTER_INTERVAL = "INTERVAL '2 MINUTES'"

	POSTPONE_PERIOD_ONE_DAY  postponePeriod = "INTERVAL '1 DAY'"
	POSTPONE_PERIOD_TWO_DAYS postponePeriod = "INTERVAL '2 DAY'"
	POSTPONE_PERIOD_ONE_WEEK postponePeriod = "INTERVAL '1 WEEK'"
)

var ErrUndefinedTaskType = errors.New("Undefined task type")

type TaskType string

type StatusType int

type postponePeriod string

type Task struct {
	ID            int
	Type          TaskType
	ExecutionTime time.Time
	Status        StatusType
	Details       string
}

// RemindLection contains details (ID) about the lection to ask description from user
type RemindLection struct {
	ID int
}

func GetTasks() ([]*Task, error) {
	execTimeFilter := " AND execution_time>=(NOW()- " + EXECUTION_TIME_FILTER_INTERVAL + ") AND execution_time<=(NOW()+ " + EXECUTION_TIME_FILTER_INTERVAL + ")"
	q := "SELECT id, type, execution_time, details FROM tasks WHERE status=$1 "
	q += execTimeFilter
	rows, err := dbConn.Query(q, TASK_STATUS_NEW)
	if err != nil {
		return nil, err
	}
	tasks := make([]*Task, 0)
	for rows.Next() {
		t := &Task{}
		if err := rows.Scan(&t.ID, &t.Type, &t.ExecutionTime, &t.Details); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// LoadTask selects specified task from db by id
func LoadTask(ID int) (*Task, error) {
	t := &Task{}
	q := "SELECT id, type, execution_time, status, details FROM tasks WHERE id=$1"
	err := dbConn.QueryRow(q, ID).Scan(&t.ID, &t.Type, &t.ExecutionTime, &t.Status, &t.Details)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func FinishTask(ID int) error {
	_, err := dbConn.Exec("UPDATE tasks SET udate=NOW(), status=$1 WHERE id=$2", TASK_STATUS_DONE, ID)
	if err != nil {
		return err
	}
	return nil
}

// PostponeTask updates udate field to the current time, set status to NEW and update execution_time if a task fulfills conditions
func PostponeTask(ID int, postponePeriod postponePeriod) error {
	q := "UPDATE tasks SET udate=NOW(), execution_time=execution_time+" + string(postponePeriod) + ", status=$1 WHERE id=$2 AND status=$3"
	_, err := dbConn.Exec(q, TASK_STATUS_NEW, ID, TASK_STATUS_IN_PROGRESS)
	if err != nil {
		return err
	}
	return nil
}

// AddTask creates new task with details. Each type of task contains specific details
func AddTask(t TaskType, execTime time.Time, details string) error {
	_, err := dbConn.Exec("INSERT INTO tasks (type, execution_time, status, details) VALUES ($1, $2, $3, $4)", t, execTime, TASK_STATUS_NEW, details)
	if err != nil {
		return err
	}
	return nil
}

// LoadLection loads lection from task details
func (t *Task) LoadLection() (*Lection, error) {
	r := &RemindLection{}
	err := json.Unmarshal([]byte(t.Details), r)
	if err != nil {
		return nil, err
	}
	l, err := LoadLection(r.ID)
	return l, err
}

// TakeTask updates udate field to the current time and set status to IN_PROGRESS if a task fulfills conditions
func (t *Task) TakeTask() error {
	_, err := dbConn.Exec("UPDATE tasks SET udate=NOW(), status=$1 WHERE id=$2 AND status=$3", TASK_STATUS_IN_PROGRESS, t.ID, TASK_STATUS_NEW)
	if err != nil {
		return err
	}
	return nil
}
