package ui

import (
	"fmt"
	"runtime"

	"github.com/rivo/tview"
)

const (
	unitMB = 1024 * 1024.0
)

func CreateStatusBar(title string) *tview.TextView {
	bar := tview.NewTextView().SetText("Running")

	bar.SetBorder(true).SetBorderPadding(0, 0, 1, 0).SetTitle(title)

	return bar
}

func FormatStatusBar(status string, timerSec float32, listLen int, showSys bool, dsn string, memUsage float64) string {
	return fmt.Sprintf("%s (%.1fs) | Processes:%4d | DSN: %s | Mem: %.2fMB | Show Sys: %v | ? for Help",
		status, timerSec, listLen, dsn, memUsage, showSys)
}

func UpdateStatusBar(bar *tview.TextView, status string, listLen int, timerSec float32, showSys bool, dsn string, memUsage float64) {
	statusMessage := FormatStatusBar(status, timerSec, listLen, showSys, dsn, memUsage)
	bar.SetText(statusMessage)
}

// Returns the total allocated memory, in MB
func getMemUsage() float64 {
	var stats runtime.MemStats

	runtime.ReadMemStats(&stats)

	return float64(stats.Alloc) / unitMB
}
