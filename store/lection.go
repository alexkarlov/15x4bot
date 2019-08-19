package store

// Lecture represents a 15x4 lecture and contains information about a lector
type Lecture struct {
	ID          int
	Description string
	Name        string
	Lector      *User
}

// LoadLecture loads a lecture and lector details
func LoadLecture(ID int) (*Lecture, error) {
	q := "SELECT l.id, l.name, l.description, u.id, u.username, u.tg_id, u.role FROM lectures l LEFT JOIN users u ON u.id=l.user_id WHERE l.id=$1"
	l := &Lecture{
		Lector: &User{},
	}
	err := dbConn.QueryRow(q, ID).Scan(&l.ID, &l.Name, &l.Description, &l.Lector.ID, &l.Lector.Username, &l.Lector.TGUserID, &l.Lector.Role)
	return l, err
}

// AddLectureDescription adds description for the provided lecture
func AddLectureDescription(id int, d string) error {
	q := "UPDATE lectures SET description=$1 WHERE id=$2"
	_, err := dbConn.Exec(q, d, id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateLecture updates an existed lecture
func UpdateLecture(ID int, name string, description string) error {
	_, err := dbConn.Exec("UPDATE lectures set name=$1, description=$2 WHERE id=$3", name, description, ID)
	return err
}

// AddLecture creates a lecture and returns id of created lecture
func AddLecture(name string, description string, userID int) (int, error) {
	var ID int
	err := dbConn.QueryRow("INSERT INTO lectures (name, description, user_id) VALUES ($1, $2, $3) RETURNING id", name, description, userID).Scan(&ID)
	return ID, err
}

// Lectures return list of lectures. New lectures can be useful for creation of event
func Lectures(newOnly bool) ([]*Lecture, error) {
	typeFilter := "WHERE 1=1 "
	if newOnly {
		typeFilter = " AND l.id NOT IN (SELECT id_lecture FROM event_lectures)"
	}
	baseQuery := "SELECT l.id, l.name, l.description, u.name, u.username, u.tg_id, u.id, u.role FROM lectures l "
	baseQuery += " INNER JOIN users u ON u.id = user_id " + typeFilter
	rows, err := dbConn.Query(baseQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	lectures := make([]*Lecture, 0)
	for rows.Next() {
		lecture := &Lecture{
			Lector: &User{},
		}
		if err := rows.Scan(&lecture.ID, &lecture.Name, &lecture.Description, &lecture.Lector.Name, &lecture.Lector.Username, &lecture.Lector.TGUserID, &lecture.Lector.ID, &lecture.Lector.Role); err != nil {
			return nil, err
		}
		lectures = append(lectures, lecture)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lectures, err
}

// DeleteLecture deletes lecture by provided id
func DeleteLecture(id int) error {
	q := "DELETE FROM lectures WHERE id=$1"
	_, err := dbConn.Exec(q, id)
	return err
}
