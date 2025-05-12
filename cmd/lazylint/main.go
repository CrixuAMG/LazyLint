package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/crixuamg/pkg/config"
	"github.com/crixuamg/pkg/linters"
	"github.com/crixuamg/pkg/tui"
)

// Version information is defined in version.go

func main() {
	// Parse command line flags
	var (
		target       string
		configPath   string
		createConfig bool
		showVersion  bool
	)

	flag.StringVar(&target, "target", "", "Target file or directory to analyze")
	flag.StringVar(&configPath, "config", "", "Path to configuration file")
	flag.BoolVar(&createConfig, "create-config", false, "Create a default configuration file")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.Parse()

	// Show version information if requested
	if showVersion {
		fmt.Printf("LazyLint %s\n", GetVersion())
		os.Exit(0)
	}

	// Handle creating a default configuration file
	if createConfig {
		cfg := config.DefaultConfig()
		path := "lazylint.yaml"
		if configPath != "" {
			path = configPath
		}

		if err := config.SaveConfig(cfg, path); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating configuration file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Created configuration file at %s\n", path)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Create linter registry
	registry := linters.DefaultRegistry()

	// Configure linters from config
	for name, options := range cfg.Linters {
		linter, ok := registry.Get(name)
		if ok {
			linter.Configure(options)
		}
	}

	// Create and start the Bubble Tea program
	p := tea.NewProgram(
		tui.NewModel(cfg, registry),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
