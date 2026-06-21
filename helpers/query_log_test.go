package helpers

import (
	"os"
	"testing"
)

func TestWriteSQLLogToFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/test_queries.sql"

	tests := []struct {
		name       string
		logLine    string
		appendLine bool
	}{
		{
			name:       "Append mode",
			logLine:    "SELECT * FROM table1;",
			appendLine: true,
		},
		{
			name:       "Truncate mode",
			logLine:    "SELECT * FROM table2;",
			appendLine: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := writeSQLLogToFile(tt.logLine, tmpFile, tt.appendLine)
			if err != nil {
				t.Fatalf("writeSQLLogToFile failed: %v", err)
			}

			content, _ := os.ReadFile(tmpFile)
			got := string(content)
			if tt.appendLine {
				// In append mode, we just check if it's not empty or contains the line
				if got == "" {
					t.Errorf("Expected content for append mode to not be empty, got %q", got)
				}
			} else {
				// In truncate mode, it should be exactly the logLine
				if got != tt.logLine {
					t.Errorf("Expected content to be exactly %q, got %q", tt.logLine, got)
				}
			}
		})
	}
}
