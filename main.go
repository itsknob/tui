package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type state struct {
	pages       []page
	currentPage int
}

type page struct {
	title        string
	inputs       []Input
	focusedInput int
}

// Init implements tea.Model.
func (s state) Init() tea.Cmd {
	for idx := range s.pages {
		// Focus the `focusedInput` for each page
		s.pages[idx].inputs[s.pages[idx].focusedInput].Focus()
	}
	// Ensure current page's current input is last to be Focused
	s.pages[s.currentPage].inputs[s.pages[s.currentPage].focusedInput].Focus()
	return textinput.Blink
}

func (page page) NextInput() page {
	page.inputs[page.focusedInput].Blur() // blur input we are leaving
	page.focusedInput++
	if page.focusedInput > len(page.inputs)-1 {
		page.focusedInput = 0
	}
	page.inputs[page.focusedInput].Focus() // focus next input
	return page
}

func (page page) PrevInput() page {
	page.inputs[page.focusedInput].Blur() // blur input we are leaving
	page.focusedInput--
	if page.focusedInput < 0 {
		page.focusedInput = len(page.inputs) - 1
	}
	page.inputs[page.focusedInput].Focus() // focus next input
	return page
}

// Update implements tea.Model.
func (s state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	page := s.pages[s.currentPage]
	input := page.inputs[page.focusedInput]
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return s, tea.Quit
		case tea.KeyTab, tea.KeyEnter:
			page = page.NextInput()
			s.pages[s.currentPage] = page
			return s, textinput.Blink
		case tea.KeyShiftTab:
			page = page.PrevInput()
			s.pages[s.currentPage] = page
			return s, textinput.Blink
		case tea.KeyCtrlN:
			return s.NextPage(), textinput.Blink
		case tea.KeyCtrlP:
			return s.PrevPage(), textinput.Blink
		}
	}

	input, cmd = input.Update(msg)
	page.inputs[page.focusedInput] = input
	s.pages[s.currentPage] = page

	return s, cmd
}

func (s state) NextPage() state {
	// Blur
	curPage := s.pages[s.currentPage]
	curPage.inputs[curPage.focusedInput].Blur()

	s.currentPage++
	if s.currentPage > len(s.pages)-1 {
		s.currentPage = 0
	}

	// Focus
	curPage.inputs[curPage.focusedInput].Focus()
	return s
}

func (s state) PrevPage() state {
	// Blur
	curPage := s.pages[s.currentPage]
	curPage.inputs[curPage.focusedInput].Blur()

	// Next Page
	s.currentPage--
	if s.currentPage < 0 {
		s.currentPage = len(s.pages) - 1
	}

	// Focus
	curPage.inputs[curPage.focusedInput].Focus()
	return s
}

// View implements tea.Model.
func (s state) View() string {
	var sb strings.Builder
	page := s.pages[s.currentPage]

	sb.WriteString(page.title)
	sb.WriteString("\n")
	for _, input := range page.inputs {
		sb.WriteString(input.View())
		sb.WriteString("\n")
	}
	sb.WriteString("\n")
	return sb.String()
}

func NewState(pages []page) *state {
	return &state{
		pages:       pages,
		currentPage: 0,
	}
}

type Input interface {
	Focus() tea.Cmd
	Blur() tea.Cmd
	Value() any
	View() string
	Update(tea.Msg) (Input, tea.Cmd)
}

/*
*

Text
*/
type TextInput struct {
	textInput *huh.Input
}

// Init implements tea.Model.
func (i *TextInput) Init() tea.Cmd {
	return i.textInput.Init()
}

func (i *TextInput) View() string {
	return i.textInput.View()
}

func (i *TextInput) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	model, cmd := i.textInput.Update(msg)
	if input, ok := model.(*huh.Input); ok {
		i.textInput = input
	}
	return i, cmd
}

func (i *TextInput) Blur() tea.Cmd {
	return i.textInput.Blur()
}

// func (i *TextInput) Blur() {
// 	i.textInput.Blur()
// }

func (i *TextInput) Focus() tea.Cmd {
	return i.textInput.Focus()
}

func (i *TextInput) Value() any {
	return i.textInput.GetValue()
}

func NewTextInput(id string, prompt string, placeholder string) *TextInput {
	input := huh.NewInput().
		Prompt(prompt).
		Placeholder(placeholder).
		Key(id)

	return &TextInput{
		textInput: input,
	}
}

/*
*

Select
*/
type SelectInput struct {
	selectInput *huh.Select[string]
}

func (i *SelectInput) View() string {
	return i.selectInput.View()
}

func (i *SelectInput) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	updatedInput, cmd := i.selectInput.Update(msg)
	input, _ := updatedInput.(*huh.Select[string])
	i.selectInput = input
	return i, cmd
}

func (i *SelectInput) Blur() tea.Cmd {
	return i.selectInput.Blur()
}

func (i *SelectInput) Focus() tea.Cmd {
	return i.selectInput.Focus()
}

func (i *SelectInput) Value() any {
	return i.selectInput.GetValue()
}

func NewSelectInput(id string, prompt string, placeholder string, options []string) *SelectInput {
	selectInput := huh.NewSelect[string]().
		Key(id).
		Options(huh.NewOptions(options...)...).
		Title(prompt)

	return &SelectInput{
		selectInput: selectInput,
	}
}

func NewPage(title string, inputs ...Input) page {
	var newInputs []Input
	for _, input := range inputs {
		newInputs = append(newInputs, input)
	}
	return page{
		inputs:       newInputs,
		title:        title,
		focusedInput: 0,
	}
}

// todo: refactor so this returns an interface that can wrap any type of input
// that way we can use mixed inputs for our pages.
//
// start by creating new functions to return different types of inputs
// then work on an interface to wrap them that allows them all to
// retreive their value at a minimum
// func NewInput(prompt string, placeholder string) textinput.Model { }

func NewInput(prompt string, placeholder string) textinput.Model {
	input := textinput.New()
	input.Prompt = prompt
	input.Placeholder = placeholder
	input.Width = 30
	return input
}

// func NewTextInput(id string, prompt string, placeholder string) huh.Input {
// 	return *huh.NewInput().
// 		Prompt(prompt).
// 		Placeholder(placeholder).
// 		Key(id)
// }

func main() {
	input11 := NewTextInput("one", "Input 1 ", "")
	input12 := NewTextInput("two", "Input 2 ", "")

	// todo need interface to wrap input types into single type
	// todo preferrably also for Text (output) instead of just inputs
	input13 := NewSelectInput("user", "Select", "One", []string{"One", "Two", "Three"})

	input21 := NewTextInput("one", "Input 1 ", "")
	input22 := NewTextInput("two", "Input 2 ", "")

	page1 := NewPage("Hello from One", input11, input12, input13)
	page2 := NewPage("Hello from Two", input21, input22)

	pages := []page{page1, page2}
	model := NewState(pages)
	app := tea.NewProgram(model)
	_, err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
