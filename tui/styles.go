package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Application general colors
	MainColor = lipgloss.Color("#00FFA2")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Italic(true).
			Align(lipgloss.Center).
			Foreground(MainColor)
	SubtitleStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Foreground(lipgloss.AdaptiveColor{Light: "#969696", Dark: "#969696"})
	ErrorStyle = lipgloss.NewStyle().Inherit(TitleStyle).
			Italic(false).
			Align(lipgloss.Left).
			Foreground(lipgloss.Color("#FF0000"))
	AlertStyle = lipgloss.NewStyle().Inherit(ErrorStyle).
			Foreground(lipgloss.Color("#FFA500"))

	// Status bar colors
	StatusBarForegroundSuccessStyle = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"}
	StatusBarBackgroundSuccessStyle = lipgloss.AdaptiveColor{Light: "#178009", Dark: "#178009"}

	StatusBarForegroundErrorStyle = lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"}
	StatusBarBackgroundErrorStyle = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"}

	DefaultStatusBarStyle = ColorConfig{
		Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
		Background: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"},
	}

	// Menu card styles
	MenuCardSelected = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(MainColor).
				Padding(1, 2).
				Width(35)

	MenuCardUnselected = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#3c3836")).
				Padding(1, 2).
				Width(35)

	MenuCardTitleSelected = lipgloss.NewStyle().
				Bold(true).
				Foreground(MainColor)

	MenuCardTitleUnselected = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#969696"))

	MenuCardDescSelected = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#c8ded6"))

	MenuCardDescUnselected = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#969696"))

	GridStyle = lipgloss.NewStyle().Margin(1, 2)

	// Input colors
	NoStyle      = lipgloss.NewStyle()
	FocusedStyle = lipgloss.NewStyle().Foreground(MainColor)

	// Result card styles
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(MainColor).
			Padding(1, 2).
			MarginTop(1)

	CardHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(MainColor)

	CardOwnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#969696"))

	CardDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	CardLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#969696"))

	CardValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	CardStatStyle = lipgloss.NewStyle().
			Foreground(MainColor).
			Bold(true)

	CardURLStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#969696")).
			Italic(true)

	CardDividerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#3c3836"))
)
