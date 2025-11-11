package main

import (
	"log"
	"strings"

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
	form         *huh.Form
}

// Init implements tea.Model.
func (s state) Init() tea.Cmd {
	for idx := range s.pages {
		// Focus the `focusedInput` for each page
		// s.pages[idx].inputs[s.pages[idx].focusedInput].Focus()
		s.pages[idx].form.Init()
	}
	// Ensure current page's current input is last to be Focused
	// s.pages[s.currentPage].inputs[s.pages[s.currentPage].focusedInput].Focus()
	return s.pages[s.currentPage].form.GetFocusedField().Focus()
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
	// input := page.inputs[page.focusedInput]
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return s, tea.Quit
		case tea.KeyTab, tea.KeyEnter:
			// page = page.NextInput()
			// s.pages[s.currentPage] = page
			cmd := s.pages[s.currentPage].form.NextField()
			return s, cmd
		case tea.KeyShiftTab:
			// page = page.PrevInput()
			// s.pages[s.currentPage] = page
			cmd := s.pages[s.currentPage].form.PrevField()
			return s, cmd
		case tea.KeyCtrlN:
			return s.NextPage(), nil
		case tea.KeyCtrlP:
			return s.PrevPage(), nil
		}
	case SubmitMessage:
		log.Fatalf("This worked with msg: %+v\n", msg)
	}

	model, cmd := page.form.Update(msg)
	if model, ok := model.(*huh.Form); ok {
		s.pages[s.currentPage].form = model
	}

	if s.pages[s.currentPage].form.State == huh.StateCompleted {
		log.Printf("form.State == huh.StateCompleted")
		return s, tea.Batch(cmd, SubmitForm, tea.Quit)
	}
	// input, cmd = input.Update(msg)
	// page.inputs[page.focusedInput] = input
	s.pages[s.currentPage] = page

	return s, cmd
}

type SubmitMessage func() tea.Cmd

func SubmitForm() tea.Msg {
	msg := new(SubmitMessage)
	return msg
}

func (s state) NextPage() state {
	// Blur
	// curPage := s.pages[s.currentPage]
	// curPage.inputs[curPage.focusedInput].Blur()

	s.currentPage++
	if s.currentPage > len(s.pages)-1 {
		s.currentPage = 0
	}

	// Focus
	// curPage.inputs[curPage.focusedInput].Focus()
	return s
}

func (s state) PrevPage() state {
	// Blur
	// curPage := s.pages[s.currentPage]
	// curPage.inputs[curPage.focusedInput].Blur()

	// Next Page
	s.currentPage--
	if s.currentPage < 0 {
		s.currentPage = len(s.pages) - 1
	}

	// Focus
	// curPage.inputs[curPage.focusedInput].Focus()
	return s
}

// View implements tea.Model.
func (s state) View() string {
	page := s.pages[s.currentPage]
	if page.form.State == huh.StateCompleted {
		var sb strings.Builder
		sb.WriteString("Submitting with: \n")
		sb.WriteString("   1: " + s.pages[0].form.GetString("one") + "\n")
		sb.WriteString("   2: " + s.pages[0].form.GetString("two") + "\n")
		sb.WriteString("User: " + s.pages[0].form.GetString("user") + "\n")
		sb.WriteString("-----\n")
		sb.WriteString("   3: " + s.pages[1].form.GetString("three") + "\n")
		sb.WriteString("   4: " + s.pages[1].form.GetString("four") + "\n")
		// log.Default().Print("Submitting with: \n")
		// log.Default().Print("1: " + s.pages[0].form.GetString("One") + "\n")
		// log.Default().Print("2: " + s.pages[0].form.GetString("Two") + "\n")
		// log.Default().Print("3: " + s.pages[0].form.GetString("Three") + "\n")
		// log.Default().Print("-----\n\n")
		// log.Default().Print("1: " + s.pages[1].form.GetString("One") + "\n")
		// log.Default().Print("2: " + s.pages[1].form.GetString("Two") + "\n")
		return sb.String()
	}
	currentPage := s.pages[s.currentPage]
	return currentPage.title + "\n" + currentPage.form.View()

	// var sb strings.Builder
	// page := s.pages[s.currentPage]
	//
	// sb.WriteString(page.title)
	// sb.WriteString("\n")
	// for _, input := range page.inputs {
	// 	sb.WriteString(input.View())
	// 	sb.WriteString("\n")
	// }
	// sb.WriteString("\n")
	// return sb.String()
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
	// Value() any
	// View() string
	Update(tea.Msg) (Input, tea.Cmd)
	// Error() error
}

/* TextInput Model */
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

func (i *TextInput) Error() error {
	return i.textInput.Error()
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
		Key(id).
		Title(prompt)

	return &TextInput{
		textInput: input,
	}
}

/* SelectInput Model */
type SelectInput struct {
	selectInput *huh.Select[string]
}

func (i *SelectInput) View() string {
	return i.selectInput.View()
}

func (i *SelectInput) Error() error {
	return i.selectInput.Error()
}

func (i *SelectInput) Update(msg tea.Msg) (Input, tea.Cmd) {
	log.Printf("Updating selectinput. msg: %v\n", msg)
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
	var fields []huh.Field
	for _, i := range inputs {
		switch i := i.(type) {
		case *TextInput:
			fields = append(fields, i.textInput)
		case *SelectInput:
			fields = append(fields, i.selectInput)
		}
	}
	form := huh.NewForm(huh.NewGroup(fields...))
	return page{
		// inputs:       append([]Input{}, inputs...),
		title:        title,
		focusedInput: 0,
		form:         form,
	}
}

func main() {
	input11 := NewTextInput("one", "Input 1 ", "")
	input12 := NewTextInput("two", "Input 2 ", "")
	input13 := NewSelectInput("user", "Select", "One", []string{"One", "Two", "Three"})
	page1 := NewPage("Hello from One", input11, input12, input13)

	input21 := NewTextInput("three", "Input 1 ", "")
	input22 := NewTextInput("four", "Input 2 ", "")
	page2 := NewPage("Hello from Two", input21, input22)

	pages := []page{page1, page2}
	model := NewState(pages)
	app := tea.NewProgram(model)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
