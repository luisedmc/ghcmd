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

	help help.Model
	keys tui.KeyMap
	list tui.CustomList

	statusText     string
	statusBar      statusbar.Model
	statusBarWidth int

	service          service
	servicePerformed bool
	responseData     *Repository

	tokenInput      textinput.Model
	tokenInputState bool

	searchInputs      []textinput.Model
	searchInputsState bool

	focusIndex int
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
	ctx := context.Background()

	l := tui.CustomList{
		Choices: tui.Choices,
	}

	s := service{
		ctx:          ctx,
		errorMessage: "",
		token:        "",
	}

	sb := tui.StatusBar(s.token, s.errorMessage, s.status)

	return model{
		keys:             tui.KeyMaps(),
		help:             help.New(),
		list:             l,
		statusBar:        sb,
		statusText:       s.errorMessage,
		service:          s,
		responseData:     &Repository{},
		servicePerformed: false,

		tokenInput:      tui.TokenInput(),
		tokenInputState: true,

		searchInputs: tui.SearchInputs(),
	}
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 2)

	for i := range m.searchInputs {
		m.searchInputs[i], cmds[i] = m.searchInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
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
			// The enter key here is used to select a service from the list only if the token input is not focused, so !m.inputState
			if !m.tokenInputState {
				switch m.list.Cursor {
				case 0:
					// Token error
					if !m.service.status {
						m.responseData = nil
						m.servicePerformed = false
						m.service.errorMessage = "There's an error with your Github Token!"
						return m, nil
					}

					// Create a new client if it doesn't exist
					if m.service.client == nil {
						ts, _ := TokenSource(m.service.token)
						tc := TokenClient(m.service.ctx, ts)
						client := GithubClient(tc)
						m.service.client = client
					}

					// Search for a repository
					m.searchInputsState = true

					m.responseData = SearchRepository(m.service.ctx, m.service.client, m.searchInputs[0].Value(), m.searchInputs[1].Value())
					if m.responseData == nil {
						m.service.errorMessage = "Repository not found!"
						return m, nil
					}

					m.servicePerformed = true
					return m, nil
				}
			}
			m.tokenInputState = false
			m.tokenInput.Blur()

			token, em, s := FetchToken(m.tokenInput.Value())
			m.service.token = token
			m.service.errorMessage = em
			m.service.status = s

			// Updating status bar text
			m.statusBar = tui.StatusBar(m.service.token, m.service.errorMessage, m.service.status)
			m.statusText = m.service.errorMessage

		case "tab":
			if m.searchInputsState {
				s := msg.String()

				if s == "enter" && m.focusIndex == len(m.searchInputs) {
					return m, nil
				}

				if s == "up" {
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
						m.searchInputs[i].PromptStyle = tui.FocusedStyle
						m.searchInputs[i].TextStyle = tui.FocusedStyle
						continue
					}

					m.searchInputs[i].Blur()
					m.searchInputs[i].PromptStyle = tui.NoStyle
					m.searchInputs[i].TextStyle = tui.NoStyle
				}

				return m, tea.Batch(cmds...)
			}
		}
	}

	cmd := m.updateInputs(msg)

	m.tokenInput, _ = m.tokenInput.Update(msg)

	return m, cmd
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

	// If not typing the token, render the list of services
	if !m.tokenInputState {
		// Render list of services
		sb.WriteString(tui.ListStyle.Render(m.list.View()))

		if m.searchInputsState {
			sb.WriteString("\n" + m.searchInputs[0].View() + "\n" + m.searchInputs[1].View() + "\n")
		}

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
