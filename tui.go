package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FFA2"))

	viewStyle = lipgloss.NewStyle().
			Padding(1, 2)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#00ff04"}).
		// Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"}).
		Margin(1, 2)
)

type model struct {
	Help   help.Model
	Keys   KeyMap
	Status int
}

// StartGHCMD initialize the tui by returning a model
func StartGHCMD() model {
	return model{
		Keys: KeyMap{
			Up:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
			Down: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
			Quit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "exit")),
		},
		Help: help.New(),
	}
}

// Init run any intial IO on program start
func (m model) Init() tea.Cmd {
	return nil
}

// Update handle IO and commands
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) statusView() string {
	var (
		status string
		keys   string
	)

	switch m.Status {
	default:
		status = "Ready"
		keys = m.Help.View(m.Keys)
	}

	return statusBarStyle.Render(status, keys)
}

// View return the text UI to be output to the terminal
func (m model) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("Github CMD"))
	sb.WriteRune('\n')
	sb.WriteString("Welcome to Github CMD, a TUI for Github\n")

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, viewStyle.Render(sb.String())),
		m.statusView(),
	)
}
