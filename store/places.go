package store

import (
	"strconv"
	"strings"
)

type Place struct {
	Id          int
	Name        string
	Address     string
	Description string
	MapUrl      string
}

func GetPlaces() ([]string, error) {
	places := make([]string, 0)
	rows, err := dbConn.Query("SELECT id, name, address FROM places")
	if err != nil {
		return places, err
	}

	for rows.Next() {
		var place Place
		if err := rows.Scan(&place.Id, &place.Name, &place.Address); err != nil {
			return nil, err
		}
		placeText := []string{strconv.Itoa(place.Id), "-", place.Name, ",", place.Address}
		places = append(places, strings.Join(placeText, " "))
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return places, err
}
