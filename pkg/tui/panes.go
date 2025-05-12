package tui

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PaneType represents the type of pane
type PaneType int

const (
	PaneRepoInfo PaneType = iota
	PaneTools
	PaneExplorer
	PaneOutput
)

// Pane is the interface for all panes
type Pane interface {
	Update(msg tea.Msg) (Pane, tea.Cmd)
	View() string
	SetSize(width, height int)
	GetType() PaneType
	SetActive(active bool)
}

// RepoInfoPane displays information about the repository
type RepoInfoPane struct {
	width  int
	height int
	repoName string
	repoPath string
	isActive bool
}

// NewRepoInfoPane creates a new repo info pane
func NewRepoInfoPane(width, height int) *RepoInfoPane {
	// Get git repo name
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	repoPath := "."
	if err == nil {
		repoPath = strings.TrimSpace(string(out))
	}

	// Get repo name from path
	repoName := filepath.Base(repoPath)

	return &RepoInfoPane{
		width:  width,
		height: height,
		repoName: repoName,
		repoPath: repoPath,
	}
}

// Update updates the repo info pane
func (p *RepoInfoPane) Update(msg tea.Msg) (Pane, tea.Cmd) {
	return p, nil
}

// View renders the repo info pane
func (p *RepoInfoPane) View() string {
	// Use the appropriate style based on whether this pane is active
	var style lipgloss.Style
	if p.isActive {
		style = activePaneStyle.Copy()
	} else {
		style = inactivePaneStyle.Copy()
	}

	style = style.Width(p.width - 4).
		Height(p.height - 2)

	// Add title to the border
	style = style.Border(lipgloss.RoundedBorder())

	// Add title to the top border with pane index
	title := "1: Repository"

	content := fmt.Sprintf("Repository: %s\nPath: %s", p.repoName, p.repoPath)
	renderedContent := style.Render(content)

	// Add title to the top border
	lines := strings.Split(renderedContent, "\n")
	if len(lines) > 0 {
		// Find the position to insert the title
		firstLine := lines[0]

		// Ensure title doesn't exceed pane width
		maxTitleLen := len(firstLine) - 4 // Leave some space for borders
		if len(title) > maxTitleLen && maxTitleLen > 5 {
			// Truncate title if it's too long
			title = title[:maxTitleLen-3] + "..."
		}

		titlePos := (len(firstLine) - len(title)) / 2
		if titlePos < 0 {
			titlePos = 0
		}

		// Insert the title
		if titlePos+len(title) <= len(firstLine) {
			lines[0] = firstLine[:titlePos] + title + firstLine[titlePos+len(title):]
		}

		// Rejoin the lines
		return strings.Join(lines, "\n")
	}

	return renderedContent
}

// SetActive sets whether this pane is active
func (p *RepoInfoPane) SetActive(active bool) {
	p.isActive = active
}

// SetSize sets the size of the pane
func (p *RepoInfoPane) SetSize(width, height int) {
	p.width = width
	p.height = height
}

// GetType returns the type of pane
func (p *RepoInfoPane) GetType() PaneType {
	return PaneRepoInfo
}

// ToolsPane displays available tools
type ToolsPane struct {
	width  int
	height int
	tools  []string
	selected int
	isActive bool
}

// NewToolsPane creates a new tools pane
func NewToolsPane(width, height int, tools []string) *ToolsPane {
	return &ToolsPane{
		width:  width,
		height: height,
		tools:  tools,
		selected: 0,
	}
}

// Update updates the tools pane
func (p *ToolsPane) Update(msg tea.Msg) (Pane, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if p.selected > 0 {
				p.selected--
			}
		case "down", "j":
			if p.selected < len(p.tools)-1 {
				p.selected++
			}
		}
	}
	return p, nil
}

// View renders the tools pane
func (p *ToolsPane) View() string {
	// Use the appropriate style based on whether this pane is active
	var style lipgloss.Style
	if p.isActive {
		style = activePaneStyle.Copy()
	} else {
		style = inactivePaneStyle.Copy()
	}

	style = style.Width(p.width - 4).
		Height(p.height - 2)

	// Add title to the border
	style = style.Border(lipgloss.RoundedBorder())

	// Add title to the top border with pane index
	title := "2: Tools"

	var content strings.Builder

	for i, tool := range p.tools {
		if i == p.selected {
			content.WriteString(selectedItemStyle.Render(fmt.Sprintf("> %s", tool)))
		} else {
			content.WriteString(itemStyle.Render(fmt.Sprintf("  %s", tool)))
		}
		content.WriteString("\n")
	}

	renderedContent := style.Render(content.String())

	// Add title to the top border
	lines := strings.Split(renderedContent, "\n")
	if len(lines) > 0 {
		// Find the position to insert the title
		firstLine := lines[0]

		// Ensure title doesn't exceed pane width
		maxTitleLen := len(firstLine) - 4 // Leave some space for borders
		if len(title) > maxTitleLen && maxTitleLen > 5 {
			// Truncate title if it's too long
			title = title[:maxTitleLen-3] + "..."
		}

		titlePos := (len(firstLine) - len(title)) / 2
		if titlePos < 0 {
			titlePos = 0
		}

		// Insert the title
		if titlePos+len(title) <= len(firstLine) {
			lines[0] = firstLine[:titlePos] + title + firstLine[titlePos+len(title):]
		}

		// Rejoin the lines
		return strings.Join(lines, "\n")
	}

	return renderedContent
}

// SetActive sets whether this pane is active
func (p *ToolsPane) SetActive(active bool) {
	p.isActive = active
}

// SetSize sets the size of the pane
func (p *ToolsPane) SetSize(width, height int) {
	p.width = width
	p.height = height
}

