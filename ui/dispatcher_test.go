package ui

import (
	"database/sql"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/rivo/tview"
	"github.com/utking/mysql-ps/helpers"
)

func TestPerformUpdate_StateTransitions(t *testing.T) {
	var isRunning atomic.Bool
	isRunning.Store(false) // Start Paused
	var ShowSystem atomic.Bool
	ShowSystem.Store(true)

	config := WorkerConfig{
		TimerSec:       1.0,
		IsRunning:      &isRunning,
		StatusBar:      tview.NewTextView(),
		ListView:       tview.NewList(),
		ShowSystem:     &ShowSystem,
		SQLView:        tview.NewTextView(),
		OptionalUpdate: func(fn func()) { fn() },
		Databases:      []string{"db1"},
	}

	listFn := func(filters []string, dbs []any) ([]helpers.ProcessItem, error) {
		return []helpers.ProcessItem{
			{ID: 1, Info: sql.NullString{String: "Query 1", Valid: true}, DB: sql.NullString{String: "db1", Valid: true}, Time: 10, User: "user1", Host: "localhost", State: sql.NullString{String: "Active", Valid: true}},
		}, nil
	}

	// Test transition from Paused to Running via direct state change (simulating UI interaction)
	isRunning.Store(true)
	performUpdate(&config, listFn)

	status := config.StatusBar.GetText(true)
	if !strings.Contains(status, "Running") || !strings.Contains(status, "Processes:") || !strings.Contains(status, "1") {
		t.Errorf("expected status bar to show Running state with 1 process, got %s", status)
	}

	// Test transition from Running back to Paused
	isRunning.Store(false)
	performUpdate(&config, listFn)

	status = config.StatusBar.GetText(true)
	if !strings.Contains(status, "Paused") || !strings.Contains(status, "Processes:") || !strings.Contains(status, "0") {
		t.Errorf("expected status bar to show Paused state with 0 processes, got %s", status)
	}
}

func TestPerformUpdate_SystemFilter(t *testing.T) {
	var isRunning atomic.Bool
	isRunning.Store(true)
	var ShowSystem atomic.Bool
	ShowSystem.Store(false)
	config := WorkerConfig{
		IsRunning:      &isRunning,
		StatusBar:      tview.NewTextView(),
		ListView:       tview.NewList(),
		SQLView:        tview.NewTextView(),
		OptionalUpdate: func(fn func()) { fn() },
		Databases:      []string{"db1"},
		ShowSystem:     &ShowSystem,
	}

	listFn := func([]string, []any) ([]helpers.ProcessItem, error) {
		return []helpers.ProcessItem{
			{ID: 1, Info: sql.NullString{String: "Query 1", Valid: true}, DB: sql.NullString{String: "db1", Valid: true}, Time: 10, User: "user1", Host: "localhost", State: sql.NullString{String: "Active", Valid: true}},
			{ID: 2, Info: sql.NullString{String: "INFORMATION_SCHEMA.PROCESSLIST query", Valid: true}, DB: sql.NullString{String: "db1", Valid: true}, Time: 10, User: "user1", Host: "localhost", State: sql.NullString{String: "Active", Valid: true}},
		}, nil
	}

	performUpdate(&config, listFn)

	status := config.StatusBar.GetText(true)
	if !strings.Contains(status, "Processes:") || !strings.Contains(status, "1") {
		t.Errorf("expected status bar to show 1 process (filtered out system query), got %s", status)
	}
}