package ui

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/utking/mysql-ps/helpers"
)

// Mock data for process items
var (
	mockSuccessItems = []helpers.ProcessItem{
		{
			Host:    "127.0.0.1:3306",
			User:    "root",
			Command: sql.NullString{String: "SELECT * FROM users;"},
			State:   sql.NullString{String: "Sending data"},
			DB:      sql.NullString{String: "production"},
			Info:    sql.NullString{String: "Querying production table"},
			ID:      1,
			Time:    100,
		},
		{
			Host:    "127.0.0.1:3306",
			User:    "admin",
			Command: sql.NullString{String: "UPDATE orders SET status = 'shipped' WHERE id < 10;"},
			State:   sql.NullString{String: "Updating rows"},
			DB:      sql.NullString{String: "warehouse"},
			Info:    sql.NullString{String: "Updating warehouse records"},
			ID:      2,
			Time:    250,
		},
	}

	mockSystemItems = []helpers.ProcessItem{
		{
			Host:    "127.0.0.1:3306",
			User:    "system",
			Command: sql.NullString{String: "SHOW PROCESSLIST;"},
			State:   sql.NullString{String: "Sleep"},
			DB:      sql.NullString{}, // System queries might have null DB
			Info:    sql.NullString{String: "INFORMATION_SCHEMA.PROCESSLIST internal query"},
			ID:      3,
			Time:    500,
		},
	}
)

// listFn is a function type that returns process items and an error given filters and databases.
type listFn func([]string, []interface{}) ([]helpers.ProcessItem, error)

var (
	listFnSuccess listFn = func(filters []string, dbs []interface{}) ([]helpers.ProcessItem, error) {
		return mockSuccessItems, nil
	}

	listFnFailure listFn = func(filters []string, dbs []interface{}) ([]helpers.ProcessItem, error) {
		return nil, errors.New("connection refused")
	}

	listFnSystemOnly listFn = func(filters []string, dbs []interface{}) ([]helpers.ProcessItem, error) {
		return mockSystemItems, nil
	}
)

func TestMockHelpers(t *testing.T) {
	// Test listFnSuccess
	items, err := listFnSuccess(nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(items, mockSuccessItems) {
		t.Errorf("expected %v, got %v", mockSuccessItems, items)
	}

	// Test listFnFailure
	items, err = listFnFailure(nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "connection refused" {
		t.Errorf("expected 'connection refused' error, got %v", err)
	}

	// Test listFnSystemOnly
	items, err = listFnSystemOnly(nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !reflect.DeepEqual(items, mockSystemItems) {
		t.Errorf("expected %v, got %v", mockSystemItems, items)
	}
}
