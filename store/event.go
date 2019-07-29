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

// Events returns a list of events
func Events() ([]*Event, error) {
	q := "SELECT e.id, e.starttime, e.endtime, p.name, p.address, e.description FROM events e LEFT JOIN places p ON p.id=e.place"
	rows, err := dbConn.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := make([]*Event, 0)
	for rows.Next() {
		event := &Event{}
		if err := rows.Scan(&event.ID, &event.StartTime, &event.EndTime, &event.PlaceName, &event.Address, &event.Description); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, err
}

// DeleteEvent deletes event by provided id
func DeleteEvent(id int) error {
	q := "DELETE FROM events WHERE id=$1"
	_, err := dbConn.Exec(q, id)
	return err
}
