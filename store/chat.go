package store

import (
	"database/sql"
)

type Chat struct {
	ID       int
	TGChatID int64
	UserID   int
}

type User struct {
	ID       int
	Username string
}

func ChatUpsert(chat int64, username string) error {
	tx, err := dbConn.Begin()
	if err != nil {
		return err
	}

	user := &User{}
	row := tx.QueryRow("SELECT id FROM users WHERE username=($1) FOR UPDATE", username)

	if err := row.Scan(&user.ID); err != nil {
		if err != sql.ErrNoRows {
			tx.Rollback()
			return err
		}
		// insert new user if it doesn't exist
		err := tx.QueryRow("INSERT INTO users (username) VALUES ($1) RETURNING id", username).Scan(&user.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = tx.Exec("INSERT INTO chats (tg_chat_id, user_id) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET udate=NOW(), tg_chat_id=($1)", chat, user.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return err
}
