package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	options          []string
	cursor           int
	selected         map[int]struct{}
	searchTextFields struct {
		userName    string
		repoName    string
		showFields  bool
		activeField int
	}
	createTextFields struct {
		repoName    string
		showFields  bool
		activeField int
	}
}

func StartGHCMD() model {
	return model{
		options:  []string{"Search repository", "Create repository"},
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}

		case "enter", " ":
			if m.cursor == 0 {
				m.searchTextFields.showFields = true
			} else if m.cursor == 1 {
				m.createTextFields.showFields = true
			}

		case "esc":
			m.searchTextFields.showFields = false
			m.createTextFields.showFields = false
			m.searchTextFields.activeField = 0
			m.createTextFields.activeField = 0
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "What you want to do?\n\n"

	for i, choice := range m.options {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	if m.searchTextFields.showFields {
		s += "\n"
		s += "Github username: " + m.searchTextFields.userName + "\n"
		s += "Repository name: " + m.searchTextFields.repoName + "\n"
		s += "\n"
		s += "Press enter to search repository.\n"
		s += "Press tab to switch fields.\n"
		s += "Press esc to cancel.\n"
	} else if m.createTextFields.showFields {
		s += "\n"
		s += "Repository name: " + m.createTextFields.repoName + "\n"
		s += "\n"
		s += "Press enter to create repository.\n"
		s += "Press tab to switch fields.\n"
		s += "Press esc to cancel.\n"
	}

	s += "\nPress q to quit.\n"

	return s
}
