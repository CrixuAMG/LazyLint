package tools

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/crixuamg/pkg/config"
)

// Result represents the result of a tool execution
type Result struct {
	Tool      string
	Success   bool
	Output    string
	Error     string
	Duration  time.Duration
	Timestamp time.Time
}

// RunPHPStan runs phpstan with the given configuration and target
func RunPHPStan(cfg config.PHPStanConfig, target string, timeout time.Duration) (*Result, error) {
	if !cfg.Enabled {
		return &Result{
			Tool:      "phpstan",
			Success:   true,
			Output:    "PHPStan is disabled in configuration",
			Timestamp: time.Now(),
		}, nil
	}

	args := append([]string{}, cfg.Args...)
	if target != "" {
		args = append(args, target)
	}

	return runTool("phpstan", cfg.Path, args, timeout)
}

// RunPHPCS runs phpcs with the given configuration and target
func RunPHPCS(cfg config.PHPCSConfig, target string, timeout time.Duration) (*Result, error) {
	if !cfg.Enabled {
		return &Result{
			Tool:      "phpcs",
			Success:   true,
			Output:    "PHPCS is disabled in configuration",
			Timestamp: time.Now(),
		}, nil
	}

	args := append([]string{}, cfg.Args...)
	if target != "" {
		args = append(args, target)
	}

	return runTool("phpcs", cfg.Path, args, timeout)
}

// runTool runs a tool with the given arguments and timeout
func runTool(name, path string, args []string, timeout time.Duration) (*Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	cmd := exec.CommandContext(ctx, path, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	duration := time.Since(start)

	result := &Result{
		Tool:      name,
		Output:    stdout.String(),
		Error:     stderr.String(),
		Duration:  duration,
		Timestamp: time.Now(),
	}

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return result, fmt.Errorf("command timed out after %s", timeout)
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

// ParsePHPStanOutput parses the output of phpstan
func ParsePHPStanOutput(output string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

// ParsePHPCSOutput parses the output of phpcs
func ParsePHPCSOutput(output string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines
}
