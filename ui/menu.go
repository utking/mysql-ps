package ui

import (
	"github.com/rivo/tview"
)

const (
	FocusEnable      = true
	FocusDisable     = false
	FixedRowsAuto    = 0
	FixedRowsHeight3 = 3
	BlockHeightNone  = 0
	BlockHeight2     = 2
	BlockHeight10    = 10
)

var (
	menuLabels = []string{
		"Pause (P)",
		"Show sys DB (S)",
		"Show View (Ent)",
		"Hide View (Esc)",
		"To list (L)",
		"To view (V)",
		"Save SQL (Crtl+S)",
		"Append SQL (Ctrl+A)",
		"Quit (Q)",
	}
	MenuVisible bool
)

func CreateMenuBar() *tview.Form {
	menuBar := tview.NewForm()

	for _, lbl := range menuLabels {
		menuBar = menuBar.AddButton(lbl, func() {})
	}

	return menuBar
}
