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

	statusBarForegroundSucessColor = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"}
	statusBarBackgroundSucessColor = lipgloss.AdaptiveColor{Light: "#178009", Dark: "#178009"}

	statusBarForegroundErrorColor = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"}
	statusBarBackgroundErrorColor = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"}

	statusBarForegroundColor = lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"}
	statusBarBackgroundColor = lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"}
)

type model struct {
	height     int
	Help       help.Model
	Keys       KeyMap
	Status     int
	StatusText string
	StatusBar  statusbar.Model
}

// StartGHCMD initialize the tui by returning a model
func StartGHCMD() model {
	sb := statusbar.New(
		statusbar.ColorConfig{
			Foreground: statusBarForegroundSucessColor,
			Background: statusBarBackgroundSucessColor,
		},
		statusbar.ColorConfig{
			Foreground: statusBarForegroundColor,
			Background: statusBarBackgroundColor,
		},
		statusbar.ColorConfig{
			Foreground: statusBarForegroundColor,
			Background: statusBarBackgroundColor,
		},
		statusbar.ColorConfig{
			Foreground: statusBarForegroundColor,
			Background: statusBarBackgroundColor,
		},
	)

	apiKey, statusText := apiKey()

	if apiKey == "" {
		sb.SetColors(statusbar.ColorConfig{
			Foreground: statusBarForegroundErrorColor,
			Background: statusBarBackgroundErrorColor,
		},
			statusbar.ColorConfig{
				Foreground: statusBarForegroundColor,
				Background: statusBarBackgroundColor,
			},
			statusbar.ColorConfig{
				Foreground: statusBarForegroundColor,
				Background: statusBarBackgroundColor,
			},
			statusbar.ColorConfig{
				Foreground: statusBarForegroundColor,
				Background: statusBarBackgroundColor,
			})
	}

	return model{
		Keys: KeyMap{
			Up:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
			Down: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
			Quit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "exit")),
		},
		Help:       help.New(),
		StatusBar:  sb,
		StatusText: statusText,
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
		m.StatusBar.SetSize(msg.Width)
		m.StatusBar.SetContent(m.StatusText, fmt.Sprintf("%s %s | %s %s | %s %s", m.Keys.Up.Help().Key, m.Keys.Up.Help().Desc, m.Keys.Down.Help().Key, m.Keys.Down.Help().Desc, m.Keys.Quit.Help().Key, m.Keys.Quit.Help().Desc), "", "")
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
		m.StatusBar.View(),
	)
}
