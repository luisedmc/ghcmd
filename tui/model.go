package tui

import (
	"context"
	_ "embed"
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/luisedmc/ghcmd/domain"
	"github.com/luisedmc/ghcmd/service"
)

type serviceResultMsg struct {
	responseData *domain.Repository
	url          *string
	message      string
	err          error
}

type userFetchedMsg struct {
	user  *domain.User
	token string
	err   error
}

type Model struct {
	height int

	help    help.Model
	keys    KeyMap
	grid    CardGrid
	spinner spinner.Model
	loading bool

	statusText     string
	statusBar      StatusBarModel
	statusBarWidth int

	ghState          ghState
	servicePerformed bool
	responseData     *domain.Repository

	tokenInput      textinput.Model
	tokenInputState bool

	focusIndex        int
	searchInputs      []textinput.Model
	createInputs      []textinput.Model
	searchInputsState bool
	createInputsState bool

	authSvc        *service.AuthService
	repoSvc        *service.RepoService
	repoSvcFactory func(ctx context.Context, token string) *service.RepoService
}

// ghState holds runtime GitHub session state (context, token, user, UI messages).
type ghState struct {
	ctx    context.Context
	cancel context.CancelFunc
	token  string
	status bool
	user   *domain.User

	tokenStatus string
	message     string
	lastErr     error
	url         *string
}

// RepoSvcFactory is a function that creates a RepoService given a context and
// token. It is called lazily after token validation.
type RepoSvcFactory func(ctx context.Context, token string) *service.RepoService

// NewModel creates a Model wired with the given services. The caller must
// provide a stored token (may be empty) so the TUI can decide whether to
// show the token-input screen.
func NewModel(authSvc *service.AuthService, repoFactory RepoSvcFactory, storedToken string) Model {
	ctx, cancel := context.WithCancel(context.Background())

	gs := ghState{
		ctx:    ctx,
		cancel: cancel,
		token:  storedToken,
	}

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(MainColor)

	m := Model{
		keys:         KeyMaps(),
		help:         help.New(),
		grid:         CardGrid{Choices: Choices},
		spinner:      sp,
		statusText:   "Valid Token",
		statusBar:    StatusBar(gs.token, gs.tokenStatus, gs.status),
		ghState:      gs,
		searchInputs: SearchInputs(),
		createInputs: CreateInputs(),
		authSvc:        authSvc,
		repoSvcFactory: repoFactory,
	}

	if storedToken == "" {
		m.tokenInput = TokenInput()
		m.statusText = ""
		m.tokenInputState = true
	}

	return m
}

