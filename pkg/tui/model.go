package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/crixuamg/pkg/config"
	"github.com/crixuamg/pkg/linters"
)

// NewModel creates a new model with the given configuration and linter registry
func NewModel(cfg *config.Config, registry *linters.Registry) Model {
	// Initialize theme from config
	InitTheme(cfg)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	vp := viewport.New(80, 24)
	vp.Style = lipgloss.NewStyle().Margin(1, 2)

	// Create a default explorer with placeholder dimensions
	// (will be resized on first WindowSizeMsg)
	explorer := NewExplorer(80, 24)

	// Set the registry in the explorer
	explorer.registry = registry

	// Get available linters
	activeLinters := registry.GetAvailable()

	// Create tool names for the tools pane
	var toolNames []string
	for _, linter := range activeLinters {
		toolNames = append(toolNames, linter.Name())
	}

	// Create panes with placeholder dimensions
	// (will be resized on first WindowSizeMsg)
	panes := []Pane{
		NewRepoInfoPane(80, 8),
		NewToolsPane(80, 16, toolNames),
		NewExplorerPane(80, 24),
		NewOutputPane(80, 24),
	}

	return Model{
		config:        cfg,
		registry:      registry,
		state:         StateMultiPane, // Start with the multi-pane layout as default
		selectedTool:  0,
		results:       make(map[string]*linters.Result),
		viewport:      vp,
		spinner:       s,
		help:          help.New(),
		keys:          keys,
		explorer:      explorer,
		activeLinters: activeLinters,
		panes:         panes,
		activePaneIndex: 2, // Start with the file explorer pane active
		activeTab:     0,   // Start with the Explorer tab
		useNewUI:      true, // Always use the new UI
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
	)
}

// Update updates the model based on messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keybindings first
		switch msg.String() {
		case "q":
			// Quit the application
			return m, tea.Quit

		case "?":
			// Toggle help
			m.help.ShowAll = !m.help.ShowAll
			return m, nil

		case "tab":
			// Next tab
			m.activeTab = (m.activeTab + 1) % 4
			return m, nil

		case "shift+tab":
			// Previous tab
			m.activeTab = (m.activeTab - 1 + 4) % 4
			return m, nil

		case "t":
			// Cycle through themes
			var themes []string
			for name := range Themes {
				themes = append(themes, name)
			}

			// Find current theme index
			currentIndex := 0
			for i, name := range themes {
				if name == m.config.UI.Theme {
					currentIndex = i
					break
				}
			}

			// Switch to next theme
			nextIndex := (currentIndex + 1) % len(themes)
			m.config.UI.Theme = themes[nextIndex]
			ApplyTheme(m.config.UI.Theme)
			return m, nil
		}

		// Handle tab-specific keybindings
		switch m.activeTab {
		case 0: // Explorer tab
			// Update explorer
			var explorerCmd tea.Cmd
			m.explorer, explorerCmd = m.explorer.Update(msg)

			// Handle running tools on selected files
			if msg.String() == "r" {
				selectedFiles := m.explorer.GetSelectedFiles()
				if len(selectedFiles) > 0 {
					m.target = strings.Join(selectedFiles, " ")
					m.state = StateRunning

					// Run all available linters on the selected files
					var cmds []tea.Cmd
					for _, linter := range m.activeLinters {
						cmds = append(cmds, m.runLinter(linter))
					}
					return m, tea.Batch(cmds...)
				}
			}

			return m, explorerCmd

		case 1: // Linters tab
			// Handle navigation
			switch msg.String() {
			case "up", "k":
				if m.selectedTool > 0 {
					m.selectedTool--
				}
			case "down", "j":
				if m.selectedTool < len(m.activeLinters)-1 {
					m.selectedTool++
				}
			case "enter":
				// Run selected linter
				if m.selectedTool < len(m.activeLinters) {
					m.state = StateRunning
					return m, m.runLinter(m.activeLinters[m.selectedTool])
				}
			}
			return m, nil

		case 2: // Results tab
			// Update viewport for scrolling
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd

		case 3: // Config tab
			// No special handling needed yet
			return m, nil
		}

		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 10
		m.help.Width = msg.Width

		// Update explorer dimensions for the new UI
		m.explorer.width = msg.Width - 10
		m.explorer.height = msg.Height - 15
		var cmd tea.Cmd
		m.explorer, cmd = m.explorer.Update(msg)
		cmds = append(cmds, cmd)

		// Update pane dimensions for multi-pane layout (for backward compatibility)
		// Left column (30% of width)
		leftColWidth := msg.Width * 3 / 10
		// Right column (70% of width)
		rightColWidth := msg.Width - leftColWidth

		// Reserve space for status and shortcut bars (2 lines)
		contentHeight := msg.Height - 3 // Add extra space at the top

		// Top panes in left column
		repoInfoHeight := 6 // Reduced height for repo info
		toolsHeight := contentHeight - repoInfoHeight

		// Set pane sizes
		if len(m.panes) >= 4 {
			// Repo info pane (top-left)
			m.panes[0].SetSize(leftColWidth, repoInfoHeight)

			// Tools pane (bottom-left)
			m.panes[1].SetSize(leftColWidth, toolsHeight)

			// Explorer pane (right, when active)
			m.panes[2].SetSize(rightColWidth, contentHeight)

			// Output pane (right, when active)
			m.panes[3].SetSize(rightColWidth, contentHeight)
		}

	case spinner.TickMsg:
		if m.state == StateRunning {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case linters.Result:
		m.results[msg.Name] = &msg

		// Check if all expected linters have completed
		expectedCount := 1
		if m.selectedTool == len(m.activeLinters) { // File explorer mode with multiple files
			expectedCount = len(m.activeLinters)
		}

		if len(m.results) >= expectedCount {
			// All linters have completed, show results
			m.state = StateResults
			m.activeTab = 2 // Switch to Results tab

			// Combine results
			var content strings.Builder
			for _, result := range m.results {
				content.WriteString(fmt.Sprintf("=== %s ===\n", result.Name))
				content.WriteString(result.Output)
				content.WriteString("\n\n")
			}

			// Set viewport content
			m.viewport.SetContent(content.String())
			m.viewport.GotoTop()

			// Update output pane for backward compatibility
			if len(m.panes) >= 4 {
				outputPane, ok := m.panes[3].(*OutputPane)
				if ok {
					outputPane.SetContent(content.String())
					outputPane.SetTitle("Linter Results")
				}
			}
		}
	}

	if len(cmds) > 0 {
		return m, tea.Batch(cmds...)
	}
	return m, cmd
}


