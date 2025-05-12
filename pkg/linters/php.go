package linters

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// PHP implements the Linter interface for PHP syntax checking
type PHP struct {
	path    string
	args    []string
	enabled bool
}

// NewPHP creates a new PHP linter
func NewPHP() *PHP {
	// Default path
	phpPath := "php"

	return &PHP{
		path:    phpPath,
		args:    []string{"-l"},
		enabled: true,
	}
}

// Name returns the name of the linter
func (l *PHP) Name() string {
	return "php"
}

// Description returns a short description of the linter
func (l *PHP) Description() string {
	return "PHP syntax checker"
}

// Configure configures the linter with the given options
func (l *PHP) Configure(options map[string]interface{}) error {
	if path, ok := options["path"].(string); ok {
		l.path = path
	}

	if args, ok := options["args"].([]interface{}); ok {
		l.args = make([]string, 0, len(args))
		for _, arg := range args {
			if s, ok := arg.(string); ok {
				l.args = append(l.args, s)
			}
		}
	}

	if enabled, ok := options["enabled"].(bool); ok {
		l.enabled = enabled
	}

	return nil
}

// Run executes the linter on the given target
func (l *PHP) Run(ctx context.Context, target string) (*Result, error) {
	if !l.enabled {
		return &Result{
			Name:      l.Name(),
			Success:   true,
			Output:    "PHP syntax check is disabled in configuration",
			Timestamp: time.Now(),
		}, nil
	}

	args := append([]string{}, l.args...)
	if target != "" {
		args = append(args, target)
	}

	start := time.Now()
	cmd := exec.CommandContext(ctx, l.path, args...)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(start)

	result := &Result{
		Name:      l.Name(),
		Output:    stdout.String(),
		Error:     stderr.String(),
		Duration:  duration,
		Timestamp: time.Now(),
	}

	if err != nil {
		// Check if it's a timeout
		if ctx.Err() == context.DeadlineExceeded {
			return result, fmt.Errorf("command timed out after %s", duration)
		}

		// Check if it's an exit code error (which is expected for these tools when they find issues)
		if _, ok := err.(*exec.ExitError); ok {
			result.Success = false
			return result, nil
		}

		return result, fmt.Errorf("command failed: %w", err)
	}

	result.Success = true
	return result, nil
}

// IsAvailable checks if the linter is available
func (l *PHP) IsAvailable() bool {
	_, err := exec.LookPath(l.path)
	return err == nil
}

// FileExtensions returns the file extensions this linter can process
func (l *PHP) FileExtensions() []string {
	return []string{".php"}
}
