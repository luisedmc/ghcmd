package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v84/github"

	"github.com/luisedmc/ghcmd/db"
	"github.com/luisedmc/ghcmd/model"
	"github.com/luisedmc/ghcmd/tui"
)

type serviceResultMsg struct {
	responseData *model.Repository
	url          *string
	message      string
	err          error
}

type userFetchedMsg struct {
	user  *model.User
	token string
	err   error
}

type Model struct {
	height int

	help    help.Model
	keys    tui.KeyMap
	grid    tui.CardGrid
	spinner spinner.Model
	loading bool

	statusText     string
	statusBar      tui.StatusBarModel
	statusBarWidth int

	service          service
	servicePerformed bool
	responseData     *model.Repository

	tokenInput      textinput.Model
	tokenInputState bool

	focusIndex        int
	searchInputs      []textinput.Model
	createInputs      []textinput.Model
	searchInputsState bool
	createInputsState bool

	database *db.Database
}

type service struct {
	ctx    context.Context
	cancel context.CancelFunc
	token  string
	client *github.Client
	status bool
	user   *model.User

	tokenStatus string
	message     string
	lastErr     error
	url         *string
}

// StartGHCMD initializes the TUI and returns the model
func StartGHCMD() (Model, error) {
	ctx, cancel := context.WithCancel(context.Background())
	database, err := db.OpenDB()
	if err != nil {
		cancel()
		return Model{}, fmt.Errorf("opening database: %w", err)
	}

	token, err := database.Token()
	if err != nil {
		cancel()
		database.Close()
		return Model{}, fmt.Errorf("retrieving token: %w", err)
	}

	s := service{
		ctx:    ctx,
		cancel: cancel,
		token:  token,
	}

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(tui.MainColor)

	m := Model{
		keys:         tui.KeyMaps(),
		help:         help.New(),
		grid:         tui.CardGrid{Choices: tui.Choices},
		spinner:      sp,
		statusText:   "Valid Token",
		statusBar:    tui.StatusBar(s.token, s.tokenStatus, s.status),
		service:      s,
		searchInputs: tui.SearchInputs(),
		createInputs: tui.CreateInputs(),
		database:     database,
	}

	if token == "" {
		m.tokenInput = tui.TokenInput()
		m.statusText = ""
		m.tokenInputState = true
	}

	return m, nil
}

// Close releases resources held by the Model.
func (m Model) Close() error {
	m.service.cancel()
	return m.database.Close()
}

