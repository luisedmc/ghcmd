package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	Choices = []choice{
		{title: "Search Repository", desc: "Search for a user repository from Github"},
		{title: "Create Repository", desc: "Create a repository in your Github account"},
	}
)

type choice struct {
	title string
	desc  string
}

func (c choice) Title() string       { return c.title }
func (c choice) Description() string { return c.desc }
func (c choice) FilterValue() string { return c.title }

type CardGrid struct {
	Choices []choice
	Cursor  int
}

func (g *CardGrid) CursorLeft() {
	if g.Cursor > 0 {
		g.Cursor--
	}
}

func (g *CardGrid) CursorRight() {
	if g.Cursor < len(g.Choices)-1 {
		g.Cursor++
	}
}

// RenderRow returns the raw horizontal row of cards without centering.
func (g *CardGrid) RenderRow() string {
	parts := make([]string, 0, len(g.Choices)*2-1)

	for i, c := range g.Choices {
		cardStyle := MenuCardUnselected
		titleStyle := MenuCardTitleUnselected
		descStyle := MenuCardDescUnselected

		if i == g.Cursor {
			cardStyle = MenuCardSelected
			titleStyle = MenuCardTitleSelected
			descStyle = MenuCardDescSelected
		}

		content := titleStyle.Render(c.title) + "\n" + descStyle.Render(c.desc)
		parts = append(parts, cardStyle.Render(content))

		if i < len(g.Choices)-1 {
			parts = append(parts, "  ")
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

func (g *CardGrid) View(width int) string {
	row := g.RenderRow()
	return lipgloss.PlaceHorizontal(width, lipgloss.Center, row)
}
