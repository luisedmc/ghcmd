package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v53/github"
	"github.com/knipferrc/teacup/statusbar"
	"github.com/luisedmc/ghcmd/db"
	"github.com/luisedmc/ghcmd/tui"
	"github.com/syndtr/goleveldb/leveldb"
)

type Model struct {
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

	database *leveldb.DB
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
func StartGHCMD() Model {
	ctx := context.Background()
	database, _ := db.OpenDB()
	token, _ := database.GetToken(database.Conn)

	s := service{
		ctx:   ctx,
		token: token,
	}

	l := tui.CustomList{Choices: tui.Choices}
	sb := tui.StatusBar(s.token, s.errorMessage, s.status)

	if token == "" {
		t := tui.TokenInput()
		m := Model{
			keys:            tui.KeyMaps(),
			help:            help.New(),
			list:            l,
			statusBar:       sb,
			statusText:      s.errorMessage,
			service:         s,
			tokenInput:      t,
			tokenInputState: true,
			searchInputs:    tui.SearchInputs(),
			createInputs:    tui.CreateInputs(),
			database:        database.Conn,
		}
		return m
	} else {
		m := Model{
			keys:         tui.KeyMaps(),
			help:         help.New(),
			list:         l,
			statusBar:    sb,
			statusText:   "Valid Token",
			service:      s,
			searchInputs: tui.SearchInputs(),
			createInputs: tui.CreateInputs(),
			database:     database.Conn,
		}
		return m
	}
}

func (m Model) updateInputs(msg tea.Msg, isSearch bool) tea.Cmd {
	cmds := make([]tea.Cmd, 2)

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

func (m Model) tabKey(msg tea.KeyMsg, inputs []textinput.Model, focusIndex int) (tea.Model, tea.Cmd) {
	s := msg.String()

	if s == "up" {
		focusIndex--
		if focusIndex < 0 {
			focusIndex = len(inputs) - 1
		}
	} else {
		focusIndex++
		if focusIndex >= len(inputs) {
			focusIndex = 0
		}
	}

	cmds := make([]tea.Cmd, len(inputs))
	for i := 0; i < len(inputs); i++ {
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

// Init run any initial IO on program start
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handle IO and commands
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.searchInputsState && m.searchInputs[0].Value() != "" && m.searchInputs[1].Value() != "" {
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

			} else if m.createInputsState && m.createInputs[0].Value() != "" && m.createInputs[1].Value() != "" {
				m.createInputsState = false
				// Perform the creation
				res, msg, err := CreateRepository(m.service.ctx, m.service.client, m.createInputs[0].Value(), m.createInputs[1].Value())
				if err != nil {
					m.service.errorMessage = msg
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
			err := m.database.Put([]byte("gh_token"), []byte(token), nil)
			if err != nil {
				m.service.errorMessage = err.Error()
			}
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
			} else {
				return m, nil
			}
		}
	}

	m.tokenInput, _ = m.tokenInput.Update(msg)

	var cmd tea.Cmd
	if m.searchInputsState {
		cmd = m.updateInputs(msg, true)
	} else {
		cmd = m.updateInputs(msg, false)
	}

	return m, cmd
}

func titleASCII() string {
	t, err := os.ReadFile("docs/titleascii.txt")
	if err != nil {
		return "Github CMD"
	}

	return string(t)
}

func (m Model) statusBarKeys() string {
	return fmt.Sprintf("%s %s | %s %s | %s %s | %s %s", m.keys.Up.Help().Key, m.keys.Up.Help().Desc, m.keys.Down.Help().Key, m.keys.Down.Help().Desc, m.keys.Tab.Help().Key, m.keys.Tab.Help().Desc, m.keys.Quit.Help().Key, m.keys.Quit.Help().Desc)
}

// View returns the text UI to be output to the terminal
func (m Model) View() string {
	var sb strings.Builder

	// Render main
	sb.WriteString(tui.TitleStyle.Render(titleASCII()))
	sb.WriteRune('\n')
	sb.WriteString("Welcome to Github CMD, a TUI for Github written in Golang.\n")

	// Render token input
	if m.tokenInputState && m.service.token == "" {
		sb.WriteString("\n" + m.tokenInput.View() + "\n")
	} else {
		sb.WriteString("\n")
		// Render list of services
		sb.WriteString(tui.ListStyle.Render(m.list.View()))
		if m.searchInputsState {
			sb.WriteString("\n" + m.searchInputs[0].View() + "\n" + m.searchInputs[1].View() + "\n")
		} else if m.createInputsState {
			sb.WriteString("\n" + m.createInputs[0].View() + "\n" + m.createInputs[1].View() + "\n")
		}
	}

	// Render custom error message
	switch m.service.errorMessage {
	case "There's an error with your Github Token!":
		sb.WriteString(tui.ErrorStyle.Render("\n"+m.service.errorMessage, tui.AlertStyle.Render("\nCheck status bar for more details.")) + "\n")
	case "Repository not found!":
		sb.WriteString(tui.ErrorStyle.Render("\n"+m.service.errorMessage, tui.AlertStyle.Render("\nThe repository searched was not found!")) + "\n")
	case "Repository already exists!":
		sb.WriteString(tui.ErrorStyle.Render("\n"+m.service.errorMessage, tui.AlertStyle.Render("\nYou already have a repository with that name.")) + "\n")
	case "Repository creation failed!":
		sb.WriteString(tui.ErrorStyle.Render("\n"+m.service.errorMessage, tui.AlertStyle.Render("\nAn error has occurred.")) + "\n")
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
			sb.WriteString("\n " + m.service.message + "\n")
			sb.WriteString("Repository URL: " + *m.service.url + "\n")
		}
	}

	// Update the status bar after user input
	m.statusBar.SetSize(m.statusBarWidth)
	if m.statusText == "" {
		m.statusBar.SetContent("Token Status", m.statusBarKeys(), "", "")
	} else {
		m.statusBar.SetContent(m.statusText, m.statusBarKeys(), "", "")
	}

	// Return the final view
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Height(m.height-statusbar.Height).Render(sb.String()),
		m.statusBar.View(),
	)
}
