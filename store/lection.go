package store

import (
	"strconv"
	"strings"
)

type Lection struct {
	ID     int
	Name   string
	Lector string
}

func AddLection(name string, description string, userID int) error {
	_, err := dbConn.Exec("INSERT INTO lections (name, description, user_id) VALUES ($1, $2, $3) RETURNING id", name, description, userID)
	return err
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
		if err := rows.Scan(&lection.ID, &lection.Name, &lection.Lector); err != nil {
			return nil, err
		}
		lectionText := []string{strconv.Itoa(lection.ID), ". ", lection.Name, ".", lection.Lector}
		lections = append(lections, strings.Join(lectionText, " "))
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lections, err
}
