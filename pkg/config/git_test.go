package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindGitRoot(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "git-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a fake .git directory
	gitDir := filepath.Join(tempDir, ".git")
	err = os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}
	
	// Create a subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	
	// Create a nested subdirectory
	nestedDir := filepath.Join(subDir, "nested")
	err = os.Mkdir(nestedDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested subdirectory: %v", err)
	}
	
	// Save current directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldWd)
	
	// Test from the root directory
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	root, err := FindGitRoot()
	if err != nil {
		t.Fatalf("FindGitRoot failed: %v", err)
	}
	
	if root != tempDir {
		t.Errorf("FindGitRoot from root dir returned wrong path: got %s, want %s", root, tempDir)
	}
	
	// Test from the subdirectory
	err = os.Chdir(subDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	root, err = FindGitRoot()
	if err != nil {
		t.Fatalf("FindGitRoot failed: %v", err)
	}
	
	if root != tempDir {
		t.Errorf("FindGitRoot from subdir returned wrong path: got %s, want %s", root, tempDir)
	}
	
	// Test from the nested subdirectory
	err = os.Chdir(nestedDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	root, err = FindGitRoot()
	if err != nil {
		t.Fatalf("FindGitRoot failed: %v", err)
	}
	
	if root != tempDir {
		t.Errorf("FindGitRoot from nested dir returned wrong path: got %s, want %s", root, tempDir)
	}
}
