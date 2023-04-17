package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	UIApp       *tview.Application
	UIFlex      *tview.Flex
	UIGrid      *tview.Flex
	UIStatusBar *tview.TextView
	UISQLView   *tview.TextView
	UIListView  *tview.List
	UIMenuBar   *tview.Form
)

func CreateUIGrid() {
	UIApp = tview.NewApplication()

	UIStatusBar = CreateStatusBar("Status")
	UISQLView = CreateSQLViewer("SQL View")
	UIListView = CreateListViewer(UIApp, "Prqocess List")
	UIMenuBar = CreateMenuBar()

	UIGrid = tview.NewFlex().SetDirection(tview.FlexRow)

	UIFlex = tview.NewFlex().
		AddItem(UIGrid.
			AddItem(UIMenuBar, BlockHeightNone, BlockHeightNone, FocusDisable).
			AddItem(UIStatusBar, FixedRowsHeight3, BlockHeight2, FocusDisable).
			AddItem(UIListView, FixedRowsAuto, BlockHeight10, FocusEnable).
			AddItem(UISQLView, FixedRowsAuto, BlockHeightNone, FocusDisable),
			FixedRowsAuto, BlockHeight10, FocusEnable,
		)
}

func SetGlobalHandler(capture func(event *tcell.EventKey) *tcell.EventKey) {
	UIApp.SetInputCapture(capture)
}

func StopApp() {
	UIApp.Stop()
}

func SetFocus(p tview.Primitive) *tview.Application {
	return UIApp.SetFocus(p)
}
