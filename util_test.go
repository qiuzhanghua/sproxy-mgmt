package main

import (
	"strings"
	"testing"
	"time"
)

func TestParseExpire(t *testing.T) {

	testCases := []struct {
		name          string
		input         string
		expected      string
		shouldError   bool
		errorContains string
	}{
		{
			name:        "special case -1",
			input:       "-1",
			expected:    "-1",
			shouldError: false,
		},
		{
			name:        "special case infinite",
			input:       "infinite",
			expected:    "-1",
			shouldError: false,
		},
		{
			name:        "special case forever",
			input:       "forever",
			expected:    "-1",
			shouldError: false,
		},
		{
			name:        "valid duration 5m",
			input:       "5m",
			expected:    time.Now().Add(5 * time.Minute).Format("20060102150405"),
			shouldError: false,
		},
		{
			name:        "valid duration 1h",
			input:       "1h",
			expected:    time.Now().Add(1 * time.Hour).Format("20060102150405"),
			shouldError: false,
		},
		{
			name:        "valid duration 2h30m",
			input:       "2h30m",
			expected:    time.Now().Add(2*time.Hour + 30*time.Minute).Format("20060102150405"),
			shouldError: false,
		},
		{
			name:          "invalid duration",
			input:         "abc",
			expected:      "",
			shouldError:   true,
			errorContains: `time: invalid duration "abc"`,
		},
		{
			name:        "valid timestamp",
			input:       "20231001120000",
			expected:    "20231001120000",
			shouldError: false,
		},
		{
			name:          "invalid timestamp format",
			input:         "2023-10-01 12:00:00",
			expected:      "",
			shouldError:   true,
			errorContains: `time: unknown unit "-" in duration "2023-10-01 12:00:00"`,
		},
		{
			name:          "invalid timestamp value",
			input:         "20231301120000",
			expected:      "",
			shouldError:   true,
			errorContains: `time: missing unit in duration "20231301120000"`,
		},
		{
			name:          "empty string",
			input:         "",
			expected:      "",
			shouldError:   true,
			errorContains: `time: invalid duration ""`,
		},
		{
			name:        "mixed case special case",
			input:       "InFiNiTe",
			expected:    "-1",
			shouldError: false,
		},
		{
			name:        "duration with spaces",
			input:       " 5m ",
			expected:    time.Now().Add(5 * time.Minute).Format("20060102150405"),
			shouldError: false,
		},
		{
			name:        "timestamp with spaces",
			input:       " 20231001120000 ",
			expected:    "20231001120000",
			shouldError: false,
		},
		{
			name:        "short timestamp",
			input:       "20231001",
			expected:    "20231001000000",
			shouldError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseExpire(tc.input)

			if tc.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("expected error containing '%s', but got '%s'", tc.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				} else if result != tc.expected {
					t.Errorf("expected '%s', but got '%s'", tc.expected, result)
				}
			}
		})
	}
}
