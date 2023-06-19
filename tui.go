package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle       = focusedStyle.Copy()
	noStyle           = lipgloss.NewStyle()
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type model struct {
	list         list.Model
	cursor       int
	selected     map[int]struct{}
	focusIndex   int
	searchInputs []textinput.Model
	createInputs []textinput.Model
}

type item string

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func StartGHCMD() model {
	opt := []list.Item{
		item("Search repository"),
		item("Create repository"),
	}

	l := list.New(opt, itemDelegate{}, 20, 14)

	return model{
		list:     l,
		selected: make(map[int]struct{}),
	}
}

func (i item) FilterValue() string { return "" }

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

func (m model) showCreateInputs() model {
	m.createInputs = make([]textinput.Model, 1)

	for i := range m.createInputs {
		t := textinput.New()
		t.Cursor.Style = cursorStyle

		switch i {
		case 0:
			t.Placeholder = "Repository name"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Public or Privat"
		}

		m.createInputs[i] = t
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

		case "enter":
			switch m.cursor {
			case 0:
				return m.showSearchInputs(), nil
			case 1:
				return m.showCreateInputs(), nil
			}

		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = 0
			}

		case "down", "j":
			m.cursor++
			if m.cursor > len(m.list.Items())-1 {
				m.cursor = len(m.list.Items()) - 1
			}

		case "tab":
			if len(m.searchInputs) == 0 {
				return m, nil
			}
			m.focusIndex++

			if m.focusIndex >= len(m.searchInputs) {
				m.focusIndex = 0
			}

			cmds := make([]tea.Cmd, len(m.searchInputs))
			for i := 0; i < len(m.searchInputs); i++ {
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

	for i, input := range m.searchInputs {
		var cmd tea.Cmd
		m.searchInputs[i], cmd = input.Update(msg)
		if cmd != nil {
			return m, cmd
		}
	}

	for i, input := range m.createInputs {
		var cmd tea.Cmd
		m.createInputs[i], cmd = input.Update(msg)
		if cmd != nil {
			return m, cmd
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "What you want to do?\n\n"

	for i, choice := range m.list.Items() {
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
	} else if m.cursor == 1 && len(m.createInputs) > 0 {
		for i := range m.createInputs {
			s += "\n" + m.createInputs[i].View() + "\n"
		}
	}

	s += "\nPress q to quit.\n"

	return s
}
