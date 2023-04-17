package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/utking/mysql-ps/helpers"
)

func KeyHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Rune() == rune('q') {
		StopApp()
	} else if event.Rune() == rune('p') {
		IsRunning = !IsRunning
	} else if event.Key() == tcell.KeyF1 {
		FlipHelp()
	} else if event.Key() == tcell.KeyF2 {
		UIGrid.ResizeItem(UISQLView, 0, BlockHeight10)
		SetFocus(UIListView)
	} else if event.Key() == tcell.KeyF3 {
		UIGrid.ResizeItem(UISQLView, 0, BlockHeight10*2)
		SetFocus(UISQLView)
	} else if event.Key() == tcell.KeyESC {
		HideSQLViewer()
		SetFocus(UIListView)
	} else if event.Key() == tcell.KeyCtrlS {
		if UIListView.GetItemCount() > 0 {
			pri, sec := UIListView.GetItemText(UIListView.GetCurrentItem())

			if err := helpers.WriteSQLLog(fmt.Sprintf("-- %s\n%s", pri, sec), false); err != nil {
				UISQLView.SetText(err.Error())
			}
		}
	} else if event.Key() == tcell.KeyCtrlA {
		if UIListView.GetItemCount() > 0 {
			pri, sec := UIListView.GetItemText(UIListView.GetCurrentItem())

			if err := helpers.WriteSQLLog(fmt.Sprintf("\n--\n-- %s\n%s", pri, sec), true); err != nil {
				UISQLView.SetText(err.Error())
			}
		}
	}

	return event
}

func OpenSQLQuery(i int, s1, s2 string, r rune) {
	pri, sec := UIListView.GetItemText(UIListView.GetCurrentItem())

	PreviewSQL(UISQLView, pri, sec)
	UIApp.SetFocus(UIListView)
	UIGrid.ResizeItem(UISQLView, 0, BlockHeight10)
}
