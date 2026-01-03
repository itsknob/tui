package main

/**
*
* TODO:
* todo: Allow date to be empty
* todo: Sanitize dates eg. 2025-12-39
* todo: Pull Transactions
* 	todo: First pass All Transactions
* 	todo: local filter
* 	todo: server api param filter
* todo: Add User Select Box to Forms
* todo: update transactions filter with user select
 */

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2).Height(24)

type activePage int

const (
	MenuPage activePage = iota
	DepositPage
	WithdrawalPage
	BalancePage
	TransactionsPage
)

type transaction struct {
	amount float64
	date   string
	note   string
}

type MainState struct {
	activePage      activePage
	lastSubmit      transaction
	menuOptionsList list.Model
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string {
	return i.title
}

func (mainState MainState) Init() tea.Cmd {
	return nil
}

func (mainState MainState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Println("MainState - Update")
	var cmds []tea.Cmd
	var cmd tea.Cmd

	// switch mainState.activePage {
	// case MenuPage:
	// 	mainState.menuOptionsList, cmd = mainState.menuOptionsList.Update(msg)
	// 	cmds := append(cmds, cmd)
	// 	return mainState, tea.Batch(cmds...)
	// }

	switch msg := msg.(type) {
	case tea.KeyMsg:

		log.Println("MainState - Update - KeyMsg")
		switch msg.String() {
		case "ctrl+c":
			log.Println("MainState - Update - KeyMsg - Ctrl+C - Quitting")
			mainState.activePage = MenuPage
			cmds = append(cmds, tea.Quit)
			// return mainState, tea.Quit
		case "enter":
			mainState.activePage = activePage(mainState.menuOptionsList.Index() + 1) // don't count menuPage
			log.Printf("MainState - Update - KeyMsg - Enter - Selected: %d\n", mainState.activePage)
			switch mainState.activePage {
			case DepositPage:
				log.Println("MainState - Update - KeyMsg - Enter - Deposit - New")
				dp := NewDepositView(mainState)
				return dp, dp.Init()
			case WithdrawalPage:
				log.Println("MainState - Update - KeyMsg - Enter - WithdrawalPage - New")
				// todo: implement WithdrawalPage and update this new func
				wp := NewWithdrawalView(mainState)
				return wp, wp.Init()
			}
			// return mainState, nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		log.Printf("MainState - Update - WindowSizeMsg - W: %d, H: %d", msg.Width-h, msg.Height-v)
		mainState.menuOptionsList.SetSize(msg.Width-h, msg.Height-v)
	case ReturnToMenuMsg:
		if msg.transaction != nil {
			mainState.lastSubmit = *msg.transaction

			// prettyTransaction, _ := json.MarshalIndent(*msg.transaction, "", "  ")

			// err := huh.NewConfirm().Title("Submitting!").Description(string(prettyTransaction)).Affirmative("Okeydokey!").Negative("Sounds Good!").Run()
			// err := huh.NewConfirm().Title("Submitting!").Description(fmt.Sprintf("%+v", *msg.transaction)).Affirmative("Okeydokey!").Negative("Sounds Good!").Run()
			// if err != nil {
			// 	log.Printf("MainState - Update - ReturnToMenuMsg - Failed in Confirm - err: %#v\n", err)
			// }
		}
		log.Printf("MainState - Update - ReturnToMenuMsg - From: %s", msg.from)
		mainState.activePage = 0
		return mainState, nil
	}

	mainState.menuOptionsList, cmd = mainState.menuOptionsList.Update(msg)
	cmds = append(cmds, cmd)
	return mainState, tea.Batch(cmds...)
}

func (mainState MainState) View() string {
	string := fmt.Sprintf("Active Page: %d \n", mainState.activePage)
	return docStyle.Render(string + mainState.menuOptionsList.View())
}

func main() {
	/** Logging **/
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

	/** Setup App **/
	menuItems := []list.Item{
		item{title: "Deposit", desc: ""},
		item{title: "Withdrawal", desc: ""},
	}

	menuList := list.New(menuItems, list.NewDefaultDelegate(), 0, 0)
	menuList.Title = "Menu"

	model := MainState{
		activePage:      MenuPage,
		menuOptionsList: menuList,
	}
	// fmt.Println(menuList.Items())

	// model := NewMainState()
	log.Println("Items: ", model.menuOptionsList.Items())
	app := tea.NewProgram(model, tea.WithInputTTY())

	/** Run App **/
	_, err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

type ReturnToMenuMsg struct {
	from        string
	transaction *transaction
}

func ReturnToMenu(from string, transaction *transaction) tea.Msg {
	log.Printf("ReturnToMenu - Transaction: %+v\n", transaction)
	return ReturnToMenuMsg{
		from:        from,
		transaction: transaction,
	}
}

func SubmitDeposit(amount string, date string, note string) tea.Msg {
	// todo: POST to API
	var datestr string
	amt, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		amt = 0
		math.Copysign(amt, -1.0)
	}

	// Ensure amount is positive
	if amt <= 0 {
		amt *= -1.0
	}

	// Parse Date
	if date != "" {
		dateparsed, err := time.Parse("YYYY-MM-DD", date)
		if err != nil {
			log.Printf("SubmitDeposit - Failed to parse date &%s\n", date)
		}
		datestr = dateparsed.String()
	} else {
		datestr = time.Now().Format("2006-01-02") // weird go date format things
		log.Printf("SubmitDeposit - time.Now() - %s\n", time.Now())
		log.Printf("SubmitDeposit - Date was empty string - Using now: %s\n", datestr)
	}

	// post to API
	log.Printf("SubmitDeposit - Submitting with: \nAmount: %0.2f\nDate: %s\nNote: %s\n", amt, datestr, note)

	// Execute NodeJS Executable that takes in transaction as parameter,
	// connects to actual budget server, then imports that transaction,
	// finally shutting down the connection
	payload := fmt.Sprintf(`{"account":"kaiden","date":"%s","amount":%f,"notes":"%s"}`, datestr, amt, note)
	cmd := exec.Command("./index.js", payload)
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	log.Printf("Output: %s\n", string(out))

	t := transaction{
		amount: amt,
		date:   date,
		note:   note,
	}

	return ReturnToMenu("SubmitDeposit", &t)
}

type DepositView struct {
	mainState MainState
	title     string
	form      *huh.Form
	amount    string
	date      string
	note      string
}

func isDate(s string) error {
	if s == "" {
		log.Println("isDate - date is empty - no error")
		return nil
	}
	dateRegex, err := regexp.Compile(`\d\d\d\d-\d\d-\d\d`)
	if err != nil {
		log.Printf("isDate - could not parse date - %s\n", s)
		return err
	}
	if !dateRegex.MatchString(s) {
		return errors.New("invalid date format - required: YYYY-MM-DD")
	}

	log.Println("isDate - Valid Date")

	return nil
}

func NewDepositView(mainState MainState) *DepositView {
	depositView := &DepositView{
		mainState: mainState,
		title:     "Deposit",
		form:      nil,
		amount:    "",
		date:      "",
		note:      "",
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Amount").Value(&depositView.amount),
			huh.NewInput().Title("Date").Value(&depositView.date).Validate(isDate),
			huh.NewInput().Title("Note").Value(&depositView.note),
		),
	).WithHeight(24).WithWidth(32)
	depositView.form = form
	return depositView
}

// Init implements tea.Model.
func (d *DepositView) Init() tea.Cmd {
	d.form.Init()
	log.Println("Deposit - Init")
	return nil
}

// Update implements tea.Model.
func (d *DepositView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	log.Printf("Deposit - Update")
	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Println("Deposit - Update - KeyMsg")
		switch msg.String() {
		case tea.KeyBackspace.String():
			log.Println("Deposit - Update - KeyMsg - Backspace - ReturnToMenu()")
			model, cmd := d.mainState.Update(ReturnToMenu("Deposit", nil))
			if m, ok := model.(MainState); ok {
				d.mainState = m
			}
			cmds = append(cmds, cmd)
		}
	}

	if d.form.State == huh.StateCompleted {
		return d.mainState.Update(SubmitDeposit(d.amount, d.date, d.note))
	}

	// Form Updates
	model, cmd := d.form.Update(msg)
	if m, ok := model.(*huh.Form); ok {
		d.form = m
	}
	cmds = append(cmds, cmd)

	log.Println("Deposit - Update - Returning")
	return d, tea.Batch(cmds...)
}

