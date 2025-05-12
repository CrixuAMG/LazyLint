package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/crixuamg/pkg/linters"
)

// FileItem represents a file in the explorer
type FileItem struct {
	path     string
	name     string
	isDir    bool
	selected bool
}

// FilterValue implements list.Item interface
func (i FileItem) FilterValue() string { return i.name }

// Title returns the title of the item
func (i FileItem) Title() string {
	if i.selected {
		return selectedItemStyle.Render("✓ " + i.name)
	}
	if i.isDir {
		return dirStyle.Render("📁 " + i.name)
	}
	return fileStyle.Render("📄 " + i.name)
}

// Description returns the description of the item
func (i FileItem) Description() string {
	return i.path
}

// Explorer represents the file explorer component
type Explorer struct {
	list          list.Model
	selectedFiles map[string]bool
	currentDir    string
	rootDir       string
	preview       string
	width         int
	height        int
	showPreview   bool
	registry      *linters.Registry
	fileFilter    string
}

// NewExplorer creates a new explorer
func NewExplorer(width, height int) *Explorer {
	// Get git root directory
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	rootDir := "."
	if err == nil {
		rootDir = strings.TrimSpace(string(out))
	}

	// Create list
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = selectedItemStyle
	delegate.Styles.SelectedDesc = selectedItemStyle
	delegate.SetHeight(1)

	l := list.New([]list.Item{}, delegate, width/2, height-4)
	l.Title = "Files"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false) // Hide the help text
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.FilterPrompt = infoStyle
	l.Styles.FilterCursor = infoStyle

	e := &Explorer{
		list:          l,
		selectedFiles: make(map[string]bool),
		currentDir:    rootDir,
		rootDir:       rootDir,
		width:         width,
		height:        height,
		showPreview:   true,
		fileFilter:    "",
	}

	e.loadFiles()
	return e
}

// loadFiles loads files from the current directory
func (e *Explorer) loadFiles() {
	var items []list.Item

	// Get git tracked files
	cmd := exec.Command("git", "ls-files", "--full-name", e.currentDir)
	out, err := cmd.Output()
	if err == nil {
		files := strings.Split(strings.TrimSpace(string(out)), "\n")

		// Add parent directory if not at root
		if e.currentDir != e.rootDir {
			items = append(items, FileItem{
				path:  filepath.Dir(e.currentDir),
				name:  "..",
				isDir: true,
			})
		}

		// Process files
		dirs := make(map[string]bool)
		for _, file := range files {
			if file == "" {
				continue
			}

			// Get full path
			fullPath := filepath.Join(e.rootDir, file)

			// Check if it's in current directory or subdirectory
			rel, err := filepath.Rel(e.currentDir, fullPath)
			if err != nil || strings.HasPrefix(rel, "..") {
				continue
			}

			// Filter by extension if registry is set
			if e.registry != nil && !strings.HasPrefix(filepath.Base(fullPath), ".") {
				ext := filepath.Ext(fullPath)
				if ext != "" && len(e.registry.GetForExtension(ext)) == 0 {
					// Skip files that don't match any linter
					continue
				}
			}

			// If it's a direct child, add it
			parts := strings.Split(rel, string(os.PathSeparator))
			if len(parts) == 1 {
				info, err := os.Stat(fullPath)
				if err != nil {
					continue
				}

				items = append(items, FileItem{
					path:     fullPath,
					name:     filepath.Base(fullPath),
					isDir:    info.IsDir(),
					selected: e.selectedFiles[fullPath],
				})
			} else if len(parts) > 1 {
				// If it's in a subdirectory, add the directory if not already added
				dirPath := filepath.Join(e.currentDir, parts[0])
				if !dirs[dirPath] {
					dirs[dirPath] = true
					items = append(items, FileItem{
						path:  dirPath,
						name:  parts[0],
						isDir: true,
					})
				}
			}
		}
	}

	e.list.SetItems(items)
}

// Update updates the explorer
func (e *Explorer) Update(msg tea.Msg) (*Explorer, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		e.width = msg.Width
		e.height = msg.Height
		e.list.SetWidth(msg.Width / 2)
		e.list.SetHeight(msg.Height - 4)

	case tea.KeyMsg:
		// Handle custom keybindings
		keyMsg := msg
		switch keyMsg.String() {
		case "enter":
			if i, ok := e.list.SelectedItem().(FileItem); ok {
				if i.isDir {
					e.currentDir = i.path
					e.loadFiles()
					e.preview = ""
					return e, nil
				} else {
					// Load file preview
					content, err := os.ReadFile(i.path)
					if err == nil {
						e.preview = string(content)
					} else {
						e.preview = fmt.Sprintf("Error loading file: %s", err)
					}
				}
			}
		case " ":
			if i, ok := e.list.SelectedItem().(FileItem); ok && !i.isDir {
				// Toggle selection
				if e.selectedFiles[i.path] {
					delete(e.selectedFiles, i.path)
				} else {
					e.selectedFiles[i.path] = true
				}
				e.loadFiles()
			}
		case "tab":
			e.showPreview = !e.showPreview
		}
	}

	// Update list
	e.list, cmd = e.list.Update(msg)
	return e, cmd
}

// View renders the explorer
func (e *Explorer) View() string {
	// Create file list view
	fileList := e.list.View()

	// Create preview view
	var preview string
	if e.showPreview && e.preview != "" {
		preview = previewStyle.Width(e.width/2 - 4).Render(e.preview)
	} else if e.showPreview {
		preview = previewStyle.Width(e.width/2 - 4).Render("Select a file to preview")
	} else {
		// Show selected files when preview is hidden
		var selectedList strings.Builder
		selectedList.WriteString(titleStyle.Render("Selected Files"))
		selectedList.WriteString("\n\n")

		if len(e.selectedFiles) == 0 {
			selectedList.WriteString("No files selected")
		} else {
			for path := range e.selectedFiles {
				selectedList.WriteString(fmt.Sprintf("- %s\n", filepath.Base(path)))
			}
		}
		preview = previewStyle.Width(e.width/2 - 4).Render(selectedList.String())
	}

	// Combine views
	return lipgloss.JoinHorizontal(lipgloss.Top, fileList, preview)
}

// GetSelectedFiles returns the selected files
func (e *Explorer) GetSelectedFiles() []string {
	var files []string
	for path := range e.selectedFiles {
		files = append(files, path)
	}
	return files
}

// Additional styles for the explorer
var (
	dirStyle = lipgloss.NewStyle().
		Foreground(primary)

	fileStyle = lipgloss.NewStyle().
		Foreground(text)

	previewStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderClr).
		Padding(1, 2).
		Margin(0, 0, 0, 2)
)
