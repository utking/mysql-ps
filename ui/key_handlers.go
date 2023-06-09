package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/utking/mysql-ps/helpers"
)

func KeyHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case rune('q'):
		StopApp()
	case rune('p'):
		IsRunning = !IsRunning
	case rune('s'):
		ShowSystem = !ShowSystem
	case rune('?'):
		FlipHelp()
	case rune('l'):
		UIGrid.ResizeItem(UISQLView, 0, BlockHeight10)
		SetFocus(UIListView)
	case rune('v'):
		UIGrid.ResizeItem(UISQLView, 0, BlockHeight10*2)
		SetFocus(UISQLView)
	default:
		switch event.Key() {
		case tcell.KeyESC:
			HideSQLViewer()
			SetFocus(UIListView)
		case tcell.KeyCtrlS:
			if UIListView.GetItemCount() > 0 {
				pri, sec := UIListView.GetItemText(UIListView.GetCurrentItem())

				if err := helpers.WriteSQLLog(fmt.Sprintf("-- %s\n%s", pri, sec), false); err != nil {
					UISQLView.SetText(err.Error())
				}
			}
		case tcell.KeyCtrlA:
			if UIListView.GetItemCount() > 0 {
				pri, sec := UIListView.GetItemText(UIListView.GetCurrentItem())

				if err := helpers.WriteSQLLog(fmt.Sprintf("\n--\n-- %s\n%s", pri, sec), true); err != nil {
					UISQLView.SetText(err.Error())
				}
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