// View implements tea.Model.
func (d *DepositView) View() string {
	return fmt.Sprintf("%s\n%s\n", d.title, d.form.View())
}

/* WithdrawalView exists */
type WithdrawalView struct {
	mainState MainState
	title     string
	form      *huh.Form
	amount    string
	date      string
	note      string
}

func NewWithdrawalView(mainState MainState) *WithdrawalView {
	withdrawalView := &WithdrawalView{
		mainState: mainState,
		title:     "Withdrawal",
		form:      nil,
		amount:    "",
		date:      "",
		note:      "",
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Amount").Value(&withdrawalView.amount),
			huh.NewInput().Title("Date").Description("If left empty; will use current date").Value(&withdrawalView.date).Validate(isDate),
			huh.NewInput().Title("Note").Value(&withdrawalView.note),
		).Title("Withdrawal Form"),
	).WithHeight(24).WithWidth(32)
	withdrawalView.form = form
	return withdrawalView
}

// Init implements tea.Model.
func (d *WithdrawalView) Init() tea.Cmd {
	d.form.Init()
	log.Println("Withdrawal - Init")
	return nil
}

// Update implements tea.Model.
func (d *WithdrawalView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	log.Printf("Withdrawal - Update")
	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Println("Withdrawal - Update - KeyMsg")
		switch msg.String() {
		case tea.KeyBackspace.String():
			log.Println("Withdrawal - Update - KeyMsg - Backspace - ReturnToMenu()")
			model, cmd := d.mainState.Update(ReturnToMenu("Withdrawal", nil))
			if m, ok := model.(MainState); ok {
				d.mainState = m
			}
			cmds = append(cmds, cmd)
		}
	}

	if d.form.State == huh.StateCompleted {
		return d.mainState.Update(SubmitWithdrawal(d.amount, d.date, d.note))
	}

	// Form Updates
	model, cmd := d.form.Update(msg)
	if m, ok := model.(*huh.Form); ok {
		d.form = m
	}
	cmds = append(cmds, cmd)

	log.Println("Withdrawal - Update - Returning")
	return d, tea.Batch(cmds...)
}

