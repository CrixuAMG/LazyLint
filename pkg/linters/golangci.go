package linters

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// GolangCI implements the Linter interface for golangci-lint
type GolangCI struct {
	path    string
	args    []string
	enabled bool
}

// NewGolangCI creates a new GolangCI linter
func NewGolangCI() *GolangCI {
	// Try to find git root directory
	gitRoot, err := findGitRoot()
	if err != nil {
		gitRoot = ""
	}

	// Default path
	golangciPath := "golangci-lint"

	// If git root is found, check for Go tools
	if gitRoot != "" {
		// Check if there's a local golangci-lint in the project
		localPath := filepath.Join(gitRoot, "bin", "golangci-lint")
		if _, err := os.Stat(localPath); err == nil {
			golangciPath = localPath
		}
	}

	return &GolangCI{
		path:    golangciPath,
		args:    []string{"run", "--out-format=colored-line-number"},
		enabled: true,
	}
}

// Name returns the name of the linter
func (l *GolangCI) Name() string {
	return "golangci-lint"
}

// Description returns a short description of the linter
func (l *GolangCI) Description() string {
	return "Fast Go linters runner"
}

// Run executes the linter on the given target
func (l *GolangCI) Run(ctx context.Context, target string) (*Result, error) {
	if !l.enabled {
		return &Result{
			Name:      l.Name(),
			Success:   true,
			Output:    "golangci-lint is disabled in configuration",
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
func (l *GolangCI) IsAvailable() bool {
	_, err := exec.LookPath(l.path)
	return err == nil
}

// FileExtensions returns the file extensions this linter can process
func (l *GolangCI) FileExtensions() []string {
	return []string{".go"}
}

// Configure configures the linter with the given options
func (l *GolangCI) Configure(options map[string]interface{}) error {
	if path, ok := options["path"].(string); ok {
		l.path = path
	}
	
	if args, ok := options["args"].([]string); ok {
		l.args = args
	} else if argsInterface, ok := options["args"].([]interface{}); ok {
		args := make([]string, len(argsInterface))
		for i, arg := range argsInterface {
			if str, ok := arg.(string); ok {
				args[i] = str
			}
		}
		l.args = args
	}
	
	if enabled, ok := options["enabled"].(bool); ok {
		l.enabled = enabled
	}
	
	return nil
}
