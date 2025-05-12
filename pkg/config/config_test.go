package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	
	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}
	
	// Check that PHPStan config is set
	if cfg.PHPStan.Args == nil || len(cfg.PHPStan.Args) == 0 {
		t.Error("PHPStan args should not be empty")
	}
	
	if !cfg.PHPStan.Enabled {
		t.Error("PHPStan should be enabled by default")
	}
	
	// Check that PHPCS config is set
	if cfg.PHPCS.Args == nil || len(cfg.PHPCS.Args) == 0 {
		t.Error("PHPCS args should not be empty")
	}
	
	if !cfg.PHPCS.Enabled {
		t.Error("PHPCS should be enabled by default")
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a test config
	testConfig := &Config{
		PHPStan: PHPStanConfig{
			Path:    "/test/path/to/phpstan",
			Args:    []string{"analyse", "--level=8"},
			Enabled: true,
		},
		PHPCS: PHPCSConfig{
			Path:    "/test/path/to/phpcs",
			Args:    []string{"--standard=PSR2"},
			Enabled: false,
		},
	}
	
	// Save the config
	configPath := filepath.Join(tempDir, "crixuamg.yaml")
	err = SaveConfig(testConfig, configPath)
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}
	
	// Check that the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("Config file was not created at %s", configPath)
	}
	
	// Set up the environment to load from the temp directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(oldWd)
	
	// Load the config
	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// Verify the loaded config matches what we saved
	if loadedConfig.PHPStan.Path != testConfig.PHPStan.Path {
		t.Errorf("PHPStan path mismatch: got %s, want %s", loadedConfig.PHPStan.Path, testConfig.PHPStan.Path)
	}
	
	if len(loadedConfig.PHPStan.Args) != len(testConfig.PHPStan.Args) {
		t.Errorf("PHPStan args length mismatch: got %d, want %d", len(loadedConfig.PHPStan.Args), len(testConfig.PHPStan.Args))
	} else {
		for i, arg := range testConfig.PHPStan.Args {
			if loadedConfig.PHPStan.Args[i] != arg {
				t.Errorf("PHPStan arg mismatch at index %d: got %s, want %s", i, loadedConfig.PHPStan.Args[i], arg)
			}
		}
	}
	
	if loadedConfig.PHPStan.Enabled != testConfig.PHPStan.Enabled {
		t.Errorf("PHPStan enabled mismatch: got %t, want %t", loadedConfig.PHPStan.Enabled, testConfig.PHPStan.Enabled)
	}
	
	if loadedConfig.PHPCS.Path != testConfig.PHPCS.Path {
		t.Errorf("PHPCS path mismatch: got %s, want %s", loadedConfig.PHPCS.Path, testConfig.PHPCS.Path)
	}
	
	if loadedConfig.PHPCS.Enabled != testConfig.PHPCS.Enabled {
		t.Errorf("PHPCS enabled mismatch: got %t, want %t", loadedConfig.PHPCS.Enabled, testConfig.PHPCS.Enabled)
	}
}
