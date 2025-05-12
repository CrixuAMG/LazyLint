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

// PHPStan implements the Linter interface for PHPStan
type PHPStan struct {
	path    string
	args    []string
	enabled bool
}

// NewPHPStan creates a new PHPStan linter
func NewPHPStan() *PHPStan {
	// Try to find git root directory
	gitRoot, err := findGitRoot()
	if err != nil {
		gitRoot = ""
	}

	// Default path
	phpstanPath := "phpstan"

	// If git root is found, check vendor/bin directory
	if gitRoot != "" {
		vendorBin := filepath.Join(gitRoot, "vendor", "bin")
		vendorPath := filepath.Join(vendorBin, "phpstan")
		if _, err := os.Stat(vendorPath); err == nil {
			phpstanPath = vendorPath
		}
	}

	return &PHPStan{
		path:    phpstanPath,
		args:    []string{"analyse", "--level=5"},
		enabled: true,
	}
}

// Name returns the name of the linter
func (l *PHPStan) Name() string {
	return "phpstan"
}

// Description returns a short description of the linter
func (l *PHPStan) Description() string {
	return "PHP Static Analysis Tool"
}

// Run executes the linter on the given target
func (l *PHPStan) Run(ctx context.Context, target string) (*Result, error) {
	if !l.enabled {
		return &Result{
			Name:      l.Name(),
			Success:   true,
			Output:    "PHPStan is disabled in configuration",
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
func (l *PHPStan) IsAvailable() bool {
	_, err := exec.LookPath(l.path)
	return err == nil
}

// FileExtensions returns the file extensions this linter can process
func (l *PHPStan) FileExtensions() []string {
	return []string{".php"}
}

// Configure configures the linter with the given options
func (l *PHPStan) Configure(options map[string]interface{}) error {
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

// Helper function to find git root
func findGitRoot() (string, error) {
	// Try to find git root using git command
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)), nil
	}

	// Fallback: manually search for .git directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the root directory
			break
		}
		dir = parent
	}

	// If we can't find a git root, return the current directory
	return os.Getwd()
}
