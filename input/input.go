// Package input exists
package input

import tea "github.com/charmbracelet/bubbletea"

type Input interface {
	Focus() tea.Cmd
	Blur() tea.Cmd
	// Value() any
	// View() string
	Update(tea.Msg) (Input, tea.Cmd)
	// Error() error
}
