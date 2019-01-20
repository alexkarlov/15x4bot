package store

import "time"

type Repetition struct {
	Id        int
	Place     int
	PlaceName string
	Address   string
	MapUrl    string
	Time      time.Time
}

func AddRepetition(t time.Time, place int) error {
	_, err := dbConn.Exec("INSERT INTO repetitions (time, place) VALUES ($1, $2)", t, place)
	return err
}

func GetNextRepetition() (*Repetition, error) {
	q := `SELECT r.time, p.name, p.address, p.map_url 
	FROM repetitions r 
	LEFT JOIN places p ON p.id = r.place 
	WHERE r.time>NOW()
	ORDER BY r.id DESC 
	LIMIT 1;`
	row, err := dbConn.Query(q)
	if err != nil {
		return nil, err
	}
	r := &Repetition{}
	row.Next()
	err = row.Scan(&r.Time, &r.PlaceName, &r.Address, &r.MapUrl)
	if err != nil {
		return nil, err
	}
	return r, nil
}
