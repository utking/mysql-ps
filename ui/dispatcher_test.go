package ui

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/rivo/tview"
	"github.com/utking/mysql-ps/helpers"
)

func newTestUI() *UIComponents {
	return &UIComponents{
		StatusBar: tview.NewTextView(),
		ListView:  tview.NewList(),
		SQLView:   tview.NewTextView(),
	}
}

func TestPerformUpdate_StateTransitions(t *testing.T) {
	ui := newTestUI()
	ui.TimerSec = 1.0
	ui.IsRunning.Store(false)
	ui.ShowSystem.Store(true)

	config := WorkerConfig{
		UI:             ui,
		Databases:      []string{"db1"},
		OptionalUpdate: func(fn func()) { fn() },
	}

	listFn := func(filters []string, dbs []any) ([]helpers.ProcessItem, error) {
		return []helpers.ProcessItem{
			{ID: 1, Info: sql.NullString{String: "Query 1", Valid: true}, DB: sql.NullString{String: "db1", Valid: true}, Time: 10, User: "user1", Host: "localhost", State: sql.NullString{String: "Active", Valid: true}},
		}, nil
	}

	ui.IsRunning.Store(true)
	performUpdate(&config, listFn)

	status := config.UI.StatusBar.GetText(true)
	if !strings.Contains(status, "Running") || !strings.Contains(status, "Processes:") || !strings.Contains(status, "1") {
		t.Errorf("expected status bar to show Running state with 1 process, got %s", status)
	}

	ui.IsRunning.Store(false)
	performUpdate(&config, listFn)

	status = config.UI.StatusBar.GetText(true)
	if !strings.Contains(status, "Paused") || !strings.Contains(status, "Processes:") || !strings.Contains(status, "0") {
		t.Errorf("expected status bar to show Paused state with 0 processes, got %s", status)
	}
}

func TestPerformUpdate_SystemFilter(t *testing.T) {
	ui := newTestUI()
	ui.IsRunning.Store(true)
	ui.ShowSystem.Store(false)

	config := WorkerConfig{
		UI:             ui,
		Databases:      []string{"db1"},
		OptionalUpdate: func(fn func()) { fn() },
	}

	listFn := func([]string, []any) ([]helpers.ProcessItem, error) {
		return []helpers.ProcessItem{
			{ID: 1, Info: sql.NullString{String: "Query 1", Valid: true}, DB: sql.NullString{String: "db1", Valid: true}, Time: 10, User: "user1", Host: "localhost", State: sql.NullString{String: "Active", Valid: true}},
			{ID: 2, Info: sql.NullString{String: "INFORMATION_SCHEMA.PROCESSLIST query", Valid: true}, DB: sql.NullString{String: "db1", Valid: true}, Time: 10, User: "user1", Host: "localhost", State: sql.NullString{String: "Active", Valid: true}},
		}, nil
	}

	performUpdate(&config, listFn)

	status := config.UI.StatusBar.GetText(true)
	if !strings.Contains(status, "Processes:") || !strings.Contains(status, "1") {
		t.Errorf("expected status bar to show 1 process (filtered out system query), got %s", status)
	}
}
