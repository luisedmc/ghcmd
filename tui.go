package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	// blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle = focusedStyle.Copy()
	noStyle     = lipgloss.NewStyle()

	// focusedButton = focusedStyle.Copy().Render("[ Search ]")
	// blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type model struct {
	options      []string
	cursor       int
	selected     map[int]struct{}
	focusIndex   int
	searchInputs []textinput.Model
	// cursorMode   cursor.Mode
}

func StartGHCMD() model {
	return model{
		options:  []string{"Search repository", "Create repository"},
		selected: make(map[int]struct{}),
	}
}

func (m model) showSearchInputs() model {
	m.searchInputs = make([]textinput.Model, 2)

	for i := range m.searchInputs {
		t := textinput.New()
		t.Cursor.Style = cursorStyle

		switch i {
		case 0:
			t.Placeholder = "Github username"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Repository name"
		}

		m.searchInputs[i] = t
	}

	return m
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
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}

		case "enter":
			if m.cursor == 0 {
				return m.showSearchInputs(), nil
			}

		case "tab":
			s := msg.String()

			if s == "enter" && m.focusIndex < len(m.searchInputs)-1 {
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.searchInputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.searchInputs)
			}

			cmds := make([]tea.Cmd, len(m.searchInputs))
			for i := 0; i <= len(m.searchInputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.searchInputs[i].Focus()
					m.searchInputs[i].PromptStyle = focusedStyle
					m.searchInputs[i].TextStyle = focusedStyle
					continue
				}
				m.searchInputs[i].Blur()
				m.searchInputs[i].PromptStyle = noStyle
				m.searchInputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.searchInputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.searchInputs {
		m.searchInputs[i], cmds[i] = m.searchInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
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

	if m.cursor == 0 && len(m.searchInputs) > 0 {
		for i := range m.searchInputs {
			s += "\n" + m.searchInputs[i].View() + "\n"
		}
	}

	s += "\nPress q to quit.\n"

	return s
}
