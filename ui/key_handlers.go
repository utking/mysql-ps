package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/utking/mysql-ps/helpers"
)

func (c *UIComponents) KeyHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case rune('q'):
		c.StopApp()
	case rune('p'):
		current := c.IsRunning.Load()
		c.IsRunning.Store(!current)
		select {
		case c.updateTriggerChan <- struct{}{}:
		default:
		}
	case rune('s'):
		c.ShowSystem.Store(!c.ShowSystem.Load())
		select {
		case c.updateTriggerChan <- struct{}{}:
		default:
		}
	case rune('?'):
		c.FlipHelp()
		select {
		case c.updateTriggerChan <- struct{}{}:
		default:
		}
	case rune('l'):
		c.Grid.ResizeItem(c.SQLView, 0, BlockHeight10)
		c.SetFocus(c.ListView)
		select {
		case c.updateTriggerChan <- struct{}{}:
		default:
		}
	case rune('v'):
		c.Grid.ResizeItem(c.SQLView, 0, BlockHeight10*2)
		c.SetFocus(c.SQLView)
		select {
		case c.updateTriggerChan <- struct{}{}:
		default:
		}
	default:
		switch event.Key() {
		case tcell.KeyESC:
			c.HideSQLViewer()
			c.SetFocus(c.ListView)
			select {
			case c.updateTriggerChan <- struct{}{}:
			default:
			}
		case tcell.KeyCtrlS:
			if c.ListView.GetItemCount() > 0 {
				pri, sec := c.ListView.GetItemText(c.ListView.GetCurrentItem())

				if err := helpers.WriteSQLLog(fmt.Sprintf("-- %s\n%s", pri, sec), false); err != nil {
					c.SQLView.SetText(err.Error())
				}
			}
		case tcell.KeyCtrlA:
			if c.ListView.GetItemCount() > 0 {
				pri, sec := c.ListView.GetItemText(c.ListView.GetCurrentItem())

				if err := helpers.WriteSQLLog(fmt.Sprintf("\n--\n-- %s\n%s", pri, sec), true); err != nil {
					c.SQLView.SetText(err.Error())
				}
			}
		}
	}

	return event
}

func (c *UIComponents) OpenSQLQuery(i int, s1, s2 string, r rune) {
	pri, sec := c.ListView.GetItemText(c.ListView.GetCurrentItem())

	PreviewSQL(c.SQLView, pri, sec)
	c.App.SetFocus(c.ListView)
	if c.Grid != nil {
		c.Grid.ResizeItem(c.SQLView, 0, BlockHeight10)
	}
}
