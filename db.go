package main

import (
	"database/sql"
	"os"
)

func dbOpen() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", os.Getenv("HOME")+"/.logbook/logbook.sql")
	return db, err
}
