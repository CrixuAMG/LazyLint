package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	// Base colors
	primary    = lipgloss.Color("#7AA2F7") // Vibrant blue
	secondary  = lipgloss.Color("#BB9AF7") // Purple
	accent     = lipgloss.Color("#2AC3DE") // Cyan
	success    = lipgloss.Color("#9ECE6A") // Green
	warningClr = lipgloss.Color("#E0AF68") // Orange
	errorClr   = lipgloss.Color("#F7768E") // Red
	info       = lipgloss.Color("#7DCFFF") // Light blue

	// Grayscale
	background = lipgloss.Color("#1A1B26") // Dark blue background
	surface1   = lipgloss.Color("#24283B") // Slightly lighter than background
	surface2   = lipgloss.Color("#414868") // Even lighter for active elements
	borderClr  = lipgloss.Color("#565F89") // Border color
	text       = lipgloss.Color("#C0CAF5") // Main text color
	subtext    = lipgloss.Color("#9AA5CE") // Secondary text color
	muted      = lipgloss.Color("#565F89") // Muted text
)

// UI Components
var (
	// App title
	appTitleStyle = lipgloss.NewStyle().
		Foreground(primary).
		Background(surface1).
		Bold(true).
		Padding(0, 1).
		Margin(0, 0, 1, 0)

	// Tab styles
	activeTabStyle = lipgloss.NewStyle().
		Foreground(background).
		Background(primary).
		Bold(true).
		Padding(0, 3).
		Margin(0, 1, 0, 0)

	inactiveTabStyle = lipgloss.NewStyle().
		Foreground(text).
		Background(surface1).
		Padding(0, 3).
		Margin(0, 1, 0, 0)

	tabContentStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primary).
		Padding(1, 2).
		BorderTop(false).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true)

	// Panel styles
	panelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderClr).
		Padding(1, 2).
		Margin(1, 1)

	activePanelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primary).
		Padding(1, 2).
		Margin(1, 1)

	// Title styles
	titleStyle = lipgloss.NewStyle().
		Foreground(primary).
		Bold(true).
		MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
		Foreground(secondary).
		Bold(true)

	// Content styles
	selectedItemStyle = lipgloss.NewStyle().
		Foreground(accent).
		Bold(true)

	itemStyle = lipgloss.NewStyle().
		Foreground(text)

	// Status styles
	statusBarStyle = lipgloss.NewStyle().
		Foreground(text).
		Background(surface2).
		Padding(0, 1).
		Height(1).
		Width(100)

	// Help bar
	helpBarStyle = lipgloss.NewStyle().
		Foreground(text).
		Background(surface1).
		Padding(0, 1).
		Height(1).
		Width(100)

	helpStyle = lipgloss.NewStyle().
		Foreground(text).
		Background(surface1).
		Padding(0, 1).
		Height(1).
		Width(100)

	// Message styles
	successStyle = lipgloss.NewStyle().
		Foreground(success).
		Bold(true)

	errorStyle = lipgloss.NewStyle().
		Foreground(errorClr).
		Bold(true)

	warningStyle = lipgloss.NewStyle().
		Foreground(warningClr).
		Bold(true)

	infoStyle = lipgloss.NewStyle().
		Foreground(info)

	// Box styles for content
	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderClr).
		Padding(1, 2).
		Margin(1, 0)

	resultBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderClr).
		Padding(1, 2).
		Margin(1, 0)

	// Pane styles (for compatibility)
	activePaneStyle = activePanelStyle
	inactivePaneStyle = panelStyle

	// Badge styles
	badgeStyle = lipgloss.NewStyle().
		Foreground(background).
		Background(accent).
		Padding(0, 1).
		Margin(0, 1, 0, 0).
		Bold(true)

	// Button styles
	buttonStyle = lipgloss.NewStyle().
		Foreground(background).
		Background(primary).
		Padding(0, 3).
		Margin(0, 1).
		Bold(true)

	// Logo style
	logoStyle = lipgloss.NewStyle().
		Foreground(primary).
		Bold(true).
		Margin(1, 0).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primary).
		Padding(1, 2).
		Width(20).
		Align(lipgloss.Center)
)
