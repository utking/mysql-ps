package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

func CreateSQLViewer(title string) *tview.TextView {
	reqView := tview.NewTextView()

	reqView.SetBorder(true).SetTitle(title)

	return reqView
}

func PreviewSQL(reqView *tview.TextView, pri, sqlStr string) {
	reqView.SetTitle(fmt.Sprintf("SQL View - %s", pri))
	reqView.SetText(sqlStr)
}

func HideSQLViewer() {
	UISQLView.Clear()
	UISQLView.SetTitle("SQL View")
	UIGrid.ResizeItem(UISQLView, 0, BlockHeightNone)
}
