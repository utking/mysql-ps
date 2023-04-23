package ui

import (
	"fmt"
	"os"

	"github.com/rivo/tview"
)

func CreateStatusBar(title string) *tview.TextView {
	bar := tview.NewTextView().SetText("Running")

	bar.SetBorder(true).SetBorderPadding(0, 0, 1, 0).SetTitle(title)

	return bar
}

func UpdateStatusBar(status string, listLen int) {
	statusMessage := fmt.Sprintf("%s (%ds) | Processes:%4d | DSN: %s | ? for Help", status, TimerSec, listLen, os.Getenv("MYSQL_DSN"))
	UIStatusBar.SetText(statusMessage)
}
