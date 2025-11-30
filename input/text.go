package input

import (
	"errors"
	"log"
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

/* TextInput Model */
type TextInput struct {
	TextInput *huh.Input
}

// Init implements tea.Model.
func (i *TextInput) Init() tea.Cmd {
	return i.TextInput.Init()
}

func (i *TextInput) View() string {
	return i.TextInput.View()
}

func (i *TextInput) Error() error {
	return i.TextInput.Error()
}

func (i *TextInput) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	model, cmd := i.TextInput.Update(msg)
	if input, ok := model.(*huh.Input); ok {
		i.TextInput = input
	}
	return i, cmd
}

func (i *TextInput) Blur() tea.Cmd {
	return i.TextInput.Blur()
}

func (i *TextInput) Focus() tea.Cmd {
	return i.TextInput.Focus()
}

func (i *TextInput) Value() any {
	return i.TextInput.GetValue()
}

func NewTextInput(id string, prompt string, placeholder string) *TextInput {
	input := huh.NewInput().
		Prompt(prompt).
		Placeholder(placeholder).
		Key(id).
		Title(prompt).
		CharLimit(24).
		Validate(OnlyNumbers)

	return &TextInput{
		TextInput: input,
	}
}

func OnlyNumbers(text string) error {
	matches, err := regexp.Match("[0-9]*.[0-9]{2}", []byte(text))
	if err != nil {
		log.Default().Println("Bad Regex")
	}
	if matches {
		return nil
	} else {
		return errors.New("regex did not match")
	}
}