// View implements tea.Model.
func (d *WithdrawalView) View() string {
	return fmt.Sprintf("%s\n%s\n", d.title, d.form.View())
}

func SubmitWithdrawal(amount string, date string, note string) tea.Msg {
	// todo: POST to API
	var datestr string

	// Parse Ammount
	amt, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		amt = 0                  // If there's an error, don't affect the balance.
		math.Copysign(amt, -1.0) // IEEE754 negative 0
		note = fmt.Sprintf("SubmitWithdrawal - Error Parsing Amount '%s' - Original Description: %s", amount, note)
	}

	// log.Printf("SubmitWithdrawal - Before Negative Coersion - amt: %f\n", amt)

	// Ensure ammount is negative when submitting
	if amt >= 0 {
		amt *= -1.0 // multiply by negtive 1
	}

	// Parse Date
	if date != "" {
		dateparsed, err := time.Parse("YYYY-MM-DD", date)
		if err != nil {
			log.Printf("SubmitWidrawal - Failed to parse date &%s\n", date)
		}
		datestr = dateparsed.String()
	} else {
		datestr = time.Now().Format("2006-01-02")
		log.Printf("SubmitWithdrawal - time.Now() - %s\n", time.Now())
		log.Printf("SubmitWithdrawal - Date was empty string - Using now: %s\n", datestr)
	}
	// if can't parse use current date

	// post to API
	log.Printf("SubmitWithdrawal - Submitting with: \n\tAmount: %0.2f\n\tDate: %s\n\tNote: %s\n", amt, datestr, note)

	// Execute NodeJS Executable that takes in transaction as parameter,
	// connects to actual budget server, then imports that transaction,
	// finally shutting down the connection
	// todo:
	payload := fmt.Sprintf(`{"account":"kaiden","date":"%s","amount":%f,"notes":"%s"}`, datestr, amt, note)
	cmd := exec.Command("./index.js", payload)
	out, err := cmd.Output()
	if err != nil {
		log.Printf("SubmitWithdrawal - Error: %+v\n", err)
	}
	log.Printf("SubmitWithdrawal - Submitted Successfully - Output: %s\n\n", string(out))

	t := transaction{
		amount: amt,
		date:   date,
		note:   note,
	}

	return ReturnToMenu("SubmitWithdrawal", &t)
}
