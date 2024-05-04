package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	logbook            []logbookEntry
	cursor             int
	logbookDisplaySize int
	page               int
}

func initialModel() model {
	return model{}
}

func updatePage(m model) model {
	db, err := dbOpen()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	rows, err := retrieveUniqueCommandNameRows(db, m.logbookDisplaySize, m.logbookDisplaySize*m.page)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	m.logbook = logBookToEntryList(initLogBook(rows))
	print(m.logbook)
	if m.cursor > m.logbookDisplaySize {
		m.cursor = m.logbookDisplaySize - 1
	}
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.logbookDisplaySize = msg.Height - 5
		m = updatePage(m)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < m.logbookDisplaySize-1 {
				m.cursor++
			}
		case "left":
			if m.page > 0 {
				m.page--
				m = updatePage(m)
			}
		case "right":
			m.page++
			m = updatePage(m)
		}

	}

	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("Results %d: \n\n", m.page)
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
