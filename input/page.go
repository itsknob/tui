package input

import (
	"github.com/charmbracelet/huh"
)

type Page struct {
	Title        string
	inputs       []Input
	FocusedInput int
	Form         *huh.Form
}

func NewPage(title string, inputs ...Input) Page {
	var fields []huh.Field
	for _, i := range inputs {
		switch i := i.(type) {
		case *TextInput:
			fields = append(fields, i.TextInput)
		case *SelectInput:
			fields = append(fields, i.SelectInput)
		case *NoteInput:
			fields = append(fields, i.Note)
		}
	}
	form := huh.NewForm(huh.NewGroup(fields...)).WithWidth(32)
	if title == "__Menu__" {
		return Page{
			// inputs:       inp,
			Title:        title,
			FocusedInput: 1,
			Form:         form,
		}
	}
	return Page{
		// inputs:       inp,
		Title:        title,
		FocusedInput: 0,
		Form:         form,
	}
}
