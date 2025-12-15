package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2).Height(32)

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

		switch msg.String() {
		case "ctrl+c":
			return mainState, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		mainState.menuOptionsList.SetSize(msg.Width-h, msg.Height-v)
	}

	mainState.menuOptionsList, cmd = mainState.menuOptionsList.Update(msg)
	// cmds = append(cmds, cmd)
	return mainState, cmd
}

func (mainState MainState) View() string {
	return docStyle.Render(mainState.menuOptionsList.View())
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
