package ui

import (
	"fmt"
	"os"
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

func UpdateStatusBar(status string, listLen int) {
	statusMessage := fmt.Sprintf("%s (%ds) | Processes:%4d | DSN: %s | Mem: %.2fMB | Show Sys: %v | ? for Help",
		status, TimerSec, listLen, os.Getenv("MYSQL_DSN"), getMemUsage(), ShowSystem)
	UIStatusBar.SetText(statusMessage)
}

// Returns the total allocated memory, in MB
func getMemUsage() float64 {
	var stats runtime.MemStats

	runtime.ReadMemStats(&stats)

	return float64(stats.Alloc) / unitMB
}
