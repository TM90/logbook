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
	uuid        string
	execTime    time.Time
}

func logbookRetrieveLastEntry(db *sql.DB) logbookEntry {
	var result logbookEntry
	row := db.QueryRow("SELECT * FROM command ORDER BY ID DESC LIMIT 1")
	err := row.Scan(&result.dbId, &result.commandName, &result.historyId, &result.uuid, &result.execTime)
	if err != nil {
		result.historyId = -1
		result.commandName = ""
	}
	return result

}

func printLogbookRows(rows *sql.Rows) {
	for rows.Next() {
		var result logbookEntry
		rows.Scan(&result.dbId, &result.commandName, &result.historyId, &result.uuid, &result.execTime)
		fmt.Printf("%s\n", result.commandName)
	}
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
	uuid := addCmd.String("uuid", "", "uuid")
	rawQueryCmd := flag.NewFlagSet("raw_query", flag.ExitOnError)
	rawQuery := rawQueryCmd.String("query", "", "query")
	switch os.Args[1] {
	case "init":
		fmt.Println("Creating sqlite3 db at: " + os.Getenv("HOME") + "/.logbook/logbook.sql")

		os.MkdirAll(os.Getenv("HOME")+"/.logbook", 0775)
		os.Create(os.Getenv("HOME") + "/.logbook/logbook.sql")
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
		_, err = db.Exec("CREATE TABLE command (id INTEGER PRIMARY KEY AUTOINCREMENT, command_name TEXT, history_id INTEGER, uuid TEXT, exec_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL)")
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
			_, err = db.Exec("INSERT INTO command (command_name, history_id, uuid) VALUES (?, ?, ?)", command, id, *uuid)
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
		printLogbookRows(rows)
		db.Close()

	default:
		fmt.Println("Expected a valid subcommand!")
		os.Exit(1)
	}

}
