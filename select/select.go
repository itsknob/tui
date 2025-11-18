// Package selectinput
package selectinput

import (
	"log"

	"tui/input"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

/* SelectInput Model */
type SelectInput struct {
	SelectInput *huh.Select[string]
}

func (i *SelectInput) View() string {
	return i.SelectInput.View()
}

func (i *SelectInput) Error() error {
	return i.SelectInput.Error()
}

func (i *SelectInput) Update(msg tea.Msg) (input.Input, tea.Cmd) {
	log.Printf("Updating selectinput. msg: %v\n", msg)
	var cmd tea.Cmd
	updatedInput, cmd := i.SelectInput.Update(msg)
	input, _ := updatedInput.(*huh.Select[string])
	i.SelectInput = input
	return i, cmd
}

func (i *SelectInput) Blur() tea.Cmd {
	return i.SelectInput.Blur()
}

func (i *SelectInput) Focus() tea.Cmd {
	return i.SelectInput.Focus()
}

func (i *SelectInput) Value() any {
	return i.SelectInput.GetValue()
}

func NewSelectInput(id string, prompt string, placeholder string, options []string) *SelectInput {
	selectInput := huh.NewSelect[string]().
		Key(id).
		Options(huh.NewOptions(options...)...).
		Title(prompt)

	return &SelectInput{
		SelectInput: selectInput,
	}
}
