package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func logbookOpen() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", os.Getenv("HOME")+"/.logbook/logbook.sql")
	return db, err
}

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

func parseHistoryItem(raw_line string) (string, int) {
	trimmedVal := strings.Trim(raw_line, " ")
	values := strings.Split(trimmedVal, " ")
	id, _ := strconv.Atoi(values[0])
	command := strings.Trim(strings.Join(values[1:], " "), " ")
	return command, id
}

func main() {
	_ = flag.NewFlagSet("init", flag.ExitOnError)
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	cmdName := addCmd.String("command", "", "command")
	exitCode := addCmd.Int("exit-code", 0, "exit-code")
	uuid := addCmd.String("uuid", "", "uuid")
	rawQueryCmd := flag.NewFlagSet("raw_query", flag.ExitOnError)
	rawQuery := rawQueryCmd.String("query", "", "query")
	switch os.Args[1] {
	case "init":
		os.MkdirAll(os.Getenv("HOME")+"/.logbook", 0775)
		os.Create(os.Getenv("HOME") + "/.logbook/logbook.sql")
		os.Chmod(os.Getenv("HOME")+"/.logbook/logbook.sql", 0600)
		db, err := logbookOpen()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		_, err = db.Exec("")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		_, err = db.Exec("CREATE TABLE command (id INTEGER PRIMARY KEY AUTOINCREMENT, command_name TEXT, history_id INTEGER, exit_code INTEGER, uuid TEXT, exec_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL)")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		db.Close()
		os.Exit(0)
	case "add":
		db, err := logbookOpen()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		addCmd.Parse(os.Args[2:])
		command, id := parseHistoryItem(*cmdName)
		lastEntry := logbookRetrieveLastEntry(db)
		if !(id == lastEntry.historyId && command == lastEntry.commandName) {
			_, err = db.Exec("INSERT INTO command (command_name, history_id, exit_code, uuid) VALUES (?, ?, ?, ?)", command, id, *exitCode, *uuid)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
		db.Close()

	case "raw_query":
		rawQueryCmd.Parse(os.Args[2:])
		db, err := logbookOpen()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		rows, err := db.Query(*rawQuery)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		logBook := initLogBook(rows)
		for logBook.Next() {
			fmt.Println(logBook)
		}
		db.Close()

	default:
		fmt.Println("Expected a valid subcommand!")
		os.Exit(1)
	}

}
