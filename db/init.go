package db

import (
	"database/sql"
)

var dbConn *sql.DB

func Init(dsn string) (err error) {
	dbConn, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	return nil
}
