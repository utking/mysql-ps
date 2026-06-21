package ui

import (
	"context"
	"database/sql"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rivo/tview"
	"github.com/utking/mysql-ps/helpers"
)

func TestPSWorker_Success(t *testing.T) {
	var isRunning atomic.Bool
	isRunning.Store(true)
	var ShowSystem atomic.Bool
	ShowSystem.Store(false)
	config := WorkerConfig{
		TimerSec:       0.1,
		IsRunning:      &isRunning,
		StatusBar:      tview.NewTextView(),
		ListView:       tview.NewList(),
		ShowSystem:     &ShowSystem,
		SQLView:        tview.NewTextView(),
		OptionalUpdate: func(fn func()) { fn() }, // Added for testability and to avoid hanging in tests
		Databases:      []string{"db1"},
	}

	listFn := func(filters []string, dbs []any) ([]helpers.ProcessItem, error) {
		return []helpers.ProcessItem{
			{
				ID:      1,
				Host:    "localhost",
				User:    "user1",
				DB:      sql.NullString{String: "db1", Valid: true},
				Info:    sql.NullString{String: "SELECT * FROM table1", Valid: true},
				Time:    10,
				Command: sql.NullString{String: "...", Valid: true},
				State:   sql.NullString{String: "Querying", Valid: true},
			},
			{
				ID:      2,
				Host:    "localhost",
				User:    "system",
				DB:      sql.NullString{Valid: false},
				Info:    sql.NullString{String: "INFORMATION_SCHEMA.PROCESSLIST query", Valid: true},
				Time:    20,
				Command: sql.NullString{String: "...", Valid: true},
				State:   sql.NullString{String: "Sleep", Valid: true},
			},
		}, nil
	}

	dbInterfaces := make([]any, len(config.Databases))
	for i, v := range config.Databases {
		dbInterfaces[i] = v
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go PSWorker(ctx, listFn, dbInterfaces, config)

	// Wait for items to appear with timeout
	timeout := time.After(3 * time.Second)
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-timeout:
			t.Fatalf("timed out waiting for items in ListView, got %d", config.ListView.GetItemCount())
		case <-tick.C:
			if config.ListView.GetItemCount() == 1 {
				return
			}
		}
	}
}

func TestPSWorker_DatabaseError(t *testing.T) {
	var isRunning atomic.Bool
	isRunning.Store(true)
	config := WorkerConfig{
		TimerSec:       0.1,
		IsRunning:      &isRunning,
		StatusBar:      tview.NewTextView(),
		ListView:       tview.NewList(),
		SQLView:        tview.NewTextView(),
		OptionalUpdate: func(fn func()) { fn() }, // Added for testability and to avoid hanging in tests
		Databases:      []string{"db1"},
	}

	errToReturn := errors.New("connection failed")
	listFn := func(filters []string, dbs []interface{}) ([]helpers.ProcessItem, error) {
		return nil, errToReturn
	}

	dbInterfaces := make([]interface{}, len(config.Databases))
	for i, v := range config.Databases {
		dbInterfaces[i] = v
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go PSWorker(ctx, listFn, dbInterfaces, config)

	// Give it enough time to fire at least once and handle the error
	time.Sleep(2 * time.Second)

	if isRunning.Load() {
		t.Error("expected IsRunning to be false after database error")
	}
}

func TestPSWorker_PausedStatus(t *testing.T) {
	var isRunning atomic.Bool
	isRunning.Store(false)
	config := WorkerConfig{
		TimerSec:       0.1,
		IsRunning:      &isRunning,
		StatusBar:      tview.NewTextView(),
		ListView:       tview.NewList(),
		SQLView:        tview.NewTextView(),
		OptionalUpdate: func(fn func()) { fn() }, // Added for testability and to avoid hanging in tests
		Databases:      []string{"db1"},
	}

	listFnCalled := false
	listFn := func(filters []string, dbs []interface{}) ([]helpers.ProcessItem, error) {
		listFnCalled = true
		return nil, nil
	}

	dbInterfaces := make([]interface{}, len(config.Databases))
	for i, v := range config.Databases {
		dbInterfaces[i] = v
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go PSWorker(ctx, listFn, dbInterfaces, config)

	time.Sleep(500 * time.Millisecond)

	if listFnCalled {
		t.Error("expected listFn to not be called while paused")
	}
}
