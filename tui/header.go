package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/luisedmc/ghcmd/domain"
)

var (
	HeaderUserStyle = lipgloss.NewStyle().
			Foreground(MainColor).
			Bold(true)

	HeaderDetailStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#969696", Dark: "#969696"})

	HeaderDividerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#3c3836"))
)

// RenderHeader renders the authenticated user header with a divider.
// The divider width matches the provided gridWidth.
func RenderHeader(user *domain.User, gridWidth int) string {
	login := "?"
	repos := "?"
	if user != nil {
		login = user.Login
		repos = fmt.Sprintf("%d", user.PublicRepos)
	}

	text := HeaderUserStyle.Render("Logged in as "+login) +
		HeaderDetailStyle.Render(fmt.Sprintf(" \u00b7 %s public repos", repos))

	divider := HeaderDividerStyle.Render(strings.Repeat("\u2500", gridWidth))

	return lipgloss.JoinVertical(lipgloss.Center, text, divider)
}