// GetType returns the type of pane
func (p *ToolsPane) GetType() PaneType {
	return PaneTools
}

// GetSelectedTool returns the selected tool
func (p *ToolsPane) GetSelectedTool() string {
	if p.selected >= 0 && p.selected < len(p.tools) {
		return p.tools[p.selected]
	}
	return ""
}

// OutputPane displays the output of a command
type OutputPane struct {
	width    int
	height   int
	viewport viewport.Model
	content  string
	title    string
	isActive bool
}

// NewOutputPane creates a new output pane
func NewOutputPane(width, height int) *OutputPane {
	vp := viewport.New(width-4, height-4)
	vp.Style = lipgloss.NewStyle()

	return &OutputPane{
		width:    width,
		height:   height,
		viewport: vp,
		content:  "No output to display",
		title:    "Output",
	}
}

// Update updates the output pane
func (p *OutputPane) Update(msg tea.Msg) (Pane, tea.Cmd) {
	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	return p, cmd
}

// View renders the output pane
func (p *OutputPane) View() string {
	// Use the appropriate style based on whether this pane is active
	var style lipgloss.Style
	if p.isActive {
		style = activePaneStyle.Copy()
	} else {
		style = inactivePaneStyle.Copy()
	}

	style = style.Width(p.width - 4).
		Height(p.height - 2)

	// Add title to the border
	style = style.Border(lipgloss.RoundedBorder())

	// Add title to the top border with pane index
	title := "4: Output"

	renderedContent := style.Render(p.viewport.View())

	// Add title to the top border
	lines := strings.Split(renderedContent, "\n")
	if len(lines) > 0 {
		// Find the position to insert the title
		firstLine := lines[0]

		// Ensure title doesn't exceed pane width
		maxTitleLen := len(firstLine) - 4 // Leave some space for borders
		if len(title) > maxTitleLen && maxTitleLen > 5 {
			// Truncate title if it's too long
			title = title[:maxTitleLen-3] + "..."
		}

		titlePos := (len(firstLine) - len(title)) / 2
		if titlePos < 0 {
			titlePos = 0
		}

		// Insert the title
		if titlePos+len(title) <= len(firstLine) {
			lines[0] = firstLine[:titlePos] + title + firstLine[titlePos+len(title):]
		}

		// Rejoin the lines
		return strings.Join(lines, "\n")
	}

	return renderedContent
}

// SetActive sets whether this pane is active
func (p *OutputPane) SetActive(active bool) {
	p.isActive = active
}

// SetSize sets the size of the pane
func (p *OutputPane) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.viewport.Width = width - 8
	p.viewport.Height = height - 8
}

// GetType returns the type of pane
func (p *OutputPane) GetType() PaneType {
	return PaneOutput
}

// SetContent sets the content of the output pane
func (p *OutputPane) SetContent(content string) {
	p.content = content
	p.viewport.SetContent(content)
}

// SetTitle sets the title of the output pane
func (p *OutputPane) SetTitle(title string) {
	p.title = title
}

// ExplorerPane is a wrapper around the Explorer
type ExplorerPane struct {
	width    int
	height   int
	explorer *Explorer
	isActive bool
}

// NewExplorerPane creates a new explorer pane
func NewExplorerPane(width, height int) *ExplorerPane {
	return &ExplorerPane{
		width:    width,
		height:   height,
		explorer: NewExplorer(width, height),
	}
}

// Update updates the explorer pane
func (p *ExplorerPane) Update(msg tea.Msg) (Pane, tea.Cmd) {
	var cmd tea.Cmd
	p.explorer, cmd = p.explorer.Update(msg)
	return p, cmd
}

// View renders the explorer pane
func (p *ExplorerPane) View() string {
	// Use the appropriate style based on whether this pane is active
	var style lipgloss.Style
	if p.isActive {
		style = activePaneStyle.Copy()
	} else {
		style = inactivePaneStyle.Copy()
	}

	style = style.Width(p.width - 4).
		Height(p.height - 2)

	// Add title to the border
	style = style.Border(lipgloss.RoundedBorder())

	// Add title to the top border with pane index
	title := "3: Explorer"

	renderedContent := style.Render(p.explorer.View())

	// Add title to the top border
	lines := strings.Split(renderedContent, "\n")
	if len(lines) > 0 {
		// Find the position to insert the title
		firstLine := lines[0]

		// Ensure title doesn't exceed pane width
		maxTitleLen := len(firstLine) - 4 // Leave some space for borders
		if len(title) > maxTitleLen && maxTitleLen > 5 {
			// Truncate title if it's too long
			title = title[:maxTitleLen-3] + "..."
		}

		titlePos := (len(firstLine) - len(title)) / 2
		if titlePos < 0 {
			titlePos = 0
		}

		// Insert the title
		if titlePos+len(title) <= len(firstLine) {
			lines[0] = firstLine[:titlePos] + title + firstLine[titlePos+len(title):]
		}

		// Rejoin the lines
		return strings.Join(lines, "\n")
	}

	return renderedContent
}

// SetActive sets whether this pane is active
func (p *ExplorerPane) SetActive(active bool) {
	p.isActive = active
}

// SetSize sets the size of the pane
func (p *ExplorerPane) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.explorer.width = width - 8
	p.explorer.height = height - 8
	p.explorer.list.SetSize(width/2 - 10, height - 12)
}

// GetType returns the type of pane
func (p *ExplorerPane) GetType() PaneType {
	return PaneExplorer
}

// GetSelectedFiles returns the selected files
func (p *ExplorerPane) GetSelectedFiles() []string {
	return p.explorer.GetSelectedFiles()
}

// GetExplorer returns the underlying explorer
func (p *ExplorerPane) GetExplorer() *Explorer {
	return p.explorer
}
