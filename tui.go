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
	db, err := dbOpen()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	rows, err := retrieveUniqueCommandNameRows(db)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	logbook := logBookToEntryList(initLogBook(rows))

	return model{logbook: logbook}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.logbookDisplaySize = msg.Height - 5
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
				if len(m.logbook)-m.logbookDisplaySize*m.page-1 != m.cursor {
					m.cursor++
				}
			}
		case "left":
			if m.page > 0 {
				m.page--
			}
		case "right":
			if len(m.logbook)-m.logbookDisplaySize*(m.page+1) > 0 {
				m.page++
			}
			if len(m.logbook)-m.logbookDisplaySize*m.page < m.cursor {
				m.cursor = len(m.logbook) - (m.logbookDisplaySize * m.page) - 1
			}
		}

	}

	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("Results %d: \n\n", m.page)
	commandSliceLen := m.logbookDisplaySize
	if len(m.logbook)-m.logbookDisplaySize*m.page < m.logbookDisplaySize {
		commandSliceLen = len(m.logbook) - m.logbookDisplaySize*m.page
	}
	logbookPage := m.logbook[m.logbookDisplaySize*m.page : m.logbookDisplaySize*(m.page+1)]
	for i := 0; i < commandSliceLen; i++ {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, logbookPage[i].commandName)
	}
	for i := 0; i < m.logbookDisplaySize-commandSliceLen; i++ {
		s += "\n"
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
