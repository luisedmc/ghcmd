package tui

import (
	"github.com/charmbracelet/lipgloss"
)

const StatusBarHeight = 1

type ColorConfig struct {
	Foreground lipgloss.AdaptiveColor
	Background lipgloss.AdaptiveColor
}

type StatusBarModel struct {
	Width              int
	FirstColumn        string
	SecondColumn       string
	ThirdColumn        string
	FourthColumn       string
	FirstColumnColors  ColorConfig
	SecondColumnColors ColorConfig
	ThirdColumnColors  ColorConfig
	FourthColumnColors ColorConfig
}

func NewStatusBar(first, second, third, fourth ColorConfig) StatusBarModel {
	return StatusBarModel{
		FirstColumnColors:  first,
		SecondColumnColors: second,
		ThirdColumnColors:  third,
		FourthColumnColors: fourth,
	}
}

func (m *StatusBarModel) SetSize(width int) {
	m.Width = width
}

func (m *StatusBarModel) SetContent(first, second, third, fourth string) {
	m.FirstColumn = first
	m.SecondColumn = second
	m.ThirdColumn = third
	m.FourthColumn = fourth
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	if max <= 3 {
		return string(runes[:max])
	}
	return string(runes[:max-3]) + "..."
}

func (m StatusBarModel) column(text string, colors ColorConfig) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(colors.Foreground).
		Background(colors.Background).
		Padding(0, 1).
		Height(StatusBarHeight)
}

func (m StatusBarModel) View() string {
	col1 := m.column(m.FirstColumn, m.FirstColumnColors).
		Render(truncate(m.FirstColumn, 30))

	col3 := m.column(m.ThirdColumn, m.ThirdColumnColors).
		Align(lipgloss.Right).
		Render(m.ThirdColumn)

	col4 := m.column(m.FourthColumn, m.FourthColumnColors).
		Align(lipgloss.Right).
		Render(m.FourthColumn)

	remainingWidth := m.Width - lipgloss.Width(col1) - lipgloss.Width(col3) - lipgloss.Width(col4)
	if remainingWidth < 0 {
		remainingWidth = 0
	}

	col2 := m.column(m.SecondColumn, m.SecondColumnColors).
		Width(remainingWidth).
		Render(truncate(m.SecondColumn, remainingWidth))

	return lipgloss.JoinHorizontal(lipgloss.Top, col1, col2, col3, col4)
}

// StatusBar creates a StatusBarModel with colors based on token status.
func StatusBar(apiKey string, statusText string, status bool) StatusBarModel {
	if apiKey == "" {
		return NewStatusBar(
			ColorConfig{
				Foreground: StatusBarForegroundErrorStyle,
				Background: StatusBarBackgroundErrorStyle,
			},
			DefaultStatusBarStyle,
			DefaultStatusBarStyle,
			DefaultStatusBarStyle,
		)
	}

	return NewStatusBar(
		ColorConfig{
			Foreground: StatusBarForegroundSuccessStyle,
			Background: StatusBarBackgroundSuccessStyle,
		},
		DefaultStatusBarStyle,
		DefaultStatusBarStyle,
		DefaultStatusBarStyle,
	)
}
