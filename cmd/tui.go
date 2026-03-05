package cmd

import (
	"outline-cli/internal/api"
	"outline-cli/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func runTUI(client *api.Client) error {
	p := tea.NewProgram(tui.NewModel(client), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
