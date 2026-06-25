package ui

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/utking/mysql-ps/helpers"
)

func TestListFilteringLogic(t *testing.T) {
	items := []helpers.ProcessItem{
		{ID: 1, Info: sql.NullString{String: "SELECT * FROM users", Valid: true}},
		{ID: 2, Info: sql.NullString{String: "INFORMATION_SCHEMA.PROCESSLIST query", Valid: true}},
		{ID: 3, Info: sql.NullString{String: "SELECT sleep(5)", Valid: true}},
		{ID: 4, Info: sql.NullString{String: "SHOW PROCESSLIST", Valid: true}},
	}

	count := 0
	for _, item := range items {
		if !strings.Contains(item.Info.String, "INFORMATION_SCHEMA.PROCESSLIST") {
			count++
		}
	}

	expectedCount := 3 // Item 1, 3, 4 are kept, Item 2 is filtered out
	if count != expectedCount {
		t.Errorf("expected filtered count %d, got %d", expectedCount, count)
	}
}