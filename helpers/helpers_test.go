package helpers

import (
	"testing"
)

func TestMin(t *testing.T) {
	tests := []struct {
		a, b     int
		expected int
	}{
		{10, 5, 5},
		{5, 10, 5},
		{-1, -5, -5},
		{0, 5, 0},
		{42, 42, 42},
		{-5, -10, -10},
	}

	for _, tt := range tests {
		got := Min(tt.a, tt.b)
		if got != tt.expected {
			t.Errorf("Min(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.expected)
		}
	}
}

func TestTruncateQuery(t *testing.T) {
	// ShortQueryLen is 64.
	longString := "this_is_a_very_long_string_that_definitely_exceeds_sixty_four_characters"
	expected := longString[0:64]

	tests := []struct {
		orig     string
		expected string
	}{
		{"short", "short"},
		{longString, expected},
		{"", ""},
	}

	for _, tt := range tests {
		got := TruncateQuery(tt.orig)
		if got != tt.expected {
			t.Errorf("TruncateQuery(%q) = %q; want %q", tt.orig, got, tt.expected)
		}
	}
}

func TestHostDropPort(t *testing.T) {
	tests := []struct {
		host     string
		expected string
	}{
		{"localhost:3306", "localhost"},
		{"127.0.0.1:8080", "127.0.0.1"},
		{"mysql-server", "mysql-server"},
		{"db_host:port1:port2", "db_host"}, // Should take first part before colons
		{"", ""},
	}

	for _, tt := range tests {
		got := HostDropPort(tt.host)
		if got != tt.expected {
			t.Errorf("HostDropPort(%q) = %q; want %q", tt.host, got, tt.expected)
		}
	}
}
