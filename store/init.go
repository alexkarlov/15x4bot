package store

import (
	"database/sql"
)

var dbConn *sql.DB

func Init(dsn string) error {
	var err error
	dbConn, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	err = dbConn.Ping()
	if err != nil {
		return err
	}
	return nil
}
