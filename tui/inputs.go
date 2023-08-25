package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
)

// TokenInput defines a text input for the token
func TokenInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "your api token (you can paste it)"
	ti.Focus()
	ti.CharLimit = 40

	return ti
}

// SearchInputs defines two text inputs for the username and repository
func SearchInputs() []textinput.Model {
	ui := textinput.New()
	ui.Placeholder = "username"
	ui.Focus()
	ui.CharLimit = 39
	ui.TextStyle = FocusedStyle

	ri := textinput.New()
	ri.Placeholder = "repository"
	ri.CharLimit = 100

	return []textinput.Model{ui, ri}
}

// CreateInputs defines two text inputs for the repository name and if it's private
func CreateInputs() []textinput.Model {
	rn := textinput.New()
	rn.Placeholder = "repository name"
	rn.Focus()
	rn.CharLimit = 100
	rn.TextStyle = FocusedStyle

	p := textinput.New()
	p.Placeholder = "private (y/n)"
	p.CharLimit = 1

	return []textinput.Model{rn, p}
}
