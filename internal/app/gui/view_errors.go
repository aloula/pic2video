package gui

import "fyne.io/fyne/v2/widget"

func ShowValidationError(target *widget.Label, msg string) {
	if target == nil {
		return
	}
	target.SetText(msg)
}
