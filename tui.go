package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v53/github"
	"github.com/knipferrc/teacup/statusbar"
	"github.com/luisedmc/ghcmd/tui"
)

type model struct {
	height           int
	help             help.Model
	keys             tui.KeyMap
	list             tui.CustomList
	statusText       string
	statusBar        statusbar.Model
	service          service
	servicePerformed bool
	responseData     *Repository
}

type service struct {
	ctx          context.Context
	client       *github.Client
	status       bool
	errorMessage string
}

// StartGHCMD initialize the tui by returning a model
func StartGHCMD() model {
	apiKey, st, status := apiKey()
	ctx := context.Background()
	ts, _ := Token()
	tc := TokenClient(ctx, ts)
	client := GithubClient(tc)

	sb := tui.StatusBar(apiKey)

	l := tui.CustomList{
		Choices: tui.Choices,
	}

	s := service{
		ctx:          ctx,
		client:       client,
		status:       status,
		errorMessage: "",
	}

	return model{
		keys:             tui.KeyMaps(),
		help:             help.New(),
		list:             l,
		statusBar:        sb,
		statusText:       st,
		service:          s,
		responseData:     &Repository{},
		servicePerformed: false,
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
				if !m.service.status {
					m.responseData = nil
					m.servicePerformed = false
					m.service.errorMessage = "There's an error with your Github Token!"
					return m, nil
				}
				m.responseData = SearchRepository(m.service.ctx, m.service.client, "luisedmc", "dsa")
				if m.responseData == nil {
					m.service.errorMessage = "Repository not found!"
					return m, nil
				} else {
					m.servicePerformed = true
					return m, nil
				}
			}
		}
	}

	return m, nil
}

// View return the text UI to be output to the terminal
func (m model) View() string {
	var sb strings.Builder

	// Render main
	sb.WriteString(tui.TitleStyle.Render("Github CMD"))
	sb.WriteRune('\n')
	sb.WriteString("Welcome to Github CMD, a TUI for Github written in Golang.\n")

	// Render custom error message
	switch m.service.errorMessage {
	case "There's an error with your Github Token!":
		sb.WriteString(tui.ErrorStyle.Render(m.service.errorMessage, tui.AlertStyle.Render("\nCheck status bar for more details.")) + "\n")
	case "Repository not found!":
		sb.WriteString(tui.ErrorStyle.Render(m.service.errorMessage, tui.AlertStyle.Render("\nThe repository searched was not found!")) + "\n")
	}

	// Render list of services
	sb.WriteString(tui.ListStyle.Render(m.list.View()))
	if m.servicePerformed {
		sb.WriteString("\n\nResults\n")
		sb.WriteString("Owner: " + m.responseData.FullName + "\n")
		sb.WriteString("Repository description: " + m.responseData.Description + "\n")
		sb.WriteString("Respository URL: " + m.responseData.URL + "\n")
	}

	// Return final view
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Height(m.height-statusbar.Height).Render(sb.String()),
		m.statusBar.View(),
	)
}
