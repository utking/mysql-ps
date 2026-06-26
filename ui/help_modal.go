package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const helpContent = `[yellow]Keyboard Shortcuts[white]

  [green]P[white]          Pause/resume refresh
  [green]S[white]          Toggle system queries
  [green]Q[white]          Quit application
  [green]?[white]          Toggle this help

  [green]Enter[white]      Show full SQL preview
  [green]Esc[white]        Hide SQL preview / Close this help
  [green]L[white]          Focus process list
  [green]V[white]          Focus SQL view

  [green]Ctrl+S[white]     Save selected query to file
  [green]Ctrl+A[white]     Append selected query to file`

func CreateHelpModal(c *UIComponents) *tview.Flex {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(helpContent).
		SetTextAlign(tview.AlignLeft)

	frame := tview.NewFrame(textView).
		SetBorders(2, 2, 1, 1, 2, 2).
		AddText("Help", true, tview.AlignCenter, tcell.ColorYellow)

	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(frame, 0, 3, true).
			AddItem(nil, 0, 1, false),
			0, 3, false).
		AddItem(nil, 0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Rune() == '?' {
			c.FlipHelp()
			return nil
		}
		return event
	})

	return flex
}
