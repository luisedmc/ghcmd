package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/statusbar"
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Italic(true).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("#00FFA2"))

	ListStyle = lipgloss.NewStyle().Margin(1, 2)

	StatusBarForegroundSuccessStyle = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"}
	StatusBarBackgroundSuccessStyle = lipgloss.AdaptiveColor{Light: "#178009", Dark: "#178009"}

	StatusBarForegroundErrorStyle = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"}
	StatusBarBackgroundErrorStyle = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"}

	// Default status bar colors
	DefaultSBColors = statusbar.ColorConfig{
		Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
		Background: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"},
	}
)
