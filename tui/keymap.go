package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Up   key.Binding
	Down key.Binding
	Tab  key.Binding
	Quit key.Binding
}

func KeyMaps() KeyMap {
	k := KeyMap{
		Up:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Tab:  key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next field")),
		Quit: key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "exit")),
	}

	return k
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}