// Close releases resources held by the Model.
func (m Model) Close() error {
	m.ghState.cancel()
	return m.authSvc.Close()
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
			inputs[i].PromptStyle = FocusedStyle
			inputs[i].TextStyle = FocusedStyle
			continue
		}

		inputs[i].Blur()
		inputs[i].PromptStyle = NoStyle
		inputs[i].TextStyle = NoStyle
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
	if m.ghState.token != "" {
		return tea.Batch(textinput.Blink, fetchUserCmd(m.authSvc, m.ghState.token))
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
			m.ghState.user = msg.user
			m.ghState.status = true
			m.ghState.tokenStatus = "Valid Token"
			m.statusText = "Valid Token"
			m.statusBar = StatusBar(m.ghState.token, m.ghState.tokenStatus, m.ghState.status)
		} else if msg.err != nil {
			m.ghState.token = ""
			m.ghState.status = false

			switch {
			case errors.Is(msg.err, domain.ErrTokenInvalid):
				m.ghState.tokenStatus = "Invalid Token"
			case errors.Is(msg.err, domain.ErrTokenForbidden):
				m.ghState.tokenStatus = "Token lacks permissions"
			case errors.Is(msg.err, domain.ErrTokenRateLimited):
				m.ghState.tokenStatus = "Rate limited"
			case errors.Is(msg.err, domain.ErrTokenServerError):
				m.ghState.tokenStatus = "GitHub server error"
			default:
				m.ghState.tokenStatus = "Error validating token"
			}

			m.statusText = ""
			m.tokenInput = TokenInput()
			m.tokenInputState = true
			m.statusBar = StatusBar("", m.ghState.tokenStatus, false)
		}
		return m, nil

	case serviceResultMsg:
		m.loading = false
		m.servicePerformed = true
		m.responseData = msg.responseData
		m.ghState.url = msg.url
		m.ghState.message = msg.message
		m.ghState.lastErr = msg.err
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.ghState.cancel()
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
				owner := m.searchInputs[0].Value()
				repo := m.searchInputs[1].Value()
				return m, tea.Batch(m.spinner.Tick, searchRepoCmd(m.repoSvc, m.ghState.ctx, owner, repo))

			} else if m.createInputsState && m.createInputs[0].Value() != "" && m.createInputs[1].Value() != "" {
				m.createInputsState = false
				m.loading = true
				repoName := m.createInputs[0].Value()
				privateInput := m.createInputs[1].Value()

				private := privateInput == "y"
				if privateInput != "y" && privateInput != "n" && privateInput != "" {
					m.loading = false
					m.ghState.message = "Invalid input!"
					m.ghState.lastErr = errors.New("invalid private input: use 'y' or 'n'")
					return m, nil
				}

				return m, tea.Batch(m.spinner.Tick, createRepoCmd(m.repoSvc, m.ghState.ctx, repoName, private))

			} else if !m.tokenInputState {
				if m.ghState.token == "" {
					m.responseData = nil
					m.servicePerformed = false
					m.ghState.message = "There's an error with your Github Token!"
					m.ghState.lastErr = domain.ErrTokenEmpty
					return m, nil
				}
				// Lazily create repo service on first use
				if m.repoSvc == nil {
					m.repoSvc = m.repoSvcFactory(m.ghState.ctx, m.ghState.token)
				}
				switch m.grid.Cursor {
				// Search Repository
				case 0:
					m.ghState.message = ""
					m.ghState.lastErr = nil
					m.searchInputsState = true
					m.focusIndex = 0
					m.searchInputs[0].SetValue("")
					m.searchInputs[1].SetValue("")
					m.searchInputs[0].Focus()
					m.searchInputs[0].PromptStyle = FocusedStyle
					m.searchInputs[0].TextStyle = FocusedStyle
					m.searchInputs[1].Blur()
					m.searchInputs[1].PromptStyle = NoStyle
					m.searchInputs[1].TextStyle = NoStyle
					return m, nil

				// Create Repository
				case 1:
					m.ghState.message = ""
					m.ghState.lastErr = nil
					m.createInputsState = true
					m.focusIndex = 0
					m.createInputs[0].SetValue("")
					m.createInputs[1].SetValue("")
					m.createInputs[0].Focus()
					m.createInputs[0].PromptStyle = FocusedStyle
					m.createInputs[0].TextStyle = FocusedStyle
					m.createInputs[1].Blur()
					m.createInputs[1].PromptStyle = NoStyle
					m.createInputs[1].TextStyle = NoStyle
					return m, nil
				}
			}

			m.tokenInputState = false
			m.tokenInput.Blur()

			user, err := m.authSvc.ValidateAndStore(m.tokenInput.Value())
			if err == nil {
				m.ghState.token = m.tokenInput.Value()
			}
			m.ghState.user = user
			m.ghState.status = err == nil

			switch {
			case errors.Is(err, domain.ErrTokenEmpty):
				m.ghState.tokenStatus = "Unwritten Token"
			case errors.Is(err, domain.ErrTokenInvalid):
				m.ghState.tokenStatus = "Invalid Token"
			case errors.Is(err, domain.ErrTokenForbidden):
				m.ghState.tokenStatus = "Token lacks permissions"
			case errors.Is(err, domain.ErrTokenRateLimited):
				m.ghState.tokenStatus = "Rate limited"
			case errors.Is(err, domain.ErrTokenServerError):
				m.ghState.tokenStatus = "GitHub server error"
			case err != nil:
				m.ghState.tokenStatus = "Error validating token"
			default:
				m.ghState.tokenStatus = "Valid Token"
			}

			// Updating status bar text
			m.statusBar = StatusBar(m.ghState.token, m.ghState.tokenStatus, m.ghState.status)
			m.statusText = m.ghState.tokenStatus

		case "esc":
			if m.loading {
				m.ghState.cancel()
				ctx, cancel := context.WithCancel(context.Background())
				m.ghState.ctx = ctx
				m.ghState.cancel = cancel
				m.repoSvc = nil
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

			if m.responseData != nil || m.servicePerformed || m.ghState.message != "" {
				m.servicePerformed = false
				m.responseData = nil
				m.ghState.url = nil
				m.ghState.message = ""
				m.ghState.lastErr = nil
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
	sb.WriteString(TitleStyle.Width(m.statusBarWidth).Render(titleASCII))
	sb.WriteString("\n\n")
	sb.WriteString(SubtitleStyle.Width(m.statusBarWidth).Render("Welcome to Github CMD, a TUI for Github written in Golang."))
	sb.WriteRune('\n')

	// Render token input
	if m.tokenInputState && m.ghState.token == "" {
		sb.WriteString("\n" + m.tokenInput.View() + "\n")
	} else {
		// Render user header (with "?" placeholders until user data loads)
		gridRow := m.grid.RenderRow()
		gridWidth := lipgloss.Width(gridRow)
		header := RenderHeader(m.ghState.user, gridWidth)
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
	if m.ghState.lastErr != nil && m.ghState.message != "" {
		sb.WriteString(ErrorStyle.Render("\n"+m.ghState.message) + "\n")
	}

	// Render service response
	if m.servicePerformed {
		// Search
		if m.responseData != nil {
			card := RenderRepoCard(*m.responseData, m.statusBarWidth)
			sb.WriteString(lipgloss.PlaceHorizontal(m.statusBarWidth, lipgloss.Center, card))
			sb.WriteString("\n")
		}

		// Create
		if m.ghState.url != nil {
			sb.WriteString("\n " + m.ghState.message + "\n")
			sb.WriteString("Repository URL: " + *m.ghState.url + "\n")
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
		lipgloss.NewStyle().Height(m.height-StatusBarHeight).Render(sb.String()),
		m.statusBar.View(),
	)
}
