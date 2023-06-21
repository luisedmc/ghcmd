package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/statusbar"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Italic(true).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("#00FFA2"))

	statusBarForegroundSuccessStyle = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"}
	statusBarBackgroundSuccessStyle = lipgloss.AdaptiveColor{Light: "#178009", Dark: "#178009"}

	statusBarForegroundErrorStyle = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"}
	statusBarBackgroundErrorStyle = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"}

	statusBarForegroundStyle = lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"}
	statusBarBackgroundStyle = lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"}
)

type model struct {
	height     int
	help       help.Model
	keys       KeyMap
	statusText string
	statusBar  statusbar.Model
}

// StartGHCMD initialize the tui by returning a model
func StartGHCMD() model {
	sb := statusbar.New(
		statusbar.ColorConfig{
			Foreground: statusBarForegroundSuccessStyle,
			Background: statusBarBackgroundSuccessStyle,
		},
		statusbar.ColorConfig{
			Foreground: statusBarForegroundStyle,
			Background: statusBarBackgroundStyle,
		},
		statusbar.ColorConfig{
			Foreground: statusBarForegroundStyle,
			Background: statusBarBackgroundStyle,
		},
		statusbar.ColorConfig{
			Foreground: statusBarForegroundStyle,
			Background: statusBarBackgroundStyle,
		},
	)

	apiKey, st := apiKey()

	if apiKey == "" {
		sb.SetColors(statusbar.ColorConfig{
			Foreground: statusBarForegroundErrorStyle,
			Background: statusBarBackgroundErrorStyle,
		},
			statusbar.ColorConfig{
				Foreground: statusBarForegroundStyle,
				Background: statusBarBackgroundStyle,
			},
			statusbar.ColorConfig{
				Foreground: statusBarForegroundStyle,
				Background: statusBarBackgroundStyle,
			},
			statusbar.ColorConfig{
				Foreground: statusBarForegroundStyle,
				Background: statusBarBackgroundStyle,
			})
	}

	return model{
		keys: KeyMap{
			Up:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
			Down: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
			Quit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "exit")),
		},
		help:       help.New(),
		statusBar:  sb,
		statusText: st,
	}
}

// Init run any intial IO on program start
func (m model) Init() tea.Cmd {
	return nil
}

// Update handle IO and commands
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.statusBar.SetSize(msg.Width)
		m.statusBar.SetContent(m.statusText, fmt.Sprintf("%s %s | %s %s | %s %s", m.keys.Up.Help().Key, m.keys.Up.Help().Desc, m.keys.Down.Help().Key, m.keys.Down.Help().Desc, m.keys.Quit.Help().Key, m.keys.Quit.Help().Desc), "", "")
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

// View return the text UI to be output to the terminal
func (m model) View() string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("Github CMD"))
	sb.WriteRune('\n')
	sb.WriteString("Welcome to Github CMD, a TUI for Github written in Golang.\n")

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Height(m.height-statusbar.Height).Render(sb.String()),
		m.statusBar.View(),
	)
}
