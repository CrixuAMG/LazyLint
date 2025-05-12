package linters

import (
	"context"
	"time"
)

// Result represents the result of a linter execution
type Result struct {
	Name      string
	Success   bool
	Output    string
	Error     string
	Duration  time.Duration
	Timestamp time.Time
}

// Linter defines the interface that all linters must implement
type Linter interface {
	// Name returns the name of the linter
	Name() string
	
	// Description returns a short description of the linter
	Description() string
	
	// Run executes the linter on the given target
	Run(ctx context.Context, target string) (*Result, error)
	
	// IsAvailable checks if the linter is available in the current environment
	IsAvailable() bool
	
	// FileExtensions returns the file extensions this linter can process
	FileExtensions() []string
	
	// Configure configures the linter with the given options
	Configure(options map[string]interface{}) error
}

// Registry manages the available linters
type Registry struct {
	linters map[string]Linter
}

// NewRegistry creates a new linter registry
func NewRegistry() *Registry {
	return &Registry{
		linters: make(map[string]Linter),
	}
}

// Register adds a linter to the registry
func (r *Registry) Register(linter Linter) {
	r.linters[linter.Name()] = linter
}

// Get returns a linter by name
func (r *Registry) Get(name string) (Linter, bool) {
	linter, ok := r.linters[name]
	return linter, ok
}

// GetAll returns all registered linters
func (r *Registry) GetAll() []Linter {
	var result []Linter
	for _, linter := range r.linters {
		result = append(result, linter)
	}
	return result
}

// GetAvailable returns all available linters
func (r *Registry) GetAvailable() []Linter {
	var result []Linter
	for _, linter := range r.linters {
		if linter.IsAvailable() {
			result = append(result, linter)
		}
	}
	return result
}

// GetForExtension returns all linters that can process files with the given extension
func (r *Registry) GetForExtension(ext string) []Linter {
	var result []Linter
	for _, linter := range r.linters {
		for _, e := range linter.FileExtensions() {
			if e == ext {
				result = append(result, linter)
				break
			}
		}
	}
	return result
}
