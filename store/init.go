package store

import (
	"database/sql"
	"github.com/alexkarlov/15x4bot/config"
)

var dbConn *sql.DB

var Conf config.DB

func Init() error {
	var err error
	dbConn, err = sql.Open("postgres", Conf.DSN)
	if err != nil {
		return err
	}
	err = dbConn.Ping()
	if err != nil {
		return err
	}
	return nil
}
