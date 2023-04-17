package ui

import "github.com/rivo/tview"

func CreateListViewer(app *tview.Application, title string) *tview.List {
	listView := tview.NewList().
		ShowSecondaryText(true).
		SetHighlightFullLine(true)

	listView.SetBorder(true).SetTitle(title)

	return listView
}