func (m Model) updateInputs(msg tea.Msg, isSearch bool) tea.Cmd {
	if isSearch {
		cmds := make([]tea.Cmd, len(m.searchInputs))
		for i := range m.searchInputs {
			m.searchInputs[i], cmds[i] = m.searchInputs[i].Update(msg)
		}
		return tea.Batch(cmds...)
	} else {
		cmds := make([]tea.Cmd, len(m.createInputs))
		for i := range m.createInputs {
			m.createInputs[i], cmds[i] = m.createInputs[i].Update(msg)
		}
		return tea.Batch(cmds...)
	}
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

func fetchUserCmd(token string) tea.Cmd {
	return func() tea.Msg {
		_, user, err := FetchToken(token)
		return userFetchedMsg{user: user, token: token, err: err}
	}
}

// Init run any initial IO on program start
func (m Model) Init() tea.Cmd {
	if m.service.token != "" {
		return tea.Batch(textinput.Blink, fetchUserCmd(m.service.token))
	}
	return textinput.Blink
}

// Update handle IO and commands
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.statusBarWidth = msg.Width
		return m, nil

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil

	case userFetchedMsg:
		if msg.err == nil && msg.user != nil {
			m.service.user = msg.user
			m.service.status = true
			m.service.tokenStatus = "Valid Token"
			m.statusText = "Valid Token"
			m.statusBar = tui.StatusBar(m.service.token, m.service.tokenStatus, m.service.status)
		} else if msg.err != nil {
			m.service.token = ""
			m.service.status = false

			switch {
			case errors.Is(msg.err, ErrTokenInvalid):
				m.service.tokenStatus = "Invalid Token"
			case errors.Is(msg.err, ErrTokenForbidden):
				m.service.tokenStatus = "Token lacks permissions"
			case errors.Is(msg.err, ErrTokenRateLimited):
				m.service.tokenStatus = "Rate limited"
			case errors.Is(msg.err, ErrTokenServerError):
				m.service.tokenStatus = "GitHub server error"
			default:
				m.service.tokenStatus = "Error validating token"
			}

			m.statusText = ""
			m.tokenInput = tui.TokenInput()
			m.tokenInputState = true
			m.statusBar = tui.StatusBar("", m.service.tokenStatus, false)
		}
		return m, nil

	case serviceResultMsg:
		m.loading = false
		m.servicePerformed = true
		m.responseData = msg.responseData
		m.service.url = msg.url
		m.service.message = msg.message
		m.service.lastErr = msg.err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.service.cancel()
			return m, tea.Quit

		case "left", "h", "up", "k":
			if !m.createInputsState && !m.searchInputsState {
				m.grid.CursorLeft()
			}

		case "right", "l", "down", "j":
			if !m.createInputsState && !m.searchInputsState {
				m.grid.CursorRight()
			}

		case "enter":
			if m.loading {
				return m, nil
			}

			if m.searchInputsState && m.searchInputs[0].Value() != "" && m.searchInputs[1].Value() != "" {
				m.searchInputsState = false
				m.loading = true
				user := m.searchInputs[0].Value()
				repo := m.searchInputs[1].Value()
				ctx := m.service.ctx
				client := m.service.client
				return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
					data, err := SearchRepository(ctx, client, user, repo)
					msg := serviceResultMsg{responseData: data}
					if err != nil {
						msg.err = err
						switch {
						case errors.Is(err, ErrSearchNotFound):
							msg.message = "Repository not found!"
						default:
							msg.message = "Failed to search repository!"
						}
					}
					return msg
				})

			} else if m.createInputsState && m.createInputs[0].Value() != "" && m.createInputs[1].Value() != "" {
				m.createInputsState = false
				m.loading = true
				repoName := m.createInputs[0].Value()
				isPrivate := m.createInputs[1].Value()
				ctx := m.service.ctx
				client := m.service.client
				return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
					url, err := CreateRepository(ctx, client, repoName, isPrivate)
					msg := serviceResultMsg{url: url}
					if err != nil {
						msg.err = err
						switch {
						case errors.Is(err, ErrInvalidPrivateInput):
							msg.message = "Invalid input!"
						case errors.Is(err, ErrRepoAlreadyExists):
							msg.message = "Repository already exists!"
						case errors.Is(err, ErrRepoUnauthorized):
							msg.message = "Token lacks permission to create repositories!"
						case errors.Is(err, ErrRepoCreateFailed):
							msg.message = "Repository creation failed!"
						}
					} else {
						msg.message = "Repository created successfully!"
					}
					return msg
				})

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
				switch m.grid.Cursor {
				// Search Repository
				case 0:
					m.service.message = ""
					m.service.lastErr = nil
					m.searchInputsState = true
					m.focusIndex = 0
					m.searchInputs[0].SetValue("")
					m.searchInputs[1].SetValue("")
					m.searchInputs[0].Focus()
					m.searchInputs[0].PromptStyle = tui.FocusedStyle
					m.searchInputs[0].TextStyle = tui.FocusedStyle
					m.searchInputs[1].Blur()
					m.searchInputs[1].PromptStyle = tui.NoStyle
					m.searchInputs[1].TextStyle = tui.NoStyle
					return m, nil

				// Create Repository
				case 1:
					m.service.message = ""
					m.service.lastErr = nil
					m.createInputsState = true
					m.focusIndex = 0
					m.createInputs[0].SetValue("")
					m.createInputs[1].SetValue("")
					m.createInputs[0].Focus()
					m.createInputs[0].PromptStyle = tui.FocusedStyle
					m.createInputs[0].TextStyle = tui.FocusedStyle
					m.createInputs[1].Blur()
					m.createInputs[1].PromptStyle = tui.NoStyle
					m.createInputs[1].TextStyle = tui.NoStyle
					return m, nil
				}
			}

			m.tokenInputState = false
			m.tokenInput.Blur()

			token, user, err := FetchToken(m.tokenInput.Value())
			if err == nil {
				if dbErr := m.database.SetToken(token); dbErr != nil {
					m.service.message = "Failed to save token!"
					m.service.lastErr = dbErr
					return m, nil
				}
			}
			m.service.token = token
			m.service.user = user
			m.service.status = err == nil

			switch {
			case errors.Is(err, ErrTokenEmpty):
				m.service.tokenStatus = "Unwritten Token"
			case errors.Is(err, ErrTokenInvalid):
				m.service.tokenStatus = "Invalid Token"
			case errors.Is(err, ErrTokenForbidden):
				m.service.tokenStatus = "Token lacks permissions"
			case errors.Is(err, ErrTokenRateLimited):
				m.service.tokenStatus = "Rate limited"
			case errors.Is(err, ErrTokenServerError):
				m.service.tokenStatus = "GitHub server error"
			case err != nil:
				m.service.tokenStatus = "Error validating token"
			default:
				m.service.tokenStatus = "Valid Token"
			}

			// Updating status bar text
			m.statusBar = tui.StatusBar(m.service.token, m.service.tokenStatus, m.service.status)
			m.statusText = m.service.tokenStatus

		case "esc":
			if m.loading {
				m.service.cancel()
				ctx, cancel := context.WithCancel(context.Background())
				m.service.ctx = ctx
				m.service.cancel = cancel
				m.service.client = nil
				m.loading = false
				return m, nil
			}

			if m.searchInputsState || m.createInputsState {
				m.createInputsState = false
				m.searchInputsState = false
				for i := range len(m.searchInputs) {
					m.searchInputs[i].SetValue("")
				}
				for i := range len(m.createInputs) {
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

	var cmds []tea.Cmd

	if m.tokenInputState {
		var cmd tea.Cmd
		m.tokenInput, cmd = m.tokenInput.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	if m.searchInputsState {
		if cmd := m.updateInputs(msg, true); cmd != nil {
			cmds = append(cmds, cmd)
		}
	} else if m.createInputsState {
		if cmd := m.updateInputs(msg, false); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

//go:embed docs/titleascii.txt
var titleASCII string

func (m Model) statusBarKeys() string {
	helps := []struct {
		key  string
		desc string
	}{
		{key: m.keys.Left.Help().Key, desc: m.keys.Left.Help().Desc},
		{key: m.keys.Right.Help().Key, desc: m.keys.Right.Help().Desc},
		{key: m.keys.Tab.Help().Key, desc: m.keys.Tab.Help().Desc},
		{key: m.keys.Esc.Help().Key, desc: m.keys.Esc.Help().Desc},
		{key: m.keys.Quit.Help().Key, desc: m.keys.Quit.Help().Desc},
	}

	parts := make([]string, 0, len(helps))
	for _, help := range helps {
		parts = append(parts, help.key+" "+help.desc)
	}

	return strings.Join(parts, " | ")
}

// View returns the text UI to be output to the terminal
func (m Model) View() string {
	var sb strings.Builder

	// Render main
	sb.WriteString(tui.TitleStyle.Width(m.statusBarWidth).Render(titleASCII))
	sb.WriteString("\n\n")
	sb.WriteString(tui.SubtitleStyle.Width(m.statusBarWidth).Render("Welcome to Github CMD, a TUI for Github written in Golang."))
	sb.WriteRune('\n')

	// Render token input
	if m.tokenInputState && m.service.token == "" {
		sb.WriteString("\n" + m.tokenInput.View() + "\n")
	} else {
		// Render user header (with "?" placeholders until user data loads)
		gridRow := m.grid.RenderRow()
		gridWidth := lipgloss.Width(gridRow)
		header := tui.RenderHeader(m.service.user, gridWidth)
		sb.WriteString("\n" + lipgloss.PlaceHorizontal(m.statusBarWidth, lipgloss.Center, header))
		// Render list of services
		sb.WriteString("\n" + m.grid.View(m.statusBarWidth) + "\n")
		if m.searchInputsState {
			sb.WriteString("\n" + m.searchInputs[0].View() + "\n" + m.searchInputs[1].View() + "\n")
		} else if m.createInputsState {
			sb.WriteString("\n" + m.createInputs[0].View() + "\n" + m.createInputs[1].View() + "\n")
		}
	}

	// Render loading spinner
	if m.loading {
		sb.WriteString("\n " + m.spinner.View() + " Loading...\n")
	}

	// Render error message
	if m.service.lastErr != nil && m.service.message != "" {
		sb.WriteString(tui.ErrorStyle.Render("\n"+m.service.message) + "\n")
	}

	// Render service response
	if m.servicePerformed {
		// Search
		if m.responseData != nil {
			card := tui.RenderRepoCard(*m.responseData, m.statusBarWidth)
			sb.WriteString(lipgloss.PlaceHorizontal(m.statusBarWidth, lipgloss.Center, card))
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
