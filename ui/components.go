package ui

import (
	"sync/atomic"

	"github.com/rivo/tview"
)

type UIComponents struct {
	App       *tview.Application
	Flex      *tview.Flex
	Grid      *tview.Flex
	StatusBar *tview.TextView
	SQLView   *tview.TextView
	ListView  *tview.List
	MenuBar   *tview.Form

	TimerSec    float32
	ShowSystem  atomic.Bool
	IsRunning   atomic.Bool
	UseMouse    bool
	MenuVisible bool

	updateTriggerChan chan struct{}
}

func NewUI() *UIComponents {
	c := &UIComponents{
		App:               tview.NewApplication(),
		updateTriggerChan: make(chan struct{}, 1),
	}

	c.StatusBar = CreateStatusBar("Status")
	c.SQLView = CreateSQLViewer("SQL View")
	c.ListView = CreateListViewer(c.App, "Process List")
	c.MenuBar = CreateMenuBar()

	c.Flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(c.MenuBar, BlockHeightNone, BlockHeightNone, FocusDisable).
		AddItem(c.StatusBar, FixedRowsHeight3, BlockHeight2, FocusDisable).
		AddItem(c.ListView, FixedRowsAuto, BlockHeight10, FocusEnable).
		AddItem(c.SQLView, FixedRowsAuto, BlockHeightNone, FocusDisable)
	c.Grid = c.Flex

	return c
}
