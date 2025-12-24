package main

import (
	"fmt"
	"log"

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

type MainState struct {
	activePage      activePage
	menuOptionsList list.Model
}

// func NewMainState() MainState {
// 	menuItems := []list.Item{
// 		MenuItem("Deposit"),
// 		MenuItem("Withdrawal"),
// 	}
// 	menuList := list.New(menuItems, list.NewDefaultDelegate(), 20, 16)
// 	menuList.Title = "Menu"
// 	// fmt.Println(menuList.Items())
//
// 	return MainState{
// 		activePage:      MenuPage,
// 		menuOptionsList: menuList,
// 	}
// }

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
	// var cmds []tea.Cmd
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
			return mainState, tea.Quit
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
			return mainState, nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		log.Printf("MainState - Update - WindowSizeMsg - W: %d, H: %d", msg.Width-h, msg.Height-v)
		mainState.menuOptionsList.SetSize(msg.Width-h, msg.Height-v)
	case ReturnToMenuMsg:
		log.Printf("MainState - Update - ReturnToMenuMsg - From: %s", msg.from)
		mainState.activePage = 0
		return mainState, nil
	}

	mainState.menuOptionsList, cmd = mainState.menuOptionsList.Update(msg)
	// cmds = append(cmds, cmd)
	return mainState, cmd
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
	log.Default().Println("Items: ", model.menuOptionsList.Items())
	app := tea.NewProgram(model)

	/** Run App **/
	_, err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

type ReturnToMenuMsg struct {
	from string
}

func ReturnToMenu(from string) tea.Msg {
	return ReturnToMenuMsg{
		from: from,
	}
}

func SubmitDeposit(amount string, date string, note string) tea.Msg {
	// todo: POST to API

	return ReturnToMenu("SubmitDeposit")
}

type DepositView struct {
	mainState MainState
	title     string
	form      *huh.Form
	amount    string
	date      string
	note      string

	// amount float64
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
			huh.NewInput().Title("Date").Value(&depositView.date),
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
			model, cmd := d.mainState.Update(ReturnToMenu("Deposit"))
			if m, ok := model.(MainState); ok {
				d.mainState = m
			}
			cmds = append(cmds, cmd)
			// case "enter":
			// 	log.Println("Deposit - Update - Enter - Qutting")
			// 	return d, tea.Quit
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
	// amount float64
}

func NewWithdrawalView(mainState MainState) *WithdrawalView {
	return &WithdrawalView{
		mainState: mainState,
		title:     "Withdrawal",
	}
}

// Init implements tea.Model.
func (d *WithdrawalView) Init() tea.Cmd {
	log.Println("Withdrawal - Init")
	return nil
}

// Update implements tea.Model.
func (d *WithdrawalView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("Withdrawal - Update")
	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Println("Withdrawal - Update - KeyMsg")
		switch msg.String() {
		case tea.KeyBackspace.String():
			log.Println("Withdrawal - Update - KeyMsg - Backspace - ReturnToMenu()")
			return d.mainState.Update(ReturnToMenu("Withdrawal"))
		case "enter":
			log.Println("Withdrawal - Update - Enter - Qutting")
			return d, tea.Quit
		}
	}

	log.Println("Withdrawal - Update - Returning")
	return d, nil
}

// View implements tea.Model.
func (d *WithdrawalView) View() string {
	return fmt.Sprintf("%s\nfdsa\n", d.title)
}
