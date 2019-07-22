package store

import (
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
	PLACE_TYPE_FOR_REHEARSALS placeType = "for_rehearsals"
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
	rows, err := dbConn.Query("SELECT id, name, address FROM places " + typeFilter)
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
