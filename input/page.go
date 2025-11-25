package input

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/huh"
)

type Page struct {
	Title        string
	inputs       []Input
	FocusedInput int
	Form         *huh.Form
	Table        table.Model
}

func NewPage(title string, inputs ...Input) Page {
	var fields []huh.Field
	var table table.Model
	for _, i := range inputs {
		switch i := i.(type) {
		case *TextInput:
			fields = append(fields, i.TextInput)
		case *SelectInput:
			fields = append(fields, i.SelectInput)
		case *NoteInput:
			fields = append(fields, i.Note)
		case *TableView:
			table = i.table
		}
	}
	form := huh.NewForm(huh.NewGroup(fields...)).WithWidth(32)
	if title == "__Menu__" {
		return Page{
			// inputs:       inp,
			Title:        title,
			FocusedInput: 1,
			Form:         form,
			Table:        table,
		}
	}
	return Page{
		// inputs:       inp,
		Title:        title,
		FocusedInput: 0,
		Form:         form,
		Table:        table,
	}
}
