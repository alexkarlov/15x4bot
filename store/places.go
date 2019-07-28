package store

import (
	"database/sql"
	"fmt"
	"strings"
)

// Place represents an entity of place for rehearsals, events, etc
type Place struct {
	ID          int
	Name        string
	Address     string
	Description string
}
type placeType string

// PlaceTypes is a slice of place types
type PlaceTypes []placeType

const (
	// this place can be used for rehearsals and events as well
	PLACE_TYPE_FOR_ALL placeType = "for_all"
	// this place can be used only for rehearsals
	PLACE_TYPE_FOR_REHEARSAL placeType = "for_rehearsal"
	// this place can be used only for events
	PLACE_TYPE_FOR_EVENT placeType = "for_event"
)

// Places returns a list of all available places filtered by types
func Places(t PlaceTypes) ([]*Place, error) {
	typeFilter := ""
	typeFilters := make([]string, 0)
	if len(t) > 0 {
		for _, pl := range t {
			typeFilters = append(typeFilters, "'"+string(pl)+"'")
		}
		typeFilter = fmt.Sprintf("WHERE type IN (%s)", strings.Join(typeFilters, ","))
	}
	q := "SELECT id, name, address FROM places " + typeFilter + " ORDER BY id ASC"
	rows, err := dbConn.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	places := make([]*Place, 0)
	for rows.Next() {
		place := &Place{}
		if err := rows.Scan(&place.ID, &place.Name, &place.Address); err != nil {
			return nil, err
		}
		places = append(places, place)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return places, err
}

// DoesPlaceExist returns whether the place exists by provided id or no
func DoesPlaceExist(id int) (bool, error) {
	q := "SELECT id FROM places WHERE id=$1"
	err := dbConn.QueryRow(q, id).Scan(new(int))
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
