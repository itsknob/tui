package main

import (
	"log"
	"strings"

	"tui/input"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type state struct {
	pages       []input.Page
	currentPage int
}

type (
	nextPageMsg struct{}
	prevPageMsg struct{}
	navigateMsg struct {
		target int
	}
)

// Init implements tea.Model.
func (s state) Init() tea.Cmd {
	for idx := range s.pages {
		// Focus the `FocusedInput` for each page
		// s.pages[idx].inputs[s.pages[idx].FocusedInput].Focus()
		s.pages[idx].Form.Init()
	}
	// Ensure current page's current input is last to be Focused
	// s.pages[s.currentPage].inputs[s.pages[s.currentPage].FocusedInput].Focus()
	return s.pages[s.currentPage].Form.GetFocusedField().Focus()
}

// Update implements tea.Model.
func (s state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	page := s.pages[s.currentPage]
	// input := page.inputs[page.FocusedInput]
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return s, tea.Quit
		case tea.KeyUp:
			page.FocusedInput--
			if page.FocusedInput < 1 { // exclude menu page at 0
				page.FocusedInput = 3
			}
			s.pages[s.currentPage] = page
			return s, nil
		case tea.KeyDown:
			page.FocusedInput++
			if page.FocusedInput > 3 {
				page.FocusedInput = 1 // exclude menu page at 0
			}
			s.pages[s.currentPage] = page
			return s, nil

		case tea.KeyEsc:
			return s.Update(navigateMsg{
				target: 0,
			})

		case tea.KeyTab, tea.KeyEnter:
			if page.Title == "__Menu__" {
				if msg.Type == tea.KeyEnter {
					log.Default().Println("Pressed enter on menu page")
					log.Default().Printf("Current Index: %d\n", page.FocusedInput)
					log.Default().Printf("# Inputs: %d\n", 3)
					return s.Update(navigateMsg{
						target: page.FocusedInput,
					})
				}
				return s, cmd
			}
			// page = page.NextInput()
			// s.pages[s.currentPage] = page
			cmd := s.pages[s.currentPage].Form.NextField()
			return s, cmd
		case tea.KeyShiftTab:
			if page.Title == "__Menu__" {
				return s, page.Form.PrevField()
			}
			// page = page.PrevInput()
			// s.pages[s.currentPage] = page
			cmd := s.pages[s.currentPage].Form.PrevField()
			return s, cmd
		case tea.KeyCtrlN:
			page.Form.GetFocusedField().Update(page.Form.GetFocusedField().Blur)
			return s.Update(nextPageMsg{})
		case tea.KeyCtrlP:
			return s.Update(s.PrevPage())
		}
	case navigateMsg:
		s.currentPage = msg.target
		return s.Update(nil)
	case nextPageMsg:
		cmds := []tea.Cmd{}
		cmds = append(cmds, s.pages[s.currentPage].Form.GetFocusedField().Blur())

		s.currentPage++
		if s.currentPage > len(s.pages)-1 {
			s.currentPage = 0
		}
		cmds = append(cmds, s.pages[s.currentPage].Form.GetFocusedField().Focus())
		return s, tea.Batch(cmds...)

	case prevPageMsg:
		cmds := []tea.Cmd{}
		cmds = append(cmds, s.pages[s.currentPage].Form.GetFocusedField().Blur())
		s.currentPage--
		if s.currentPage < 0 {
			s.currentPage = len(s.pages) - 1
		}
		cmds = append(cmds, s.pages[s.currentPage].Form.GetFocusedField().Focus())
		return s, tea.Batch(cmds...)

	case SubmitMessage:
		log.Fatalf("This worked with msg: %+v\n", msg)
	}

	model, cmd := page.Form.Update(msg)
	if model, ok := model.(*huh.Form); ok {
		s.pages[s.currentPage].Form = model
	}

	if s.pages[s.currentPage].Form.State == huh.StateCompleted {
		log.Printf("Form.State == huh.StateCompleted")
		return s, tea.Batch(cmd, SubmitForm, tea.Quit)
	}
	// input, cmd = input.Update(msg)
	// page.inputs[page.FocusedInput] = input
	s.pages[s.currentPage] = page

	return s, cmd
}

