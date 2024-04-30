package main

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type logbookEntry struct {
	dbId        int
	commandName string
	historyId   int
	exitCode    int
	uuid        string
	execTime    time.Time
}

func logbookRetrieveLastEntry(db *sql.DB) logbookEntry {
	var l logbookEntry
	row := db.QueryRow("SELECT * FROM command ORDER BY ID DESC LIMIT 1")
	err := row.Scan(&l.dbId, &l.commandName, &l.historyId, &l.exitCode, &l.uuid, &l.execTime)
	if err != nil {
		l.historyId = -1
		l.commandName = ""
	}
	return l
}

func logbookEntryFromRow(rows *sql.Rows) logbookEntry {
	var l logbookEntry
	rows.Scan(&l.dbId, &l.commandName, &l.historyId, &l.exitCode, &l.uuid, &l.execTime)
	return l
}

type logBook struct {
	rows *sql.Rows
}

func (l logBook) Next() bool {
	return l.rows.Next()
}

func (l logBook) Value() logbookEntry {
	return logbookEntryFromRow(l.rows)
}

func initLogBook(rows *sql.Rows) logBook {
	return logBook{rows}
}

func (l logBook) String() string {
	return l.Value().commandName
}
