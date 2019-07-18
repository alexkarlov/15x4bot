package store

import (
	"database/sql"
	"errors"
	"time"
)

var ErrUndefinedNextEvent = errors.New("Next event is undeffined")

type Event struct {
	Id          int
	StartTime   time.Time
	EndTime     time.Time
	PlaceName   string
	Address     string
	Description string
	Letions     []int
}

func AddEvent(startTime time.Time, endTime time.Time, place int, description string, lections []int) (int, error) {
	var eventID int
	tx, err := dbConn.Begin()
	if err != nil {
		return eventID, err
	}
	err = tx.QueryRow("INSERT INTO events (starttime, endtime, place, description) VALUES ($1, $2, $3, $4) RETURNING id", startTime, endTime, place, description).Scan(&eventID)
	if err != nil {
		tx.Rollback()
		return eventID, err
	}
	for _, lection := range lections {
		_, err = tx.Exec("INSERT INTO event_lections (id_event, id_lection) VALUES ($1, $2)", eventID, lection)
		if err != nil {
			tx.Rollback()
			return eventID, err
		}
	}
	err = tx.Commit()
	return eventID, err
}

func GetNextEvent() (*Event, error) {
	q := `SELECT e.starttime, e.endtime, e.description, p.name, p.address
	FROM events e
	LEFT JOIN places p ON p.id = e.place 
	WHERE e.starttime>NOW()
	ORDER BY e.id DESC 
	LIMIT 1;`
	row := dbConn.QueryRow(q)
	e := &Event{}
	if err := row.Scan(&e.StartTime, &e.EndTime, &e.Description, &e.PlaceName, &e.Address); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUndefinedNextEvent
		}
		return nil, err
	}
	return e, nil
}
