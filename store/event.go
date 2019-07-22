package store

import (
	"database/sql"
	"errors"
	"time"
)

// ErrUndefinedNextEvent needs for determining case where there is no next event
var ErrUndefinedNextEvent = errors.New("Next event is undeffined")

// Event represents an event with general information and lections
type Event struct {
	ID          int
	StartTime   time.Time
	EndTime     time.Time
	PlaceName   string
	Address     string
	Description string
	Letions     []int
}

// AddEvent creates a new event and adds lections to it
func AddEvent(startTime time.Time, endTime time.Time, place int, description string, lections []int) (int, error) {
	tx, err := dbConn.Begin()
	if err != nil {
		return 0, err
	}
	var eventID int
	qEvents := "INSERT INTO events (starttime, endtime, place, description) VALUES ($1, $2, $3, $4) RETURNING id"
	err = tx.QueryRow(qEvents, startTime, endTime, place, description).Scan(&eventID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	for _, lection := range lections {
		_, err = tx.Exec("INSERT INTO event_lections (id_event, id_lection) VALUES ($1, $2)", eventID, lection)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	err = tx.Commit()
	return eventID, err
}

// NextEvent returns next event with info about place name and address
func NextEvent() (*Event, error) {
	q := `SELECT e.starttime, e.endtime, e.description, p.name, p.address
	FROM events e
	LEFT JOIN places p ON p.id = e.place 
	WHERE e.starttime>NOW()
	ORDER BY e.id DESC`
	e := &Event{}
	err := dbConn.QueryRow(q).Scan(&e.StartTime, &e.EndTime, &e.Description, &e.PlaceName, &e.Address)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUndefinedNextEvent
		}
		return nil, err
	}
	return e, nil
}
