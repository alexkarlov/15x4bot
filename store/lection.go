package store

type Lection struct {
	ID int
}

func AddLection(name string, description string, userID int) error {
	_, err := dbConn.Exec("INSERT INTO lections (name, description, user_id) VALUES ($1, $2, $3) RETURNING id", name, description, userID)
	return err
}
