package main

import (
	"database/sql"
	"os"
)

func dbOpen() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", os.Getenv("HOME")+"/.logbook/logbook.sql")
	return db, err
}

func retrieveUniqueCommandNameRows(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT * FROM command GROUP BY command_name ORDER BY id DESC")
	db.Close()
	return rows, err
}
