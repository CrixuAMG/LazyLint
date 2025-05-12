package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/crixuamg/pkg/linters"
)

// runLinter runs a linter and returns a command
func (m Model) runLinter(linter linters.Linter) tea.Cmd {
	return func() tea.Msg {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		// Run the linter
		result, err := linter.Run(ctx, m.target)
		if err != nil {
			return errorMsg{err: err.Error()}
		}
		return *result
	}
}

// updateViewportContent updates the viewport content based on the results
func (m *Model) updateViewportContent() {
	var content strings.Builder

	// Process each result
	for name, result := range m.results {
		content.WriteString(titleStyle.Render(fmt.Sprintf("%s Results", name)))
		content.WriteString("\n\n")

		statusLine := fmt.Sprintf("Completed in %.2fs", result.Duration.Seconds())
		if result.Success {
			content.WriteString(successStyle.Render(fmt.Sprintf("✓ %s completed successfully", name)))
		} else {
			content.WriteString(errorStyle.Render(fmt.Sprintf("✗ %s found issues", name)))
		}

		content.WriteString(" " + infoStyle.Render(statusLine) + "\n\n")

		lines := linters.ParseOutput(result)
		if len(lines) > 0 {
			for _, line := range lines {
				// Highlight error lines
				if strings.Contains(strings.ToLower(line), "error") {
					content.WriteString(errorStyle.Render(line))
				} else if strings.Contains(strings.ToLower(line), "warning") {
					content.WriteString(warningStyle.Render(line))
				} else {
					content.WriteString(line)
				}
				content.WriteString("\n")
			}
		} else {
			content.WriteString(infoStyle.Render(fmt.Sprintf("No output from %s", name)))
			content.WriteString("\n")
		}

		if result.Error != "" {
			content.WriteString("\n")
			content.WriteString(errorStyle.Render("Errors:"))
			content.WriteString("\n")
			content.WriteString(errorStyle.Render(result.Error))
			content.WriteString("\n")
		}

		content.WriteString("\n")
	}

	m.viewport.SetContent(content.String())
	m.viewport.GotoTop()
}
