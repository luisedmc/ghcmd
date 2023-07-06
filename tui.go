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

	focusIndex        int
	searchInputs      []textinput.Model
	searchInputsState bool
	createInputs      []textinput.Model
	createInputsState bool
}

type service struct {
	ctx          context.Context
	token        string
	client       *github.Client
	status       bool
	message      string
	errorMessage string
	url          *string
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

		createInputs: tui.CreateInputs(),
	}
}

func (m *model) updateInputs(msg tea.Msg, isSearch bool) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.searchInputs))

	if isSearch {
		for i := range m.searchInputs {
			m.searchInputs[i], cmds[i] = m.searchInputs[i].Update(msg)
		}
	}

	for i := range m.createInputs {
		m.createInputs[i], cmds[i] = m.createInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) tabKey(msg tea.KeyMsg, inputs []textinput.Model, focusIndex int) (tea.Model, tea.Cmd) {
	s := msg.String()

	if s == "enter" && focusIndex == len(inputs) {
		return m, nil
	}

	if s == "up" {
		focusIndex--
	} else {
		focusIndex++
	}

	if focusIndex > len(inputs) {
		focusIndex = 0
	} else if focusIndex < 0 {
		focusIndex = len(inputs)
	}

	cmds := make([]tea.Cmd, len(inputs))
	for i := 0; i <= len(inputs)-1; i++ {
		if i == focusIndex {
			cmds[i] = inputs[i].Focus()
			inputs[i].PromptStyle = tui.FocusedStyle
			inputs[i].TextStyle = tui.FocusedStyle
			continue
		}

		inputs[i].Blur()
		inputs[i].PromptStyle = tui.NoStyle
		inputs[i].TextStyle = tui.NoStyle
	}

	if m.searchInputsState {
		m.searchInputs = inputs
		m.focusIndex = focusIndex
	} else if m.createInputsState {
		m.searchInputs = inputs
		m.focusIndex = focusIndex
	}

	return m, tea.Batch(cmds...)
}

// Init run any intial IO on program start
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// / Update handle IO and commands
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

		case "enter":
			if m.searchInputsState {
				m.searchInputsState = false
				// Perform the search
				m.responseData = SearchRepository(m.service.ctx, m.service.client, m.searchInputs[0].Value(), m.searchInputs[1].Value())
				if m.responseData == nil {
					m.service.errorMessage = "Repository not found!"
					return m, nil
				}
				m.service.url = nil
				m.servicePerformed = true
				return m, nil

			} else if m.createInputsState {
				m.createInputsState = false
				// Perform the creation
				res, msg, err := CreateRepository(m.service.ctx, m.service.client, m.createInputs[0].Value(), m.createInputs[1].Value())
				if err != nil {
					m.service.errorMessage = "Repository already exists."
				}
				m.service.url = res
				m.service.message = msg
				m.responseData = nil
				m.servicePerformed = true
				return m, nil

			} else if !m.tokenInputState {
				if m.service.token == "" {
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
				switch m.list.Cursor {
				// Search Repository
				case 0:
					m.service.errorMessage = ""
					m.searchInputsState = true
					m.searchInputs[0].SetValue("")
					m.searchInputs[1].SetValue("")
					return m, nil

				// Create Repository
				case 1:
					m.service.errorMessage = ""
					m.createInputsState = true
					m.createInputs[0].SetValue("")
					m.createInputs[1].SetValue("")
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

		case "esc":
			if m.searchInputsState || m.createInputsState {
				m.createInputsState = false
				m.searchInputsState = false
				for i := 0; i < 2; i++ {
					m.searchInputs[i].SetValue("")
					m.createInputs[i].SetValue("")
				}
				return m, nil
			}

		case "tab":
			if m.searchInputsState {
				return m.tabKey(msg, m.searchInputs, m.focusIndex)
			} else if m.createInputsState {
				return m.tabKey(msg, m.createInputs, m.focusIndex)
			}
		}
	}

	var cmd tea.Cmd
	if m.searchInputsState {
		cmd = m.updateInputs(msg, true)
	} else {
		cmd = m.updateInputs(msg, false)
	}

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
	if m.tokenInputState {
		sb.WriteString("\n" + m.tokenInput.View() + "\n")
	} else {
		sb.WriteString("\n")
	}

	// Render custom error message
	switch m.service.errorMessage {
	case "There's an error with your Github Token!":
		sb.WriteString(tui.ErrorStyle.Render(m.service.errorMessage, tui.AlertStyle.Render("\nCheck status bar for more details.")) + "\n")
	case "Repository not found!":
		sb.WriteString(tui.ErrorStyle.Render(m.service.errorMessage, tui.AlertStyle.Render("\nThe repository searched was not found!")) + "\n")
	case "Repository already exists!":
		sb.WriteString(tui.ErrorStyle.Render(m.service.errorMessage, tui.AlertStyle.Render("\nYou already have a repository with that name.")) + "\n")
	case "Repository creation failed!":
		sb.WriteString(tui.ErrorStyle.Render(m.service.errorMessage, tui.AlertStyle.Render("\nAn error has occured.")) + "\n")
	}

	// If not typing the token, render the list of services
	if !m.tokenInputState {
		// Render list of services
		sb.WriteString(tui.ListStyle.Render(m.list.View()))
		if m.searchInputsState {
			sb.WriteString("\n" + m.searchInputs[0].View() + "\n" + m.searchInputs[1].View() + "\n")
		} else if m.createInputsState {
			sb.WriteString("\n" + m.createInputs[0].View() + "\n" + m.createInputs[1].View() + "\n")
		}
	}

	// Render service response
	if m.servicePerformed {
		// Search
		if m.responseData != nil {
			sb.WriteString("\nOwner: " + m.responseData.FullName + "\n")
			sb.WriteString("Repository description: " + m.responseData.Description + "\n")
			sb.WriteString("Repository URL: " + m.responseData.URL + "\n")
		}

		// Create
		if m.service.url != nil {
			sb.WriteString("\n" + m.service.message + "\n")
			sb.WriteString("Repository URL: " + *m.service.url + "\n")
		}
	}

	// Update the status bar after user input
	m.statusBar.SetSize(m.statusBarWidth)
	if m.statusText == "" {
		m.statusBar.SetContent("Token Status", fmt.Sprintf("%s %s | %s %s | %s %s | %s %s", m.keys.Up.Help().Key, m.keys.Up.Help().Desc, m.keys.Down.Help().Key, m.keys.Down.Help().Desc, m.keys.Tab.Help().Key, m.keys.Tab.Help().Desc, m.keys.Quit.Help().Key, m.keys.Quit.Help().Desc), "", "")
	} else {
		m.statusBar.SetContent(m.statusText, fmt.Sprintf("%s %s | %s %s | %s %s | %s %s", m.keys.Up.Help().Key, m.keys.Up.Help().Desc, m.keys.Down.Help().Key, m.keys.Down.Help().Desc, m.keys.Tab.Help().Key, m.keys.Tab.Help().Desc, m.keys.Quit.Help().Key, m.keys.Quit.Help().Desc), "", "")
	}

	// Return the final view
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Height(m.height-statusbar.Height).Render(sb.String()),
		m.statusBar.View(),
	)
}
