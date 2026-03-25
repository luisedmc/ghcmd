package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v53/github"
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
	statusBar      tui.StatusBarModel
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
	ctx    context.Context
	token  string
	client *github.Client
	status bool

	tokenStatus string
	message     string
	lastErr     error
	url         *string
}

// StartGHCMD initializes the TUI
func StartGHCMD() Model {
	ctx := context.Background()
	database, err := db.OpenDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	token, err := database.GetToken(database.Conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving token from database: %v\n", err)
		os.Exit(1)
	}

	s := service{
		ctx:   ctx,
		token: token,
	}

	m := Model{
		keys:         tui.KeyMaps(),
		help:         help.New(),
		list:         tui.CustomList{Choices: tui.Choices},
		statusText:   "Valid Token",
		statusBar:    tui.StatusBar(s.token, s.tokenStatus, s.status),
		service:      s,
		searchInputs: tui.SearchInputs(),
		createInputs: tui.CreateInputs(),
		database:     database.Conn,
	}

	if token == "" {
		m.tokenInput = tui.TokenInput()
		m.statusText = ""
		m.tokenInputState = true
	}

	return m
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
	for i := range inputs {
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
		m.createInputs = inputs
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
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if !m.createInputsState && !m.searchInputsState {
				m.list.CursorUp()
			}

		case "down", "j":
			if !m.createInputsState && !m.searchInputsState {
				m.list.CursorDown()
			}

		case "enter":
			if m.searchInputsState && m.searchInputs[0].Value() != "" && m.searchInputs[1].Value() != "" {
				m.searchInputsState = false
				// Perform the search
				responseData, err := SearchRepository(m.service.ctx, m.service.client, m.searchInputs[0].Value(), m.searchInputs[1].Value())
				if err != nil {
					m.service.lastErr = err
					switch {
					case errors.Is(err, ErrSearchNotFound):
						m.service.message = "Repository not found!"
					default:
						m.service.message = "Failed to search repository!"
					}
					return m, nil
				}
				m.responseData = responseData
				m.service.url = nil
				m.service.lastErr = nil
				m.service.message = ""
				m.servicePerformed = true
				return m, nil

			} else if m.createInputsState && m.createInputs[0].Value() != "" && m.createInputs[1].Value() != "" {
				m.createInputsState = false
				// Perform the creation
				res, err := CreateRepository(m.service.ctx, m.service.client, m.createInputs[0].Value(), m.createInputs[1].Value())
				if err != nil {
					m.service.lastErr = err
					switch {
					case errors.Is(err, ErrInvalidPrivateInput):
						m.service.message = "Invalid input!"
					case errors.Is(err, ErrRepoAlreadyExists):
						m.service.message = "Repository already exists!"
					case errors.Is(err, ErrRepoUnauthorized):
						m.service.message = "Token lacks permission to create repositories!"
					case errors.Is(err, ErrRepoCreateFailed):
						m.service.message = "Repository creation failed!"
					}
				} else {
					m.service.lastErr = nil
					m.service.message = "Repository created successfully!"
				}
				m.service.url = res
				m.responseData = nil
				m.servicePerformed = true
				return m, nil

			} else if !m.tokenInputState {
				if m.service.token == "" {
					m.responseData = nil
					m.servicePerformed = false
					m.service.message = "There's an error with your Github Token!"
					m.service.lastErr = ErrTokenEmpty
					return m, nil
				}
				// Create a new client if it doesn't exist
				if m.service.client == nil {
					ts := TokenSource(m.service.token)
					tc := TokenClient(m.service.ctx, ts)
					client := GithubClient(tc)
					m.service.client = client
				}
				switch m.list.Cursor {
				// Search Repository
				case 0:
					m.service.message = ""
					m.service.lastErr = nil
					m.searchInputsState = true
					m.searchInputs[0].SetValue("")
					m.searchInputs[1].SetValue("")
					return m, nil

				// Create Repository
				case 1:
					m.service.message = ""
					m.service.lastErr = nil
					m.createInputsState = true
					m.createInputs[0].SetValue("")
					m.createInputs[1].SetValue("")
					return m, nil
				}
			}

			m.tokenInputState = false
			m.tokenInput.Blur()

			token, err := FetchToken(m.tokenInput.Value())
			if err == nil {
				if dbErr := m.database.Put([]byte("gh_token"), []byte(token), nil); dbErr != nil {
					m.service.message = "Failed to save token!"
					m.service.lastErr = dbErr
					return m, nil
				}
			}
			m.service.token = token
			m.service.status = err == nil

			switch {
			case errors.Is(err, ErrTokenEmpty):
				m.service.tokenStatus = "Unwritten Token"
			case errors.Is(err, ErrTokenTest):
				m.service.tokenStatus = "Error validating token"
			case errors.Is(err, ErrTokenInvalid):
				m.service.tokenStatus = "Invalid Token"
			default:
				m.service.tokenStatus = "Valid Token"
			}

			// Updating status bar text
			m.statusBar = tui.StatusBar(m.service.token, m.service.tokenStatus, m.service.status)
			m.statusText = m.service.tokenStatus

		case "esc":
			if m.searchInputsState || m.createInputsState {
				m.createInputsState = false
				m.searchInputsState = false
				for i := range 2 {
					m.searchInputs[i].SetValue("")
					m.createInputs[i].SetValue("")
				}
				return m, nil
			}

			if m.responseData != nil || m.servicePerformed || m.service.message != "" {
				m.servicePerformed = false
				m.responseData = nil
				m.service.url = nil
				m.service.message = ""
				m.service.lastErr = nil
				return m, tea.ClearScreen
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
	return fmt.Sprintf("%s %s | %s %s | %s %s | %s %s | %s %s", m.keys.Up.Help().Key, m.keys.Up.Help().Desc, m.keys.Down.Help().Key, m.keys.Down.Help().Desc, m.keys.Tab.Help().Key, m.keys.Tab.Help().Desc, m.keys.Esc.Help().Key, m.keys.Esc.Help().Desc, m.keys.Quit.Help().Key, m.keys.Quit.Help().Desc)
}

// View returns the text UI to be output to the terminal
func (m Model) View() string {
	var sb strings.Builder

	// Render main
	sb.WriteString(tui.TitleStyle.Width(m.statusBarWidth).Render(titleASCII()))
	sb.WriteString("\n\n")
	sb.WriteString(tui.SubtitleStyle.Width(m.statusBarWidth).Render("Welcome to Github CMD, a TUI for Github written in Golang."))
	sb.WriteRune('\n')

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

	// Render error message
	if m.service.lastErr != nil && m.service.message != "" {
		sb.WriteString(tui.ErrorStyle.Render("\n"+m.service.message) + "\n")
	}

	// Render service response
	if m.servicePerformed {
		// Search
		if m.responseData != nil {
			cardData := tui.RepositoryCard{
				Name:        m.responseData.Name,
				Owner:       m.responseData.Owner,
				OwnerURL:    m.responseData.OwnerURL,
				Description: m.responseData.Description,
				URL:         m.responseData.URL,
				Stars:       m.responseData.Stars,
				Forks:       m.responseData.Forks,
				Language:    m.responseData.Language,
				OpenIssues:  m.responseData.OpenIssues,
				CreatedAt:   m.responseData.CreatedAt,
				License:     m.responseData.License,
			}
			sb.WriteString(tui.RenderRepoCard(cardData, m.statusBarWidth))
			sb.WriteString("\n")
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
		lipgloss.NewStyle().Height(m.height-tui.StatusBarHeight).Render(sb.String()),
		m.statusBar.View(),
	)
}