type SubmitMessage func() tea.Cmd

func SubmitForm() tea.Msg {
	msg := new(SubmitMessage)
	return msg
}

func (s state) NextPage() tea.Msg {
	s.Update(s.pages[s.currentPage].Form.GetFocusedField().Blur())
	return nextPageMsg{}
}

func (s state) PrevPage() tea.Msg {
	s.Update(s.pages[s.currentPage].Form.GetFocusedField().Blur())
	return prevPageMsg{}
}

// View implements tea.Model.
func (s state) View() string {
	var sb strings.Builder
	page := s.pages[s.currentPage]

	// Menu Page UI
	if page.Title == "__Menu__" {
		for idx, menuOption := range []string{"__Menu__", "Deposit", "Withdrawal", "Balance"} {
			if page.FocusedInput == idx {
				if menuOption == "__Menu__" {
					sb.WriteString(" ")
				}
				sb.WriteString(">")
			} else {
				sb.WriteString(" ")
			}
			sb.WriteString(menuOption)
			sb.WriteString("\n")
		}
		return sb.String()
	}

	if page.Form.State == huh.StateCompleted {
		var (
			form1 string
			form2 string
			form3 string
		)
		sb.WriteString("Submitting with: \n")
		log.Println("Page: " + page.Title)
		switch page.Title {
		case "Deposit":
			form1 = s.pages[s.currentPage].Form.GetString("depositAmount")
			form2 = s.pages[s.currentPage].Form.GetString("depositDescription")
			form3 = s.pages[s.currentPage].Form.GetString("depositUser")
			sb.WriteString("Amount: " + form1 + "\n")
			sb.WriteString("Description: " + form2 + "\n")
			sb.WriteString("User: " + form3 + "\n")
		case "Withdrawal":
			form1 = s.pages[s.currentPage].Form.GetString("withdrawalAmount")
			form2 = s.pages[s.currentPage].Form.GetString("withdrawalDescription")
			form3 = s.pages[s.currentPage].Form.GetString("withdrawalUser")
			sb.WriteString("Amount: " + form1 + "\n")
			sb.WriteString("Description: " + form2 + "\n")
			sb.WriteString("User: " + form3 + "\n")
		case "Balance":
			sb.WriteString("How did you submit the balance page?\n")
		}

		return sb.String()
	}

	// Print out Current Form View
	currentPage := s.pages[s.currentPage]
	return currentPage.Title + "\n" + currentPage.Form.View()
}

func NewState(pages []input.Page) *state {
	return &state{
		pages:       pages,
		currentPage: 0,
	}
}

func main() {
	noteDeposit := input.NewNoteInput("Deposit")
	noteWithdrawal := input.NewNoteInput("Withdrawal")
	noteBalance := input.NewNoteInput("Balance")
	menuPage := input.NewPage("__Menu__", noteDeposit, noteWithdrawal, noteBalance)

	input11 := input.NewTextInput("depositAmount", "Amount ", "$0.00")
	input12 := input.NewTextInput("depositDescription", "Description ", "Reason...")
	input13 := input.NewSelectInput("depositUser", "User ", "One", []string{"One", "Two", "Three"})
	pageDeposit := input.NewPage("Deposit", input11, input12, input13)

	input21 := input.NewTextInput("one", "Amount", "")
	input22 := input.NewTextInput("two", "Input 2 ", "")
	input23 := input.NewSelectInput("user", "Select", "One", []string{"One", "Two", "Three"})
	pageWithdrawal := input.NewPage("Withdrawal", input21, input22, input23)

	input31 := input.NewTextInput("three", "Input 1 ", "")
	input32 := input.NewTextInput("four", "Input 2 ", "")
	pageBalance := input.NewPage("Balance", input31, input32)

	pages := []input.Page{menuPage, pageDeposit, pageWithdrawal, pageBalance}
	model := NewState(pages)
	app := tea.NewProgram(model)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	_, err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
