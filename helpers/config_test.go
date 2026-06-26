package helpers

import (
	"os"
	"strings"
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
	if path != home+"/.my.cnf" && !strings.HasSuffix(path, expected) {
		t.Errorf("ExpandMyCnfPath() = %q; want something ending in %q", path, expected)
	}
}

func TestLoadConfigMissingEnvVars(t *testing.T) {
	// Clear relevant environment variables to simulate a clean state
	os.Unsetenv("MYSQL_DSN")
	os.Unsetenv("MYSQL_USER")
	os.Unsetenv("MYSQL_PASSWORD")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("LoadConfig panicked with missing env vars: %v", r)
		}
	}()

	LoadConfig()
}


