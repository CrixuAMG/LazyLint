package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/crixuamg/pkg/config"
	"github.com/crixuamg/pkg/linters"
)

// State represents the current state of the application
type State int

const (
	StateMenu State = iota
	StateRunning
	StateResults
	StateConfig
	StateExplorer
	StateMultiPane
)

// errorMsg is a message that contains an error
type errorMsg struct {
	err string
}

// Error returns the error string
func (e errorMsg) Error() string {
	return e.err
}

// Model represents the application state
type Model struct {
	config      *config.Config
	registry    *linters.Registry
	state       State
	width       int
	height      int
	selectedTool int
	target      string
	results     map[string]*linters.Result
	viewport    viewport.Model
	spinner     spinner.Model
	help        help.Model
	keys        keyMap
	err         string
	explorer    *Explorer
	activeLinters []linters.Linter

	// Multi-pane layout
	panes       []Pane
	activePaneIndex int

	// New UI layout
	activeTab   int // 0: Explorer, 1: Linters, 2: Results, 3: Config
	useNewUI    bool // Whether to use the new UI
}
