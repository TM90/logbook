package main

import (
	"database/sql"
	"fmt"
	"os"
)

func dbOpen() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", os.Getenv("HOME")+"/.logbook/logbook.sql")
	return db, err
}

func retrieveUniqueCommandNameRows(db *sql.DB, limit, offset int) (*sql.Rows, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM command GROUP BY command_name ORDER BY id DESC LIMIT %d OFFSET %d", limit, offset))
	db.Close()
	return rows, err
}
