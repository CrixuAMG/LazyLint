package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// ThemeColors holds the colors for a theme
type ThemeColors struct {
	Subtle     string `mapstructure:"subtle"`
	Highlight  string `mapstructure:"highlight"`
	Special    string `mapstructure:"special"`
	Error      string `mapstructure:"error"`
	Warning    string `mapstructure:"warning"`
	Border     string `mapstructure:"border"`
	Text       string `mapstructure:"text"`
	DimText    string `mapstructure:"dim_text"`
	Background string `mapstructure:"background"`
}

// ThemeConfig holds the theme configuration
type ThemeConfig struct {
	Name   string      `mapstructure:"name"`
	Colors ThemeColors `mapstructure:"colors"`
}

// UIConfig holds the UI configuration
type UIConfig struct {
	Theme   string                 `mapstructure:"theme"`
	Themes  map[string]ThemeConfig `mapstructure:"themes"`
}

// Config holds the application configuration
type Config struct {
	Linters map[string]map[string]interface{} `mapstructure:"linters"`
	UI      UIConfig                         `mapstructure:"ui"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Linters: map[string]map[string]interface{}{
			"phpstan": {
				"path":    "phpstan",
				"args":    []string{"analyse", "--level=5"},
				"enabled": true,
			},
			"phpcs": {
				"path":    "phpcs",
				"args":    []string{"--standard=PSR12"},
				"enabled": true,
			},
			"php": {
				"path":    "php",
				"args":    []string{"-l"},
				"enabled": true,
			},
			"golangci-lint": {
				"path":    "golangci-lint",
				"args":    []string{"run", "--out-format=colored-line-number"},
				"enabled": true,
			},
			"eslint": {
				"path":    "eslint",
				"args":    []string{"--format=stylish"},
				"enabled": true,
			},
		},
		UI: UIConfig{
			Theme: "tokyo-night",
			Themes: map[string]ThemeConfig{
				"tokyo-night": {
					Name: "tokyo-night",
					Colors: ThemeColors{
						Subtle:     "#565F89",
						Highlight:  "#7AA2F7",
						Special:    "#9ECE6A",
						Error:      "#F7768E",
						Warning:    "#E0AF68",
						Border:     "#565F89",
						Text:       "#C0CAF5",
						DimText:    "#9AA5CE",
						Background: "#1A1B26",
					},
				},
				"light": {
					Name: "light",
					Colors: ThemeColors{
						Subtle:     "#D9DCCF",
						Highlight:  "#874BFD",
						Special:    "#43BF6D",
						Error:      "#FF5F87",
						Warning:    "#FFA500",
						Border:     "#4A4A4A",
						Text:       "#1A1A1A",
						DimText:    "#666666",
						Background: "#FFFFFF",
					},
				},
				"high-contrast": {
					Name: "high-contrast",
					Colors: ThemeColors{
						Subtle:     "#000000",
						Highlight:  "#0000FF",
						Special:    "#008000",
						Error:      "#FF0000",
						Warning:    "#FFA500",
						Border:     "#000000",
						Text:       "#000000",
						DimText:    "#333333",
						Background: "#FFFFFF",
					},
				},
				"catppuccin": {
					Name: "catppuccin",
					Colors: ThemeColors{
						Subtle:     "#6E6C7E",
						Highlight:  "#89B4FA", // Blue
						Special:    "#A6E3A1", // Green
						Error:      "#F38BA8", // Red
						Warning:    "#FAB387", // Peach
						Border:     "#89B4FA", // Blue
						Text:       "#CDD6F4", // Text
						DimText:    "#BAC2DE", // Subtext0
						Background: "#1E1E2E", // Base
					},
				},
			},
		},
	}
}

// LoadConfig loads the configuration from a file
func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	// Set up viper
	v := viper.New()
	v.SetConfigName("lazylint")
	v.SetConfigType("yaml")

	// Look for config in current directory and all parent directories
	dir, err := os.Getwd()
	if err != nil {
		return config, fmt.Errorf("failed to get current directory: %w", err)
	}

	// First check for user config in $HOME/.config/lazylint/
	homeDir, err := os.UserHomeDir()
	if err == nil {
		userConfigDir := filepath.Join(homeDir, ".config", "lazylint")
		v.AddConfigPath(userConfigDir)
	}

	// Then look in current directory and all parent directories
	for {
		v.AddConfigPath(dir)
		if _, err := os.Stat(filepath.Join(dir, "lazylint.yaml")); err == nil {
			break
		}
		if _, err := os.Stat(filepath.Join(dir, "lazylint.yml")); err == nil {
			break
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the root directory
			break
		}
		dir = parent
	}

	// Try to read the config file
	if err := v.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		// Config file found, unmarshal it
		if err := v.Unmarshal(config); err != nil {
			return config, fmt.Errorf("failed to unmarshal config: %w", err)
		}
	}

	return config, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *Config, path string) error {
	v := viper.New()
	v.SetConfigFile(path)

	// Set the config values
	v.Set("linters", config.Linters)
	v.Set("ui", config.UI)

	// Save the config
	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// SaveUserConfig saves the configuration to the user's config directory
func SaveUserConfig(config *Config) error {
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create config directory if it doesn't exist
	userConfigDir := filepath.Join(homeDir, ".config", "lazylint")
	if err := os.MkdirAll(userConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Save config to file
	configPath := filepath.Join(userConfigDir, "config.yaml")
	return SaveConfig(config, configPath)
}
