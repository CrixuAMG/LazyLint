package linters

import "strings"

// DefaultRegistry creates a registry with all default linters
func DefaultRegistry() *Registry {
	registry := NewRegistry()

	// Register PHP linters
	registry.Register(NewPHPStan())
	registry.Register(NewPHPCS())
	registry.Register(NewPHP())

	// Register Go linters
	registry.Register(NewGolangCI())

	// Register JS/TS linters
	registry.Register(NewESLint())

	return registry
}

// ParseOutput parses the output of a linter
func ParseOutput(result *Result) []string {
	if result == nil || result.Output == "" {
		return []string{}
	}

	// Split the output into lines
	lines := strings.Split(result.Output, "\n")

	// Remove empty lines
	var filtered []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			filtered = append(filtered, line)
		}
	}

	return filtered
}
