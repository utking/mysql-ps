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
		name     string
		input    string
		expected string
	}{
		{"Standard case", "127.0.0.1:3306", "127.0.0.1"},
		{"Localhost port", "localhost:8080", "localhost"},
		{"No port", "mysql-server", "mysql-server"},
		{"Multiple colons", "db_host.example.com:3306:8080", "db_host.example.com"},
		{"Missing host", ":3306", ""},
		{"Empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HostDropPort(tt.input)
			if got != tt.expected {
				t.Errorf("HostDropPort(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
