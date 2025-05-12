package tools

import (
	"strings"
	"testing"
	"time"

	"github.com/crixuamg/pkg/config"
)

func TestParsePHPStanOutput(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Empty output",
			input:    "",
			expected: 0,
		},
		{
			name:     "Single line",
			input:    "Error: Something went wrong",
			expected: 1,
		},
		{
			name:     "Multiple lines",
			input:    "Error: Something went wrong\nWarning: This is a warning\nInfo: Just some info",
			expected: 3,
		},
		{
			name:     "Lines with whitespace",
			input:    "Error: Something went wrong\n\n\nWarning: This is a warning",
			expected: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ParsePHPStanOutput(tc.input)
			if len(result) != tc.expected {
				t.Errorf("Expected %d lines, got %d", tc.expected, len(result))
			}
		})
	}
}

func TestParsePHPCSOutput(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Empty output",
			input:    "",
			expected: 0,
		},
		{
			name:     "Single line",
			input:    "ERROR: Something went wrong",
			expected: 1,
		},
		{
			name:     "Multiple lines",
			input:    "ERROR: Something went wrong\nWARNING: This is a warning\nFILE: test.php",
			expected: 3,
		},
		{
			name:     "Lines with whitespace",
			input:    "ERROR: Something went wrong\n\n\nWARNING: This is a warning",
			expected: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ParsePHPCSOutput(tc.input)
			if len(result) != tc.expected {
				t.Errorf("Expected %d lines, got %d", tc.expected, len(result))
			}
		})
	}
}

func TestRunPHPStanDisabled(t *testing.T) {
	cfg := config.PHPStanConfig{
		Path:    "phpstan",
		Args:    []string{"analyse"},
		Enabled: false,
	}

	result, err := RunPHPStan(cfg, "test.php", 1*time.Second)
	if err != nil {
		t.Fatalf("RunPHPStan returned error: %v", err)
	}

	if !result.Success {
		t.Error("Result should be successful when PHPStan is disabled")
	}

	if !strings.Contains(result.Output, "disabled") {
		t.Errorf("Output should mention that PHPStan is disabled, got: %s", result.Output)
	}
}

func TestRunPHPCSDisabled(t *testing.T) {
	cfg := config.PHPCSConfig{
		Path:    "phpcs",
		Args:    []string{"--standard=PSR12"},
		Enabled: false,
	}

	result, err := RunPHPCS(cfg, "test.php", 1*time.Second)
	if err != nil {
		t.Fatalf("RunPHPCS returned error: %v", err)
	}

	if !result.Success {
		t.Error("Result should be successful when PHPCS is disabled")
	}

	if !strings.Contains(result.Output, "disabled") {
		t.Errorf("Output should mention that PHPCS is disabled, got: %s", result.Output)
	}
}
