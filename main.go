package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	Amount      float64
	Description string
	Date        string
	User        string
)

func main() {
	// pages := []string{"Deposit", "Withdrawal", "Balance"}

	pages := []string{
		"Menu",
		"Deposit",
		"Widthdrawal",
		"Balance",
	}

	m := NewModel(pages)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := tea.NewProgram(m)

	_, err = p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

type model struct {
	pages       []string
	cursor      int
	currentPage string
}

func NewModel(pages []string) *model {
	m := new(model)
	m.pages = pages
	m.cursor = 0
	m.currentPage = "Menu"
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

var count int = 0

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	count++
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			m.Prev()
		case "down":
			m.Next()
		case "enter":
			if m.currentPage == m.pages[m.cursor] {
				return m, nil
			}
			m.currentPage = m.pages[m.cursor]
			return m, nil
		case "q":
			return NewModel(m.pages), nil
		}
	}
	return m, nil
}

func (m *model) Prev() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.pages) - 1
	}
}

func (m *model) Next() {
	m.cursor++
	if m.cursor > len(m.pages)-1 {
		m.cursor = 0
	}
}

func (m model) View() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%s\n", m.currentPage))
	s.WriteString(fmt.Sprintf("Count: %d\n\n", count))
	for idx, p := range m.pages {
		if idx == m.cursor {
			s.WriteString(">")
		} else {
			s.WriteString(" ")
		}
		s.WriteString(" ")
		s.WriteString(p)
		s.WriteString("\n")
	}
	return s.String()
}
