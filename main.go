package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m, err := StartGHCMD()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting ghcmd: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
