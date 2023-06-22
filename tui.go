package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v53/github"
	"github.com/knipferrc/teacup/statusbar"
	"github.com/luisedmc/ghcmd/tui"
)

type model struct {
	height     int
	help       help.Model
	keys       tui.KeyMap
	list       tui.CustomList
	statusText string
	statusBar  statusbar.Model
	service    service
}

type service struct {
	ctx    context.Context
	client *github.Client
	user   string
	repo   string
}

// StartGHCMD initialize the tui by returning a model
func StartGHCMD() model {
	apiKey, st := apiKey()
	ctx := context.Background()
	ts, _ := Token()
	tc := TokenClient(ctx, ts)
	client := GithubClient(tc)

	var sb statusbar.Model

	if apiKey == "" {
		sb = statusbar.New(
			statusbar.ColorConfig{
				Foreground: tui.StatusBarForegroundErrorStyle,
				Background: tui.StatusBarBackgroundErrorStyle,
			},
			tui.DefaultSBColors,
			tui.DefaultSBColors,
			tui.DefaultSBColors,
		)
	} else {
		sb = statusbar.New(
			statusbar.ColorConfig{
				Foreground: tui.StatusBarForegroundSuccessStyle,
				Background: tui.StatusBarBackgroundSuccessStyle,
			},
			tui.DefaultSBColors,
			tui.DefaultSBColors,
			tui.DefaultSBColors,
		)
	}

	l := tui.CustomList{
		Choices: tui.Choices,
	}

	s := service{
		ctx:    ctx,
		client: client,
	}

	return model{
		keys: tui.KeyMap{
			Up:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
			Down: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
			Quit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "exit")),
		},
		help:       help.New(),
		list:       l,
		statusBar:  sb,
		statusText: st,
		service:    s,
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

		case "up", "k":
			m.list.CursorUp()

		case "down", "j":
			m.list.CursorDown()

		// Confirm selection
		case "enter":
			switch m.list.Cursor {
			case 0:
				SearchRepository(m.service.ctx, m.service.client, "luisedmc", "dsa")
				return m, nil
			}
		}
	}

	return m, nil
}

// View return the text UI to be output to the terminal
func (m model) View() string {
	var sb strings.Builder

	sb.WriteString(tui.TitleStyle.Render("Github CMD"))
	sb.WriteRune('\n')
	sb.WriteString("Welcome to Github CMD, a TUI for Github written in Golang.\n")
	sb.WriteString(tui.ListStyle.Render(m.list.View()))

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Height(m.height-statusbar.Height).Render(sb.String()),
		m.statusBar.View(),
	)
}
