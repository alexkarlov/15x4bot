package store

import (
	"fmt"
	"github.com/alexkarlov/simplelog"
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
type placeType string

type PlaceTypes []placeType

const (
	PLACE_TYPE_FOR_ALL        placeType = "for_all"
	PLACE_TYPE_FOR_REPETITION placeType = "for_repetition"
	PLACE_TYPE_FOR_EVENT      placeType = "for_event"
)

func GetPlaces(t PlaceTypes) ([]string, error) {
	places := make([]string, 0)
	typeFilter := ""
	if len(t) > 0 {
		for _, pl := range t {
			typeFilter += "'" + string(pl) + "'" + ","
		}
		typeFilter = fmt.Sprintf("WHERE type IN (%s)", typeFilter[:len(typeFilter)-1])
	}
	log.Info(typeFilter)
	rows, err := dbConn.Query("SELECT id, name, address FROM places " + typeFilter)
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
