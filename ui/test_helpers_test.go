package ui

import (
	"database/sql"

	"github.com/utking/mysql-ps/helpers"
)

var mockItems = []helpers.ProcessItem{
	{
		ID:      1,
		Host:    "localhost",
		User:    "root",
		Command: sql.NullString{String: "SELECT * FROM users;", Valid: true},
		State:   sql.NullString{String: "Sending data", Valid: true},
		DB:      sql.NullString{String: "mysql", Valid: true},
		Info:    sql.NullString{String: "SELECT * FROM users", Valid: true},
		Time:    100,
	},
	{
		ID:      2,
		Host:    "localhost",
		User:    "admin",
		Command: sql.NullString{String: "UPDATE orders SET status = 'shipped' WHERE id < 10;", Valid: true},
		State:   sql.NullString{String: "Updating rows", Valid: true},
		DB:      sql.NullString{String: "warehouse", Valid: true},
		Info:    sql.NullString{String: "Updating warehouse records", Valid: true},
		Time:    250,
	},
	{
		ID:      3,
		Host:    "localhost",
		User:    "system",
		Command: sql.NullString{String: "SHOW PROCESSLIST;", Valid: true},
		State:   sql.NullString{String: "Sleep", Valid: true},
		DB:      sql.NullString{Valid: false}, // System queries might not have a DB name
		Info:    sql.NullString{String: "INFORMATION_SCHEMA.PROCESSLIST internal query", Valid: true},
		Time:    500,
	},
}

func MockListFn() func([]string, []interface{}) ([]helpers.ProcessItem, error) {
	return func(filters []string, dbs []interface{}) ([]helpers.ProcessItem, error) {
		// For now, just return the mockItems.
		// In tests, you can choose to filter them or return an error.
		return mockItems, nil
	}
}
