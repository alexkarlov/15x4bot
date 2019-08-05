package store

import (
	"database/sql"
	"errors"
	"time"
)

var ErrUndefinedNextRehearsal = errors.New("Next rehearsal is undeffined")

type Rehearsal struct {
	ID        int
	Place     int
	PlaceName string
	Address   string
	MapUrl    string
	Time      time.Time
}

func AddRehearsal(t time.Time, place int) (int, error) {
	var ID int
	err := dbConn.QueryRow("INSERT INTO rehearsals (time, place) VALUES ($1, $2) RETURNING id", t, place).Scan(&ID)
	return ID, err
}

func NextRehearsal() (*Rehearsal, error) {
	q := `SELECT r.time, p.name, p.address, p.map_url 
	FROM rehearsals r 
	LEFT JOIN places p ON p.id = r.place 
	WHERE r.time>NOW()
	ORDER BY r.id DESC 
	LIMIT 1;`
	row := dbConn.QueryRow(q)
	r := &Rehearsal{}
	if err := row.Scan(&r.Time, &r.PlaceName, &r.Address, &r.MapUrl); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUndefinedNextRehearsal
		}
		return nil, err
	}
	return r, nil
}

// Rehearsals returns a list of rehearsals
func Rehearsals() ([]*Rehearsal, error) {
	q := "SELECT r.id, r.time, p.name FROM rehearsals r LEFT JOIN places p ON p.id=r.place"
	rows, err := dbConn.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rehearsals := make([]*Rehearsal, 0)
	for rows.Next() {
		rehearsal := &Rehearsal{}
		if err := rows.Scan(&rehearsal.ID, &rehearsal.Time, &rehearsal.PlaceName); err != nil {
			return nil, err
		}
		rehearsals = append(rehearsals, rehearsal)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rehearsals, err
}

// DeleteRehearsal deletes rehearsal by provided id
func DeleteRehearsal(id int) error {
	q := "DELETE FROM rehearsals WHERE id=$1"
	_, err := dbConn.Exec(q, id)
	return err
}

// LoadRehearsal returns a rehearsal loaded by id
func LoadRehearsal(ID int) (*Rehearsal, error) {
	r := &Rehearsal{}
	q := "SELECT r.id, r.time, p.name, p.address, p.map_url FROM rehearsals r LEFT JOIN places p ON p.id = r.place WHERE r.id=$1"
	err := dbConn.QueryRow(q, ID).Scan(&r.ID, &r.Time, &r.PlaceName, &r.Address, &r.MapUrl)
	if err == sql.ErrNoRows {
		return nil, ErrNoUser
	}
	return r, err
}
