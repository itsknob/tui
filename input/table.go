package input

import "github.com/charmbracelet/bubbles/table"

type TableView struct {
	table table.Model
}

func NewTable(id string, columns []table.Column) *TableView {
	return &TableView{
		table: table.New(table.WithColumns(columns)),
	}
}
