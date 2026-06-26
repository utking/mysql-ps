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

func (c *UIComponents) triggerUpdate() {
	select {
	case <-c.updateTriggerChan:
	default:
	}
	c.updateTriggerChan <- struct{}{}
}

func (c *UIComponents) FlipHelp() {
	c.helpVisible = !c.helpVisible
	if c.helpVisible {
		c.Pages.ShowPage("help")
		c.App.SetFocus(c.Pages)
	} else {
		c.Pages.HidePage("help")
		c.App.SetFocus(c.ListView)
	}
}

