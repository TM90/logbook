package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	logbook []logbookEntry
	cursor  int
}

func initialModel() model {
	db, err := dbOpen()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	rows, err := db.Query("SELECT * FROM command GROUP BY command_name ORDER BY id DESC LIMIT 0,20")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	db.Close()
	logbook := initLogBook(rows)
	entries := []logbookEntry{}
	for logbook.Next() {
		entries = append(entries, logbook.Value())
	}
	return model{
		logbook: entries,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < 20 {
				m.cursor++
			}
		}

	}

	return m, nil
}

func (m model) View() string {
	s := "Results: \n\n"
	for i, entry := range m.logbook {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, entry.commandName)
	}
	s += "\nPress q to quit\n"
	return s
}

func runTuiBrowser() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error!")
		os.Exit(1)
	}
}
