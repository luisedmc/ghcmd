package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/statusbar"
)

var (
	// Application general colors
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Italic(true).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("#00FFA2"))

	// Status bar colors
	StatusBarForegroundSuccessStyle = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"}
	StatusBarBackgroundSuccessStyle = lipgloss.AdaptiveColor{Light: "#178009", Dark: "#178009"}

	StatusBarForegroundErrorStyle = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"}
	StatusBarBackgroundErrorStyle = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"}

	DefaultSBColors = statusbar.ColorConfig{
		Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
		Background: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"},
	}

	// List colors
	TitleListStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#969696"}).
			Padding(0, 0, 0, 2)
	DescListStyle = TitleListStyle.Copy().Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#969696"})

	TitleListSelected = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"}).
				Foreground(lipgloss.AdaptiveColor{Light: "#00FFA2", Dark: "#00FFA2"}).
				Padding(0, 0, 0, 1)

	DescListSelected = TitleListSelected.Copy().
				Foreground(lipgloss.AdaptiveColor{Light: "#c8ded6", Dark: "#c8ded6"})

	ListStyle = lipgloss.NewStyle().Margin(1, 2)
)
