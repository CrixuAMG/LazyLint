package config

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FindGitRoot finds the root directory of the git repository
func FindGitRoot() (string, error) {
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
