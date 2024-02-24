package utils

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	testCases := []struct {
		input    string
		expected time.Duration
		err      bool
	}{
		{"1:0:0", 1 * time.Hour, false},
		{"1:61:64", 2*time.Hour + 2*time.Minute + 4*time.Second, false},
		{"0:1:0", 1*time.Minute, false},
		{"1:23:40", 1*time.Hour + 23*time.Minute + 40*time.Second, false},
		{"23:40", 23*time.Minute + 40*time.Second, false},
		{"40", 40 * time.Second, false},
		{"1:23:40:50", 0, true}, // Invalid format
		{"", 0, true},           // Empty string
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			duration, err := ParseTime(tc.input)
			if tc.err {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if duration != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, duration)
				}
			}
		})
	}
}
