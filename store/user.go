package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrNoUser     = errors.New("no users found")
	ErrUserExists = errors.New("user with this username is aleready exists!")
)

// UserRole describes what user can do
type UserRole string

const (
	// admin can do anything
	USER_ROLE_ADMIN UserRole = "admin"
	// lector can create lections and ask general questions
	USER_ROLE_LECTOR UserRole = "lector"
	// guest can only ask general questions
	USER_ROLE_GUEST UserRole = "guest"
)

// NewUserRole creates a user role parsed from string
// if there is no matching - guest role will be used
func NewUserRole(r string) UserRole {
	switch UserRole(r) {
	case USER_ROLE_ADMIN:
		return USER_ROLE_ADMIN
	case USER_ROLE_LECTOR:
		return USER_ROLE_LECTOR
	case USER_ROLE_GUEST:
		return USER_ROLE_GUEST
	}
	return USER_ROLE_GUEST
}

// User represents a user entity in DB
type User struct {
	ID        int
	TGUserID  int
	Username  string
	Name      string
	FB        string
	VK        string
	PictureID string
	BDate     time.Time
	Role      UserRole
}

// DoesUserExist returns whether user exists by provided id or no
func DoesUserExist(id int) (bool, error) {
	q := "SELECT id FROM users WHERE id=$1"
	err := dbConn.QueryRow(q, id).Scan(new(int))
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// LoadUserByUsername returns a user loaded by username
func LoadUserByUsername(username string) (*User, error) {
	u := &User{}
	q := "SELECT id, role, tg_id, username, name, fb, vk, picture_id, bdate FROM users WHERE username=$1"
	err := dbConn.QueryRow(q, username).Scan(&u.ID, &u.Role, &u.TGUserID, &u.Username, &u.Name, &u.FB, &u.VK, &u.PictureID, &u.BDate)
	if err == sql.ErrNoRows {
		return nil, ErrNoUser
	}
	return u, err
}

// LoadUserByTGID returns a user loaded by username
func LoadUserByTGID(tgID int) (*User, error) {
	u := &User{}
	q := "SELECT id, role, tg_id, username, name, fb, vk, picture_id, bdate FROM users WHERE tg_id=$1"
	err := dbConn.QueryRow(q, tgID).Scan(&u.ID, &u.Role, &u.TGUserID, &u.Username, &u.Name, &u.FB, &u.VK, &u.PictureID, &u.BDate)
	if err == sql.ErrNoRows {
		return nil, ErrNoUser
	}
	return u, err
}

func UpdateTGIDUser(ID int, TGID int) error {
	_, err := dbConn.Exec("UPDATE users SET tg_id=$1, udate=NOW() WHERE id=$2", TGID, ID)
	return err
}

func UpdateNameUser(ID int, name string) error {
	_, err := dbConn.Exec("UPDATE users SET name=$1, udate=NOW() WHERE id=$2", name, ID)
	return err
}

func UpdateFBUser(ID int, fb string) error {
	_, err := dbConn.Exec("UPDATE users SET fb=$1, udate=NOW() WHERE id=$2", fb, ID)
	return err
}

func UpdateVKUser(ID int, vk string) error {
	_, err := dbConn.Exec("UPDATE users SET vk=$1, udate=NOW() WHERE id=$2", vk, ID)
	return err
}

func UpdateBDateUser(ID int, bdate time.Time) error {
	_, err := dbConn.Exec("UPDATE users SET bdate=$1, udate=NOW() WHERE id=$2", bdate, ID)
	return err
}

func UpdatePictureUser(ID int, picture string) error {
	_, err := dbConn.Exec("UPDATE users SET picture_id=$1, udate=NOW() WHERE id=$2", picture, ID)
	return err
}

func UpdateUser(ID int, username string, role UserRole, name string, fb string, vk string, bdate time.Time) error {
	_, err := dbConn.Exec("UPDATE users SET username=$1, role=$2, name=$3, fb=$4, vk=$5, bdate=$6, udate=NOW() WHERE id=$7", username, role, name, fb, vk, bdate, ID)
	return err
}

// AddUserByAdmin creates a new record in users table by admin via tg command
func AddUserByAdmin(username string, role UserRole, name string, fb string, vk string, bdate time.Time) error {
	tx, err := dbConn.Begin()
	if err != nil {
		return err
	}
	err = tx.QueryRow("SELECT id FROM users WHERE username=$1 FOR UPDATE", username).Scan(new(int))
	if err == nil {
		tx.Rollback()
		return ErrUserExists
	}
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("INSERT INTO users (username, role, name, fb, vk, bdate) VALUES ($1, $2, $3, $4, $5, $6)", username, role, name, fb, vk, bdate)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// AddGuestUser creates a new guest in users table and returns created user
func AddGuestUser(username string, tgID int, name string) (*User, error) {
	u := &User{
		Username: username,
		TGUserID: tgID,
		Role:     USER_ROLE_GUEST,
		Name:     name,
	}
	err := dbConn.QueryRow("INSERT INTO users (username, tg_id, role, name) VALUES ($1, $2, $3, $4) RETURNING id", u.Username, u.TGUserID, u.Role, u.Name).Scan(&u.ID)
	if err != nil {
		return nil, err
	}
	return u, err
}

// Users returns a list of users by particular user roles
func Users(roles []UserRole) ([]*User, error) {
	roleFilters := make([]string, 0)
	roleFilter := ""
	if len(roles) > 0 {
		for _, role := range roles {
			roleFilters = append(roleFilters, "'"+string(role)+"'")
		}
		// here used a plain string instead of prepared statment because roles aren't a 3-rd party data
		roleFilter = fmt.Sprintf("WHERE role IN (%s)", strings.Join(roleFilters, ","))
	}
	q := "SELECT id, tg_id, username, name, role FROM users " + roleFilter + " ORDER BY id"
	rows, err := dbConn.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*User, 0)
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.TGUserID, &user.Username, &user.Name, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, err
}

// DeleteUser deletes user by provided id
func DeleteUser(id int) error {
	q := "DELETE FROM users WHERE id=$1"
	_, err := dbConn.Exec(q, id)
	return err
}

func GuestUser() *User {
	return &User{
		Role: USER_ROLE_GUEST,
	}
}
