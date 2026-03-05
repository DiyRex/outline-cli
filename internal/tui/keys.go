package tui

import tea "github.com/charmbracelet/bubbletea"

type keyMap struct {
	Up     tea.Key
	Down   tea.Key
	Enter  tea.Key
	Back   tea.Key
	Quit   tea.Key
	Search tea.Key
	Help   tea.Key
}
