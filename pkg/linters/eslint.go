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

// ESLint implements the Linter interface for ESLint
type ESLint struct {
	path    string
	args    []string
	enabled bool
}

// NewESLint creates a new ESLint linter
func NewESLint() *ESLint {
	// Try to find git root directory
	gitRoot, err := findGitRoot()
	if err != nil {
		gitRoot = ""
	}

	// Default path
	eslintPath := "eslint"

	// If git root is found, check for node_modules
	if gitRoot != "" {
		// Check if there's a local eslint in the project
		localPath := filepath.Join(gitRoot, "node_modules", ".bin", "eslint")
		if _, err := os.Stat(localPath); err == nil {
			eslintPath = localPath
		}
	}

	return &ESLint{
		path:    eslintPath,
		args:    []string{"--format=stylish"},
		enabled: true,
	}
}

// Name returns the name of the linter
func (l *ESLint) Name() string {
	return "eslint"
}

// Description returns a short description of the linter
func (l *ESLint) Description() string {
	return "JavaScript/TypeScript linter"
}

// Run executes the linter on the given target
func (l *ESLint) Run(ctx context.Context, target string) (*Result, error) {
	if !l.enabled {
		return &Result{
			Name:      l.Name(),
			Success:   true,
			Output:    "ESLint is disabled in configuration",
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
func (l *ESLint) IsAvailable() bool {
	_, err := exec.LookPath(l.path)
	return err == nil
}

// FileExtensions returns the file extensions this linter can process
func (l *ESLint) FileExtensions() []string {
	return []string{".js", ".jsx", ".ts", ".tsx"}
}

// Configure configures the linter with the given options
func (l *ESLint) Configure(options map[string]interface{}) error {
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
