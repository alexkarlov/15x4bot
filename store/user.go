package store

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ErrUserEmptyID = errors.New("empty user ID")

type UserRole string

const (
	USER_ROLE_ADMIN  UserRole = "admin"
	USER_ROLE_LECTOR UserRole = "lector"
	USER_ROLE_GUEST  UserRole = "guest"
)

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

type User struct {
	ID       int
	Username string
	Name     string
	Role     UserRole
}

func (u *User) TGChat() (*Chat, error) {
	if u.ID == 0 {
		return nil, ErrUserEmptyID
	}
	c := &Chat{}
	q := "SELECT c.id, c.tg_chat_id FROM users u LEFT JOIN chats c ON c.user_id=u.id WHERE u.id=$1"
	err := dbConn.QueryRow(q, u.ID).Scan(&c.ID, &c.TGChatID)
	return c, err
}

func AddUser(username string, role UserRole, name string, fb string, vk string, bdate time.Time) error {
	// TODO: check existense of username
	_, err := dbConn.Exec("INSERT INTO users (username, role, name, fb, vk, bdate) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", username, role, name, fb, vk, bdate)
	return err
}

func GetUsers(roles []UserRole) ([]string, error) {
	users := make([]string, 0)
	typeFilter := ""
	if len(roles) > 0 {
		for _, role := range roles {
			typeFilter += "'" + string(role) + "'" + ","
		}
		typeFilter = fmt.Sprintf("WHERE role IN (%s)", typeFilter[:len(typeFilter)-1])
	}
	rows, err := dbConn.Query("SELECT id, username, name FROM users " + typeFilter)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Name); err != nil {
			return nil, err
		}
		userText := []string{strconv.Itoa(user.ID), "-", user.Username, ",", user.Name}
		users = append(users, strings.Join(userText, " "))
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, err
}
