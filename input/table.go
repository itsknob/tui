package input

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type TableView struct {
	table table.Model
}

func NewTable(id string, columns []table.Column, rows []table.Row) *TableView {
	table := table.New()
	table.SetColumns(columns)
	table.SetRows(rows)
	return &TableView{table}
}

func (t TableView) Focus() tea.Cmd {
	t.table.Focus()
	return nil
}

func (t TableView) Blur() tea.Cmd {
	t.table.Blur()
	return nil
}

func (t TableView) Update(msg tea.Msg) (Input, tea.Cmd) {
	model, cmd := t.table.Update(msg)
	t.table = model
	return t, cmd
}
