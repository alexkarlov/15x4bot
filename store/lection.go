package store

import (
	"errors"
	"strconv"
	"strings"
)

var ErrEmptyLectionID = errors.New("empty lection ID")

type Lection struct {
	ID             int
	Name           string
	LectorUsername string
	LectorID       int
}

// TODO: think about checking LectorUsername
func (l *Lection) OwnedBy(u string) bool {
	q := "SELECT l.id FROM lections l LEFT JOIN users u ON l.user_id=u.id WHERE l.id = $1 AND u.username = $2"
	var ID int
	dbConn.QueryRow(q, l.ID, u).Scan(&ID)
	return ID != 0
}

func (l *Lection) AddDescriptionLection(d string) error {
	q := "UPDATE lections SET description=$1 WHERE id=$2"
	_, err := dbConn.Exec(q, d, l.ID)
	if err != nil {
		return err
	}
	return nil
}

func (l *Lection) Lector() (*User, error) {
	u := &User{}
	q := "SELECT u.id, u.username, u.role, u.name FROM lections l LEFT JOIN users u ON u.id=l.user_id WHERE l.id=$1"
	err := dbConn.QueryRow(q, l.ID).Scan(&u.ID, &u.Username, &u.Role, &u.Name)
	return u, err
}

func AddLection(name string, description string, userID int) (int, error) {
	var ID int
	err := dbConn.QueryRow("INSERT INTO lections (name, description, user_id) VALUES ($1, $2, $3) RETURNING id", name, description, userID).Scan(&ID)
	return ID, err
}

func GetLections(newOnly bool) ([]string, error) {
	lections := make([]string, 0)
	typeFilter := ""
	if newOnly {
		typeFilter = "WHERE l.id NOT IN (SELECT id_lection FROM event_lections)"
	}
	baseQuery := "SELECT l.id, l.name, u.name FROM lections l "
	baseQuery += " LEFT JOIN users u ON u.id = user_id " + typeFilter
	rows, err := dbConn.Query(baseQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		lection := &Lection{}
		if err := rows.Scan(&lection.ID, &lection.Name, &lection.LectorUsername); err != nil {
			return nil, err
		}
		lectionText := []string{strconv.Itoa(lection.ID), ". ", lection.Name, ".", lection.LectorUsername}
		lections = append(lections, strings.Join(lectionText, " "))
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lections, err
}

func GetEmptyLections(eventID int) ([]*Lection, error) {
	// TODO: fix filter for empty lections
	q := `SELECT l.name, l.id, u.username, u.id
		FROM lections l 
		LEFT JOIN users u ON u.id=l.user_id 
		WHERE l.id IN (SELECT id_lection FROM event_lections WHERE id_event=$1) AND l.description='-'`
	rows, err := dbConn.Query(q, eventID)
	if err != nil {
		return nil, err
	}
	lections := make([]*Lection, 0)
	for rows.Next() {
		lection := &Lection{}
		err = rows.Scan(&lection.Name, &lection.ID, &lection.LectorUsername, &lection.LectorID)
		lections = append(lections, lection)
	}
	return lections, nil
}
