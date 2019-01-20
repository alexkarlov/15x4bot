package store

import (
	"strconv"
	"strings"
)

type Place struct {
	Id          int    `db:"id"`
	Name        string `db:"name"`
	Address     string `db:"address"`
	Description string `db:"description"`
	MapUrl      string `db:"map_url"`
}

func GetPlaces() ([]string, error) {
	places := make([]string, 0)
	rows, err := dbConn.Query("SELECT id, name, address FROM places")
	if err != nil {
		return places, err
	}

	defer rows.Close()
	for rows.Next() {
		var place Place
		if err := rows.Scan(&place.Id, &place.Name, &place.Address); err != nil {
			return nil, err
		}
		placeText := []string{strconv.Itoa(place.Id), "-", place.Name, ",", place.Address}
		places = append(places, strings.Join(placeText, " "))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return places, err
}
