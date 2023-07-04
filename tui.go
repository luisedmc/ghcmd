package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v53/github"
	"github.com/knipferrc/teacup/statusbar"
	"github.com/luisedmc/ghcmd/tui"
)

type model struct {
	height int

	help             help.Model
	keys             tui.KeyMap
	list             tui.CustomList
	statusText       string
	statusBar        statusbar.Model
	statusBarWidth   int
	service          service
	servicePerformed bool
	responseData     *Repository

	tokenInput textinput.Model
	inputState bool
}

type service struct {
	ctx          context.Context
	client       *github.Client
	status       bool
	errorMessage string
	token        string
}

// StartGHCMD initialize the tui by returning a model
func StartGHCMD() model {
	// apiKey, st, status := apiKey()
	ctx := context.Background()
	ts, _ := Token()
	tc := TokenClient(ctx, ts)
	client := GithubClient(tc)

	l := tui.CustomList{
		Choices: tui.Choices,
	}

	s := service{
		ctx:    ctx,
		client: client,
		// status:       ,
		errorMessage: "",
		token:        "",
	}

	sb := tui.StatusBar(s.token, s.errorMessage, s.status)

	ti := textinput.New()
	ti.Placeholder = "you can paste it :)"
	ti.Focus()
	ti.CharLimit = 40

	return model{
		keys:             tui.KeyMaps(),
		help:             help.New(),
		list:             l,
		statusBar:        sb,
		statusText:       s.errorMessage,
		service:          s,
		responseData:     &Repository{},
		servicePerformed: false,

		tokenInput: ti,
		inputState: true,
	}
}

// Init run any intial IO on program start
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handle IO and commands
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.statusBarWidth = msg.Width
		m.statusBar.SetSize(msg.Width)
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
			if !m.inputState {
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
			} else {
				m.inputState = false
				m.tokenInput.Blur()

				token, em, s := FetchToken(m.tokenInput.Value())
				m.service.token = token
				m.service.errorMessage = em
				m.service.status = s

				// Updating status bar text
				m.statusBar = tui.StatusBar(m.service.token, m.service.errorMessage, m.service.status)
				m.statusText = m.service.errorMessage

			}
		}
	}

	m.tokenInput, _ = m.tokenInput.Update(msg)

	return m, nil
}

// View returns the text UI to be output to the terminal
func (m model) View() string {
	var sb strings.Builder

	// Render main
	sb.WriteString(tui.TitleStyle.Render("Github CMD"))
	sb.WriteRune('\n')
	sb.WriteString("Welcome to Github CMD, a TUI for Github written in Golang.\n")

	// Render token input
	sb.WriteString("\n" + m.tokenInput.View() + "\n\n")

	// Render custom error message
	switch m.service.errorMessage {
	case "There's an error with your Github Token!":
		sb.WriteString(tui.ErrorStyle.Render(m.service.errorMessage, tui.AlertStyle.Render("\nCheck status bar for more details.")) + "\n")
	case "Repository not found!":
		sb.WriteString(tui.ErrorStyle.Render(m.service.errorMessage, tui.AlertStyle.Render("\nThe repository searched was not found!")) + "\n")
	}

	// Not typing...
	if !m.inputState {
		// Render list of services
		sb.WriteString(tui.ListStyle.Render(m.list.View()))
		if m.servicePerformed {
			sb.WriteString("\n\nResults\n")
			sb.WriteString("Owner: " + m.responseData.FullName + "\n")
			sb.WriteString("Repository description: " + m.responseData.Description + "\n")
			sb.WriteString("Repository URL: " + m.responseData.URL + "\n")
		}
	}

	// Update the status bar after user input
	m.statusBar.SetSize(m.statusBarWidth)
	m.statusBar.SetContent(m.statusText, fmt.Sprintf("%s %s | %s %s | %s %s", m.keys.Up.Help().Key, m.keys.Up.Help().Desc, m.keys.Down.Help().Key, m.keys.Down.Help().Desc, m.keys.Quit.Help().Key, m.keys.Quit.Help().Desc), "", "")

	// Return the final view
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Height(m.height-statusbar.Height).Render(sb.String()),
		m.statusBar.View(),
	)
}
