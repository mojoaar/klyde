package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mojoaar/klyde/internal/data"
	"github.com/mojoaar/klyde/internal/ui"
)

func main() {
	store, err := data.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "klyde: failed to load data: %v\n", err)
		os.Exit(1)
	}

	app := ui.NewApp(store)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "klyde: %v\n", err)
		os.Exit(1)
	}
}
