package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type FormView struct {
	InputEntry  *widget.Entry
	OutputEntry *widget.Entry
	Container   *fyne.Container
}

func NewFormView(cfg GuiRunConfiguration) *FormView {
	input := widget.NewEntry()
	input.SetText(cfg.InputFolder)
	output := widget.NewEntry()
	output.SetText(cfg.OutputFolder)
	output.CursorColumn = len([]rune(output.Text))

	c := container.NewVBox(
		widget.NewLabel("Input folder"),
		input,
		widget.NewLabel("Output folder"),
		output,
	)
	return &FormView{InputEntry: input, OutputEntry: output, Container: c}
}
