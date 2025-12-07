package main

// package main
//
// import (
// 	"fmt"
// 	"log"
// 	"strconv"
// 	"strings"
//
// 	"tui/input"
//
// 	"github.com/charmbracelet/bubbles/list"
// 	"github.com/charmbracelet/bubbles/table"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/huh"
// )
//
// type activePage int
//
// const (
// 	MenuPage activePage = iota
// 	DepositPage
// 	WithdrawalPage
// 	BalancePage
// 	TransactionsPage
// )
//
// type MainState struct {
// 	activePage      activePage
// 	menuOptionsList list.Model
// }
//
// /** MenuItems **/
// type MenuItem struct {
// 	title       string
// 	description string
// }
//
// func (mi MenuItem) Title() string {
// 	return mi.title
// }
//
// func (mi MenuItem) Description() string {
// 	return mi.description
// }
//
// func (mi MenuItem) FilterValue() string {
// 	return mi.title
// }
//
// func NewMenuItem(title string, description string) MenuItem {
// 	return MenuItem{title: title, description: description}
// }
//
// func (mainState MainState) Init() {
// 	l := new(list.Model)
// 	l.Title = "Menu"
// 	l.SetItems([]list.Item{
// 		MenuItem{title: "Deposit", description: "Make a Deposit"},
// 		MenuItem{title: "Withdrawal", description: "Make a Withdrawal"},
// 		MenuItem{title: "Balance", description: "See your balance"},
// 		MenuItem{title: "Transactions", description: "View your transactions"},
// 	})
// }
//
// type state struct {
// 	pages       []input.Page
// 	currentPage int
// }
//
// type (
// 	nextPageMsg struct{}
// 	prevPageMsg struct{}
// 	navigateMsg struct {
// 		target int
// 	}
// )
//
// // Init implements tea.Model.
// func (s state) Init() tea.Cmd {
// 	for idx := range s.pages {
// 		// Focus the `FocusedInput` for each page
// 		// s.pages[idx].inputs[s.pages[idx].FocusedInput].Focus()
// 		s.pages[idx].Form.Init()
// 	}
// 	// Ensure current page's current input is last to be Focused
// 	// s.pages[s.currentPage].inputs[s.pages[s.currentPage].FocusedInput].Focus()
// 	// return s.pages[s.currentPage].Form.GetFocusedField().Focus()
// 	return nil
// }
//
// // Update implements tea.Model.
// func (s state) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var cmd tea.Cmd
// 	page := s.pages[s.currentPage]
// 	// input := page.inputs[page.FocusedInput]
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.Type {
// 		/** Ensusure you can quit */
// 		case tea.KeyCtrlC:
// 			return s, tea.Quit
//
// 		/**
// 		* These Move through options outside of forms
// 		 */
// 		case tea.KeyUp:
// 			page.FocusedInput--
// 			if page.FocusedInput < 1 { // exclude menu page at 0
// 				page.FocusedInput = 3
// 			}
// 			s.pages[s.currentPage] = page
// 			return s, nil
// 		case tea.KeyDown:
// 			page.FocusedInput++
// 			if page.FocusedInput > 3 {
// 				page.FocusedInput = 1 // exclude menu page at 0
// 			}
// 			s.pages[s.currentPage] = page
// 			return s, nil
//
// 		/** Return to Menu */
// 		case tea.KeyEsc:
// 			return s.Update(navigateMsg{
// 				target: 0,
// 			})
//
// 		case tea.KeyEnter:
// 			// Naviage to Selected Page
// 			if page.Title == "__Menu__" {
// 				return s.Update(navigateMsg{
// 					target: page.FocusedInput,
// 				})
// 			}
// 			cmd := s.pages[s.currentPage].Form.NextField()
// 			return s, cmd
// 		case tea.KeyShiftTab:
// 			if page.Title == "__Menu__" {
// 				return s, page.Form.PrevField()
// 			}
// 			cmd := s.pages[s.currentPage].Form.PrevField()
// 			return s, cmd
//
// 			// case tea.KeyCtrlN:
// 			// 	page.Form.GetFocusedField().Update(page.Form.GetFocusedField().Blur)
// 			// 	return s.Update(nextPageMsg{})
// 			// case tea.KeyCtrlP:
// 			// 	return s.Update(s.PrevPage())
// 		}
//
// 	/** Jump to specified target page */
// 	case navigateMsg:
// 		s.currentPage = msg.target
// 		return s.Update(nil)
// 	case nextPageMsg:
// 		cmds := []tea.Cmd{}
// 		cmds = append(cmds, s.pages[s.currentPage].Form.GetFocusedField().Blur())
//
// 		s.currentPage++
// 		if s.currentPage > len(s.pages)-1 {
// 			s.currentPage = 0
// 		}
// 		cmds = append(cmds, s.pages[s.currentPage].Form.GetFocusedField().Focus())
// 		return s, tea.Batch(cmds...)
//
// 	case prevPageMsg:
// 		cmds := []tea.Cmd{}
// 		cmds = append(cmds, s.pages[s.currentPage].Form.GetFocusedField().Blur())
// 		s.currentPage--
// 		if s.currentPage < 0 {
// 			s.currentPage = len(s.pages) - 1
// 		}
// 		cmds = append(cmds, s.pages[s.currentPage].Form.GetFocusedField().Focus())
// 		return s, tea.Batch(cmds...)
//
// 	case SubmitMessage:
// 		log.Fatalf("This worked with msg: %+v\n", msg)
// 		s.currentPage = 0
// 		return s, nil
// 	}
//
// 	if page.Title == "Balance" {
// 		if page.Form.State == huh.StateCompleted {
// 			log.Default().Printf("Updating Table\n")
// 			model, cmd := page.Table.Update(msg)
// 			page.Table = model
// 			s.pages[s.currentPage] = page
// 			return s, cmd
// 		} else {
// 			return page.Form.Update(msg)
// 		}
// 	}
//
// 	model, cmd := page.Form.Update(msg)
// 	if model, ok := model.(*huh.Form); ok {
// 		s.pages[s.currentPage].Form = model
// 	}
//
// 	if s.pages[s.currentPage].Form.State == huh.StateCompleted {
// 		log.Printf("Form.State == huh.StateCompleted")
// 		return s, tea.Batch(cmd, SubmitForm, tea.Quit)
// 	}
// 	// input, cmd = input.Update(msg)
// 	// page.inputs[page.FocusedInput] = input
// 	s.pages[s.currentPage] = page
//
// 	return s, cmd
// }
//
// type SubmitMessage func() tea.Cmd
//
// func SubmitForm() tea.Msg {
// 	msg := new(SubmitMessage)
// 	return msg
// }
//
// func (s state) NextPage() tea.Msg {
// 	s.Update(s.pages[s.currentPage].Form.GetFocusedField().Blur())
// 	return nextPageMsg{}
// }
//
// func (s state) PrevPage() tea.Msg {
// 	s.Update(s.pages[s.currentPage].Form.GetFocusedField().Blur())
// 	return prevPageMsg{}
// }
//
// // View implements tea.Model.
// func (s state) View() string {
// 	var sb strings.Builder
// 	page := s.pages[s.currentPage]
//
// 	// Menu Page UI
// 	if page.Title == "__Menu__" {
// 		for idx, menuOption := range []string{"__Menu__", "Deposit", "Withdrawal", "Balance"} {
// 			if page.FocusedInput == idx {
// 				if menuOption == "__Menu__" {
// 					sb.WriteString(" ")
// 				}
// 				sb.WriteString(">")
// 			} else {
// 				sb.WriteString(" ")
// 			}
// 			sb.WriteString(menuOption)
// 			sb.WriteString("\n")
// 		}
// 		return sb.String()
// 	}
// 	if page.Title == "Balance" {
// 		if page.Form.State != huh.StateCompleted {
// 			return page.Form.View()
// 		} else {
// 			return page.Table.View()
// 		}
// 		//
// 		// log.Default().Println("Printing Balance Page")
// 		// table33 := table.New()
// 		// table33.SetColumns([]table.Column{
// 		// 	{Width: 16, Title: "Date"},
// 		// 	{Width: 16, Title: "Amount"},
// 		// 	{Width: 16, Title: "Description"},
// 		// })
// 		// // table33.SetRows([]table.Row{{"2025-01-01", "123.45", "First Deposit"}, {"2025-02-02", "234.56", "Second Deposit"}})
// 		// // log.Default().Printf("%s", strings.Join(table33.Rows()[0], ", "))
// 		// table33.FromValues("2025-01-01,123.45,First Deposit\n2025-02-02,234.56,Second Deposit", ",")
// 		// // log.Default().Printf("cols %+v\n", table33.Columns())
// 		// // log.Default().Printf("rows %+v\n", table33.Rows())
// 		// table33.View()
// 	}
//
// 	if page.Form.State == huh.StateCompleted {
// 		// Generic Containers to store formdata in
// 		var (
// 			value any
// 		)
// 		sb.WriteString("Submitting with: \n")
// 		log.Println("Page: " + page.Title)
// 		switch page.Title {
// 		case "Deposit":
// 			value = s.pages[s.currentPage].Form.GetString("depositAmount")
// 			valueFloat, _ := strconv.ParseFloat(value.(string), 64)
// 			sb.WriteString(fmt.Sprintf("Amount: $%.2f\n", valueFloat))
//
// 			value = s.pages[s.currentPage].Form.GetString("depositDescription")
// 			sb.WriteString("Description: " + value.(string) + "\n")
//
// 			value = s.pages[s.currentPage].Form.GetString("depositUser")
// 			sb.WriteString("User: " + value.(string) + "\n")
// 		case "Withdrawal":
// 			value = s.pages[s.currentPage].Form.GetString("withdrawalAmount")
// 			valueFloat, _ := strconv.ParseFloat(value.(string), 64)
// 			sb.WriteString(fmt.Sprintf("Amount: $%.2f\n", valueFloat))
//
// 			value = s.pages[s.currentPage].Form.GetString("withdrawalDescription")
// 			sb.WriteString("Description: " + value.(string) + "\n")
//
// 			value = s.pages[s.currentPage].Form.GetString("withdrawalUser")
// 			sb.WriteString("User: " + value.(string) + "\n")
// 			// case "Balance":
// 			// 	// When form completed show table
// 			// 	// sb.WriteString(page.Table.View())
// 			// 	return "Not Dead\n\n" + page.Table.View()
// 		}
//
// 		return sb.String()
// 	}
//
// 	// Print out Current Form View
// 	currentPage := s.pages[s.currentPage]
// 	return currentPage.Title + "\n" + currentPage.Form.View()
// }
//
// func NewState(pages []input.Page) *state {
// 	return &state{
// 		pages:       pages,
// 		currentPage: 0,
// 	}
// }
//
// func main() {
// 	noteDeposit := input.NewNoteInput("Deposit")
// 	noteWithdrawal := input.NewNoteInput("Withdrawal")
// 	noteBalance := input.NewNoteInput("Balance")
// 	menuPage := input.NewPage("__Menu__", noteDeposit, noteWithdrawal, noteBalance)
//
// 	input11 := input.NewTextInput("depositAmount", "Amount ", "$0.00")
// 	input12 := input.NewTextInput("depositDescription", "Description ", "Reason...")
// 	input13 := input.NewSelectInput("depositUser", "User ", "One", []string{"One", "Two", "Three"})
// 	pageDeposit := input.NewPage("Deposit", input11, input12, input13)
//
// 	input21 := input.NewTextInput("withdrawalAmount", "Amount", "")
// 	input22 := input.NewTextInput("withdrawalDescription", "Input 2 ", "")
// 	input23 := input.NewSelectInput("withdrawalUser", "Select", "One", []string{"One", "Two", "Three"})
// 	pageWithdrawal := input.NewPage("Withdrawal", input21, input22, input23)
//
// 	input31 := input.NewTextInput("balanceBalance", "Balance", "")
// 	// input32 := input.NewTextInput("balanceTransactions", "Transactions", "")
// 	table33rows := []table.Row{{"2025-01-01", "123.45", "First Deposit"}}
// 	table33 := input.NewTable("Transactions",
// 		[]table.Column{
// 			{Width: 16, Title: "Date"},
// 			{Width: 16, Title: "Amount"},
// 			{Width: 16, Title: "Description"},
// 		},
// 		table33rows)
//
// 	// table33 := table.New(
// 	// 	table.WithColumns(
// 	// 		[]table.Column{
// 	// 			{Title: "Date"},
// 	// 			{Title: "Amount"},
// 	// 			{Title: "Description"},
// 	// 		}),
// 	// )
//
// 	// pageBalance := input.NewPage("Balance", input31, input32, table33)
// 	pageBalance := input.NewPage("Balance", input31, table33)
//
// 	pages := []input.Page{menuPage, pageDeposit, pageWithdrawal, pageBalance}
// 	model := NewState(pages)
// 	app := tea.NewProgram(model)
//
// 	f, err := tea.LogToFile("debug.log", "debug")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	defer func() {
// 		err := f.Close()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}()
//
// 	_, err = app.Run()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
