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

func AddRehearsal(t time.Time, place int) error {
	_, err := dbConn.Exec("INSERT INTO rehearsals (time, place) VALUES ($1, $2)", t, place)
	return err
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
