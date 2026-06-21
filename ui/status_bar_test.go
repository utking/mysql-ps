package ui

import (
	"testing"

	"github.com/rivo/tview"
)

func TestFormatStatusBar(t *testing.T) {
	status := "Running"
	timerSec := float32(1.5)
	listLen := 10
	showSys := true
	dsn := "mysql://user:pass@tcp(host:3306)/db"
	memUsage := 123.45

	msg := FormatStatusBar(status, timerSec, listLen, showSys, dsn, memUsage)

	expected := "Running (1.5s) | Processes:  10 | DSN: mysql://user:pass@tcp(host:3306)/db | Mem: 123.45MB | Show Sys: true | ? for Help"
	if msg != expected {
		t.Errorf("Expected %q, got %q", expected, msg)
	}
}

func TestUpdateStatusBar(t *testing.T) {
	bar := tview.NewTextView()
	status := "Running"
	listLen := 10
	timerSec := float32(1.5)
	showSys := true
	dsn := "mysql://user:pass@tcp(host:3306)/db"
	memUsage := 123.45

	UpdateStatusBar(bar, status, listLen, timerSec, showSys, dsn, memUsage)

	expected := "Running (1.5s) | Processes:  10 | DSN: mysql://user:pass@tcp(host:3306)/db | Mem: 123.45MB | Show Sys: true | ? for Help"
	if bar.GetText(true) != expected {
		t.Errorf("Expected %q, got %q", expected, bar.GetText(true))
	}
}
