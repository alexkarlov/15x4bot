package store

// Lection represents a 15x4 lection and contains information about a lector
type Lection struct {
	ID          int
	Description string
	Name        string
	Lector      *User
}

// LoadLection loads a lection and lector details
func LoadLection(ID int) (*Lection, error) {
	q := "SELECT l.id, l.name, l.description, u.id, u.username, u.role FROM lections l LEFT JOIN users u ON u.id=l.user_id WHERE l.id=$1"
	l := &Lection{
		Lector: &User{},
	}
	err := dbConn.QueryRow(q, ID).Scan(&l.ID, &l.Name, &l.Description, &l.Lector.ID, &l.Lector.Username, &l.Lector.Role)
	return l, err
}

// AddLectionDescription adds description for the provided lection
func AddLectionDescription(id int, d string) error {
	q := "UPDATE lections SET description=$1 WHERE id=$2"
	_, err := dbConn.Exec(q, d, id)
	if err != nil {
		return err
	}
	return nil
}

// AddLection creates a lection and returns id of created lection
func AddLection(name string, description string, userID int) (int, error) {
	var ID int
	err := dbConn.QueryRow("INSERT INTO lections (name, description, user_id) VALUES ($1, $2, $3) RETURNING id", name, description, userID).Scan(&ID)
	return ID, err
}

// Lections return list of lections. New lections can be useful for creation of event
func Lections(newOnly bool) ([]*Lection, error) {
	typeFilter := ""
	if newOnly {
		typeFilter = "WHERE l.id NOT IN (SELECT id_lection FROM event_lections)"
	}
	baseQuery := "SELECT l.id, l.name, u.name, u.username, u.id, u.role FROM lections l "
	baseQuery += " LEFT JOIN users u ON u.id = user_id " + typeFilter
	rows, err := dbConn.Query(baseQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	lections := make([]*Lection, 0)
	for rows.Next() {
		lection := &Lection{
			Lector: &User{},
		}
		if err := rows.Scan(&lection.ID, &lection.Name, &lection.Lector.Name, &lection.Lector.Username, &lection.Lector.ID, &lection.Lector.Role); err != nil {
			return nil, err
		}
		lections = append(lections, lection)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lections, err
}

// LectionsWithoutDescriptions returns lections without description for provided speaker
func LectionsWithoutDescriptions(u int) ([]*Lection, error) {
	q := `SELECT l.name, l.id, u.username, u.id, u.role
		FROM lections l 
		LEFT JOIN users u ON u.id=l.user_id 
		WHERE u.id=$1 AND l.description=''`
	rows, err := dbConn.Query(q, u)
	if err != nil {
		return nil, err
	}
	lections := make([]*Lection, 0)
	for rows.Next() {
		lection := &Lection{
			Lector: &User{},
		}
		err = rows.Scan(&lection.Name, &lection.ID, &lection.Lector.Username, &lection.Lector.ID, &lection.Lector.Role)
		lections = append(lections, lection)
	}
	return lections, nil
}
