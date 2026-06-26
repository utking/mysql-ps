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

func (c *UIComponents) HideSQLViewer() {
	c.SQLView.Clear()
	c.SQLView.SetTitle("SQL View")
	c.Grid.ResizeItem(c.SQLView, 0, BlockHeightNone)
}
