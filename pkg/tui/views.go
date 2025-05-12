package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the current view
func (m Model) View() string {
	return m.newView()
}

// renderLogo renders the LazyLint logo
func renderLogo() string {
	logo := "LazyLint"
	return logoStyle.Render(logo)
}

// renderTabs renders the tab bar
func (m Model) renderTabs() string {
	tabs := []string{"Explorer", "Linters", "Results", "Config"}
	renderedTabs := []string{}

	for i, tab := range tabs {
		if i == m.activeTab {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(tab))
		} else {
			renderedTabs = append(renderedTabs, inactiveTabStyle.Render(tab))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

// renderTabContent renders the content for the active tab
func (m Model) renderTabContent() string {
	var content string
	width := m.width - 4 // Account for margins

	switch m.activeTab {
	case 0: // Explorer tab
		content = m.renderExplorerTab(width)
	case 1: // Linters tab
		content = m.renderLintersTab(width)
	case 2: // Results tab
		content = m.renderResultsTab(width)
	case 3: // Config tab
		content = m.renderConfigTab(width)
	}

	return tabContentStyle.Width(width).Render(content)
}

// renderExplorerTab renders the explorer tab content
func (m Model) renderExplorerTab(width int) string {
	// Adjust explorer dimensions
	m.explorer.width = width - 4 // Account for padding
	m.explorer.height = m.height - 10 // Account for other UI elements

	// Get explorer view
	explorerView := m.explorer.View()

	// Add title
	title := titleStyle.Render("File Explorer")

	// Add path info
	pathInfo := infoStyle.Render(fmt.Sprintf("Current Directory: %s", m.explorer.currentDir))

	// Join all components
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		pathInfo,
		explorerView,
	)
}

// renderLintersTab renders the linters tab content
func (m Model) renderLintersTab(width int) string {
	// Add title
	title := titleStyle.Render("Available Linters")

	// Build linter list
	var lintersContent strings.Builder
	for i, linter := range m.activeLinters {
		if i == m.selectedTool {
			lintersContent.WriteString(selectedItemStyle.Render(fmt.Sprintf("> %s", linter.Name())))
		} else {
			lintersContent.WriteString(itemStyle.Render(fmt.Sprintf("  %s", linter.Name())))
		}
		lintersContent.WriteString("\n")
	}

	if len(m.activeLinters) == 0 {
		lintersContent.WriteString(infoStyle.Render("No linters available"))
	}

	// Add action buttons
	buttons := lipgloss.JoinHorizontal(
		lipgloss.Top,
		buttonStyle.Render("Run Selected"),
		buttonStyle.Render("Run All"),
	)

	// Join all components
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		lintersContent.String(),
		"\n",
		buttons,
	)
}

// renderResultsTab renders the results tab content
func (m Model) renderResultsTab(width int) string {
	// Add title
	title := titleStyle.Render("Linter Results")

	// Check if we have results
	if len(m.results) == 0 {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			infoStyle.Render("No results yet. Run a linter to see results here."),
		)
	}

	// Build results content
	var resultsContent strings.Builder
	for name, result := range m.results {
		resultsContent.WriteString(subtitleStyle.Render(fmt.Sprintf("%s Results", name)))
		resultsContent.WriteString("\n\n")

		if result.Output != "" {
			resultsContent.WriteString(result.Output)
		} else {
			resultsContent.WriteString(infoStyle.Render("No output from linter"))
		}

		resultsContent.WriteString("\n\n")
	}

	// Join all components
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		resultsContent.String(),
	)
}

// renderConfigTab renders the config tab content
func (m Model) renderConfigTab(width int) string {
	// Add title
	title := titleStyle.Render("Configuration")

	// Build config content
	var configContent strings.Builder

	// Theme selection
	configContent.WriteString(subtitleStyle.Render("Theme Settings"))
	configContent.WriteString("\n")
	configContent.WriteString(fmt.Sprintf("Current Theme: %s\n", m.config.UI.Theme))
	configContent.WriteString("Available Themes:\n")

	// Get all theme names and sort them alphabetically
	var themeNames []string
	for name := range m.config.UI.Themes {
		themeNames = append(themeNames, name)
	}

	// Sort theme names alphabetically
	sort.Strings(themeNames)

	// Display themes with consistent indentation
	for _, name := range themeNames {
		if name == m.config.UI.Theme {
			configContent.WriteString(selectedItemStyle.Render(fmt.Sprintf("  > %s\n", name)))
		} else {
			configContent.WriteString(itemStyle.Render(fmt.Sprintf("    %s\n", name)))
		}
	}

	// Join all components
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		configContent.String(),
	)
}

// renderStatusBar renders the status bar
func (m Model) renderStatusBar() string {
	// Create status text based on current state
	var statusText string

	switch m.state {
	case StateRunning:
		statusText = fmt.Sprintf("%s Running linters...", m.spinner.View())
	case StateResults:
		statusText = "Linter results ready"
	default:
		statusText = fmt.Sprintf("LazyLint - Theme: %s", m.config.UI.Theme)
	}

	return statusBarStyle.Width(m.width).Render(statusText)
}

// renderHelpBar renders the help bar with keyboard shortcuts
func (m Model) renderHelpBar() string {
	var shortcuts string

	// Base shortcuts that are always available
	baseShortcuts := "q: Quit • ?: Help • t: Change theme"

	// Add tab-specific shortcuts
	switch m.activeTab {
	case 0: // Explorer tab
		shortcuts = baseShortcuts + " • ↑/↓: Navigate • Space: Select file • Enter: Open • r: Run tools"
	case 1: // Linters tab
		shortcuts = baseShortcuts + " • ↑/↓: Navigate • Enter: Run selected linter"
	case 2: // Results tab
		shortcuts = baseShortcuts + " • ↑/↓: Scroll results"
	case 3: // Config tab
		shortcuts = baseShortcuts + " • ↑/↓: Navigate"
	}

	// Add tab navigation shortcuts
	shortcuts += " • Tab: Next tab • Shift+Tab: Previous tab"

	return helpBarStyle.Width(m.width).Render(shortcuts)
}

// newView renders the new UI layout
func (m Model) newView() string {
	// Create a background style with the theme's background color
	backgroundStyle := lipgloss.NewStyle().
		Background(background).
		Width(m.width).
		Height(m.height)

	// Render the logo
	logo := renderLogo()

	// Render the tabs
	tabs := m.renderTabs()

	// Render the tab content
	content := m.renderTabContent()

	// Render the status bar
	statusBar := m.renderStatusBar()

	// Render the help bar
	helpBar := m.renderHelpBar()

	// Create a container for the main content
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		logo,
		tabs,
		content,
	)

	// Create the final UI
	ui := lipgloss.JoinVertical(
		lipgloss.Left,
		mainContent,
		lipgloss.NewStyle().Height(m.height - lipgloss.Height(mainContent) - 2).Render(""),
		helpBar,
		statusBar,
	)

	// Apply the background color to the entire view
	return backgroundStyle.Render(ui)
}
