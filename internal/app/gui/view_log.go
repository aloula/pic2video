package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type LogView struct {
	Label     *widget.Label
	Scroll    *container.Scroll
	Container *fyne.Container
}

func NewLogView() *LogView {
	label := widget.NewLabel("")
	label.Wrapping = fyne.TextWrapWord
	scroll := container.NewVScroll(label)
	scroll.SetMinSize(fyne.NewSize(0, 88))
	c := container.NewBorder(widget.NewLabel("Log"), nil, nil, nil, scroll)
	return &LogView{Label: label, Scroll: scroll, Container: c}
}

func (v *LogView) SetText(txt string) {
	v.Label.SetText(txt)
	v.Scroll.ScrollToBottom()
}
