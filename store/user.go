package store

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
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
	ID       int
	Username string
	Name     string
	Role     UserRole
}

// TGChat returns a chat associated with the user
func (u *User) TGChat() (*Chat, error) {
	c := &Chat{}
	q := "SELECT c.id, c.tg_chat_id FROM users u LEFT JOIN chats c ON c.user_id=u.id WHERE u.id=$1"
	err := dbConn.QueryRow(q, u.ID).Scan(&c.ID, &c.TGChatID)
	return c, err
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

// LoadUser returns a user loaded by username
func LoadUser(username string) (*User, error) {
	u := &User{}
	q := "SELECT u.id, u.role, u.username FROM users u WHERE u.username=$1"
	r := ""
	err := dbConn.QueryRow(q, username).Scan(&u.ID, &r, &u.Username)
	u.Role = UserRole(r)
	return u, err
}

// AddUser creates a new record in users table
func AddUser(username string, role UserRole, name string, fb string, vk string, bdate time.Time) error {
	// TODO: check existense of username
	_, err := dbConn.Exec("INSERT INTO users (username, role, name, fb, vk, bdate) VALUES ($1, $2, $3, $4, $5, $6)", username, role, name, fb, vk, bdate)
	return err
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
	q := "SELECT id, username, name FROM users " + roleFilter
	rows, err := dbConn.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*User, 0)
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Name); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, err
}
