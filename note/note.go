// Package noteinput exists
package noteinput

import (
	"tui/input"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type NoteInput struct {
	Title string
	Note  *huh.Note
}

func NewNoteInput(title string) *NoteInput {
	input := NoteInput{
		Note:  huh.NewNote().Title(title).Next(true),
		Title: title,
	}
	return &input
}

// Init implements tea.Model.
func (i *NoteInput) Init() tea.Cmd {
	return i.Note.Init()
}

func (i *NoteInput) View() string {
	return i.Title
}

func (i *NoteInput) Error() error {
	return i.Note.Error()
}

func (i *NoteInput) Update(msg tea.Msg) (input.Input, tea.Cmd) {
	var cmd tea.Cmd
	model, cmd := i.Note.Update(msg)
	if input, ok := model.(*huh.Note); ok {
		i.Note = input
	}
	return i, cmd
}

func (i *NoteInput) Blur() tea.Cmd {
	return i.Note.Blur()
}

func (i *NoteInput) Focus() tea.Cmd {
	return i.Note.Focus()
}

func (i *NoteInput) Value() any {
	return i.Note.GetValue()
}
