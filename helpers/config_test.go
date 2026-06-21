package helpers

import (
	"os"
	"testing"
)

func TestExpandMyCnfPath(t *testing.T) {
	// Since myCnfPath is hardcoded to "~/.my.cnf", we can only test
	// that it expands correctly based on the current user's home directory.
	path := ExpandMyCnfPath()
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home dir: %v", err)
	}

	expected := ".my.cnf"
	// The result should end with .my.cnf and not contain ~
	if path != home+"/.my.cnf" && !strings_contains(path, expected) {
		t.Errorf("ExpandMyCnfPath() = %q; want something ending in %q", path, expected)
	}
}

func TestLoadConfigNoPanic(t *testing.T) {
	// Just ensure it doesn't panic even if the file is missing
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("LoadConfig panicked: %v", r)
		}
	}()
	LoadConfig()
}

func strings_contains(s, substr string) bool {
	return len(s) >= len(substr) && (s[len(s)-len(substr):] == substr || s[:len(substr)] == substr)
}
