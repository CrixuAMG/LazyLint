package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/crixuamg/pkg/config"
)

// ThemeColors defines the colors used in a theme
type ThemeColors struct {
	Subtle     lipgloss.Color
	Highlight  lipgloss.Color
	Special    lipgloss.Color
	Error      lipgloss.Color
	Warning    lipgloss.Color
	Border     lipgloss.Color
	Text       lipgloss.Color
	DimText    lipgloss.Color
	Background lipgloss.Color
}

// Theme defines a complete theme with colors and styles
type Theme struct {
	Name   string
	Colors ThemeColors
}

// Theme variables
var (
	// Available themes
	Themes = make(map[string]Theme)

	// Current theme
	CurrentTheme Theme
)

// ApplyTheme applies the given theme to the UI
func ApplyTheme(themeName string) {
	// Get theme from available themes
	theme, ok := Themes[themeName]
	if !ok {
		// Use first available theme if the requested theme doesn't exist
		for _, t := range Themes {
			theme = t
			break
		}

		// If no themes are available, create a default theme
		if !ok && len(Themes) == 0 {
			theme = Theme{
				Name: "tokyo-night",
				Colors: ThemeColors{
					Subtle:     lipgloss.Color("#565F89"),
					Highlight:  lipgloss.Color("#7AA2F7"),
					Special:    lipgloss.Color("#9ECE6A"),
					Error:      lipgloss.Color("#F7768E"),
					Warning:    lipgloss.Color("#E0AF68"),
					Border:     lipgloss.Color("#565F89"),
					Text:       lipgloss.Color("#C0CAF5"),
					DimText:    lipgloss.Color("#9AA5CE"),
					Background: lipgloss.Color("#1A1B26"),
				},
			}
		}
	}

	// Set current theme
	CurrentTheme = theme

	// Update global styles
	// Map theme colors to our new color scheme
	muted = CurrentTheme.Colors.Subtle
	primary = CurrentTheme.Colors.Highlight
	secondary = CurrentTheme.Colors.Highlight
	accent = CurrentTheme.Colors.Special
	success = CurrentTheme.Colors.Special
	errorClr = CurrentTheme.Colors.Error
	warningClr = CurrentTheme.Colors.Warning
	borderClr = CurrentTheme.Colors.Border
	background = CurrentTheme.Colors.Background
	text = CurrentTheme.Colors.Text
	subtext = CurrentTheme.Colors.DimText
	surface1 = lipgloss.Color("#24283B")
	surface2 = lipgloss.Color("#414868")
	info = CurrentTheme.Colors.Highlight

	// Update all styles based on the current theme

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
		Foreground(text).
		MarginLeft(2)

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
		Margin(1, 0)

	// Update explorer styles
	dirStyle = lipgloss.NewStyle().
		Foreground(primary)

	fileStyle = lipgloss.NewStyle().
		Foreground(text)

	previewStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderClr).
		Padding(1, 2).
		Margin(0, 0, 0, 2)
}

// parseColor parses a color string
func parseColor(colorStr string) lipgloss.Color {
	// Default color
	defaultColor := "#FFFFFF"

	// If we have a color, use it
	if colorStr != "" {
		return lipgloss.Color(colorStr)
	}

	return lipgloss.Color(defaultColor)
}

// convertConfigTheme converts a config theme to a UI theme
func convertConfigTheme(configTheme config.ThemeConfig) Theme {
	return Theme{
		Name: configTheme.Name,
		Colors: ThemeColors{
			Subtle:     parseColor(configTheme.Colors.Subtle),
			Highlight:  parseColor(configTheme.Colors.Highlight),
			Special:    parseColor(configTheme.Colors.Special),
			Error:      parseColor(configTheme.Colors.Error),
			Warning:    parseColor(configTheme.Colors.Warning),
			Border:     parseColor(configTheme.Colors.Border),
			Text:       parseColor(configTheme.Colors.Text),
			DimText:    parseColor(configTheme.Colors.DimText),
			Background: parseColor(configTheme.Colors.Background),
		},
	}
}

// InitTheme initializes the theme from the configuration
func InitTheme(cfg *config.Config) {
	// Load themes from config
	for name, configTheme := range cfg.UI.Themes {
		Themes[name] = convertConfigTheme(configTheme)
	}

	// Apply theme from config
	themeName := cfg.UI.Theme

	// Validate theme name
	_, ok := Themes[themeName]
	if !ok {
		// Default to "default" theme if the configured theme doesn't exist
		themeName = "default"
		cfg.UI.Theme = themeName
	}

	ApplyTheme(themeName)
}

// convertThemeToConfig converts a UI theme to a config theme
func convertThemeToConfig(theme Theme) config.ThemeConfig {
	// Convert colors to string format
	subtleStr := string(theme.Colors.Subtle)
	highlightStr := string(theme.Colors.Highlight)
	specialStr := string(theme.Colors.Special)
	errorStr := string(theme.Colors.Error)
	warningStr := string(theme.Colors.Warning)
	borderStr := string(theme.Colors.Border)
	textStr := string(theme.Colors.Text)
	dimTextStr := string(theme.Colors.DimText)
	backgroundStr := string(theme.Colors.Background)

	return config.ThemeConfig{
		Name: theme.Name,
		Colors: config.ThemeColors{
			Subtle:     subtleStr,
			Highlight:  highlightStr,
			Special:    specialStr,
			Error:      errorStr,
			Warning:    warningStr,
			Border:     borderStr,
			Text:       textStr,
			DimText:    dimTextStr,
			Background: backgroundStr,
		},
	}
}

// SaveTheme saves the current theme to the configuration
func SaveTheme(cfg *config.Config, themeName string) error {
	// Update config
	cfg.UI.Theme = themeName

	// Make sure the theme exists in the config
	theme, ok := Themes[themeName]
	if ok {
		// Convert theme to config format
		configTheme := convertThemeToConfig(theme)

		// Add or update theme in config
		if cfg.UI.Themes == nil {
			cfg.UI.Themes = make(map[string]config.ThemeConfig)
		}
		cfg.UI.Themes[themeName] = configTheme
	}

	// Save to user config
	return config.SaveUserConfig(cfg)
}
