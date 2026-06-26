package ui

import (
	"github.com/rivo/tview"
)

func (c *UIComponents) SetGlobalHandler() {
	c.App.SetInputCapture(c.KeyHandler)
}

func (c *UIComponents) StopApp() {
	c.App.Stop()
}

func (c *UIComponents) SetFocus(p tview.Primitive) *tview.Application {
	return c.App.SetFocus(p)
}

func (c *UIComponents) FlipHelp() {
	c.MenuVisible = !c.MenuVisible

	if c.MenuVisible {
		c.Grid.ResizeItem(c.MenuBar, FixedRowsHeight3, BlockHeight2)
	} else {
		c.Grid.ResizeItem(c.MenuBar, BlockHeightNone, BlockHeightNone)
	}
}

