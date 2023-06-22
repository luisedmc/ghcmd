package tui

import (
	"strings"
)

var (
	Choices = []choice{
		{title: "Search Repository", desc: "Search for an user repository from Github."},
		{title: "Create Repository", desc: "Create a repository for your Github account."},
	}
)

type choice struct {
	title string
	desc  string
}

func (c choice) Title() string       { return c.title }
func (c choice) Description() string { return c.desc }
func (c choice) FilterValue() string { return c.title }

type CustomList struct {
	Choices []choice
	Cursor  int
}

func (l *CustomList) CursorUp() {
	if l.Cursor > 0 {
		l.Cursor--
	}
}

func (l *CustomList) CursorDown() {
	if l.Cursor < len(l.Choices)-1 {
		l.Cursor++
	}
}

func (l *CustomList) View() string {
	var sb strings.Builder

	for i, choice := range l.Choices {
		styledTitle := TitleListStyle.Render(choice.title)
		styledDesc := DescListStyle.Render(choice.desc)

		if i == l.Cursor {
			styledTitle = TitleListSelected.Render(choice.title)
			styledDesc = DescListSelected.Render(choice.desc)
			sb.WriteString("  " + styledTitle + "\n")
			sb.WriteString("  " + styledDesc + "\n\n")
		} else {
			sb.WriteString("  " + styledTitle + "\n")
			sb.WriteString("  " + styledDesc + "\n\n")
		}
	}

	return sb.String()
}
