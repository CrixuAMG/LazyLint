# LazyLint

A beautiful terminal user interface (TUI) for running various linters and code quality tools, inspired by Lazygit.

![LazyLint Screenshot](https://example.com/screenshot.png)

## Features

- Run multiple linters from a user-friendly terminal interface
- View formatted results in a scrollable viewport with syntax highlighting
- Configure tool paths and arguments per project
- File explorer with preview for selecting specific files to lint
- Automatic detection of tools in your project
- Beautiful UI with borders, colors, and intuitive layout
- Support for multiple languages and linters:
  - PHP: PHPStan, PHPCS
  - Go: golangci-lint
  - JavaScript/TypeScript: ESLint

## Installation

### Prerequisites

- Git (for repository detection)
- Linters for your languages of choice

### Using Homebrew (macOS and Linux)

```bash
# Add the tap
brew tap crixuamg/lazylint

# Install LazyLint
brew install lazylint
```

### Local Installation

```bash
# Clone the repository
git clone https://github.com/crixuamg/lazylint.git
cd lazylint

# Build the application
make build

# Optional: Install the binary to your GOPATH
make install

# Or move it manually to a directory in your PATH
sudo mv lazylint /usr/local/bin/
```

### Global Installation with Go

You can install the tool globally using Go:

```bash
go install github.com/crixuamg/lazylint/cmd/lazylint@latest
```

### Using a Pre-built Binary

Download the latest release for your platform from the [Releases page](https://github.com/crixuamg/lazylint/releases).

```bash
# Make it executable
chmod +x lazylint

# Move to a directory in your PATH
sudo mv lazylint /usr/local/bin/
```

## Configuration

LazyLint automatically looks for a configuration file named `lazylint.yaml` or `lazylint.yml` in the current directory or any parent directory.

### Default Configuration

By default, LazyLint will look for linters in standard locations:

- PHP linters in `vendor/bin/` directory
- Go linters in your PATH or local `bin/` directory
- JavaScript/TypeScript linters in `node_modules/.bin/` directory

### Creating a Configuration File

You can create a default configuration file using:

```bash
lazylint --create-config
```

### Theme Customization

LazyLint supports custom themes that can be defined in your configuration file. Themes are stored in `$HOME/.config/lazylint/config.yaml` and can be switched using the `t` key while the application is running.

Each theme defines colors for various UI elements using single hex color values. The application background color will also update based on the selected theme.

Here's an example of a custom theme configuration:

```yaml
ui:
  theme: "catppuccin"  # Currently selected theme
  themes:
    catppuccin:
      name: "catppuccin"
      colors:
        subtle: "#6E6C7E"      # Used for subtle elements
        highlight: "#F5C2E7"   # Used for highlighting active elements
        special: "#ABE9B3"     # Used for special elements
        error: "#F28FAD"       # Used for error messages
        warning: "#FAE3B0"     # Used for warnings
        border: "#96CDFB"      # Used for borders
        text: "#D9E0EE"        # Used for normal text
        dim_text: "#988BA2"    # Used for dimmed text
        background: "#1E1E2E"  # Used for backgrounds
```

LazyLint comes with several built-in themes:
- `default`: A dark theme with purple accents
- `light`: A light theme with purple accents
- `high-contrast`: A high contrast theme for accessibility
- `catppuccin`: A soothing pastel theme

See the `examples/config.yaml` file for more theme examples.

### Example Configuration

```yaml
# UI Configuration
ui:
  # Current theme
  theme: "default"

  # Custom themes
  themes:
    # Custom blue theme
    custom-blue:
      name: "custom-blue"
      colors:
        subtle: "#BBDEFB,#263238"      # Light mode,Dark mode
        highlight: "#2196F3,#64B5F6"
        special: "#00BCD4,#4DD0E1"
        error: "#F44336,#EF5350"
        warning: "#FF9800,#FFB74D"
        border: "#1976D2,#42A5F5"
        text: "#212121,#ECEFF1"
        dim_text: "#757575,#B0BEC5"
        background: "#E3F2FD"

# Linter Configuration
linters:
  phpstan:
    # Path to the PHPStan executable
    path: /path/to/your/project/vendor/bin/phpstan
    # Arguments to pass to PHPStan
    args:
      - analyse
      - --level=5
      - --no-progress
    # Enable or disable PHPStan
    enabled: true

  phpcs:
    # Path to the PHPCS executable
    path: /path/to/your/project/vendor/bin/phpcs
    # Arguments to pass to PHPCS
    args:
      - --standard=PSR12
      - --colors
    # Enable or disable PHPCS
    enabled: true

  golangci-lint:
    # Path to the golangci-lint executable
    path: golangci-lint
    # Arguments to pass to golangci-lint
    args:
      - run
      - --out-format=colored-line-number
    # Enable or disable golangci-lint
    enabled: true

  eslint:
    # Path to the ESLint executable
    path: ./node_modules/.bin/eslint
    # Arguments to pass to ESLint
    args:
      - --format=stylish
    # Enable or disable ESLint
    enabled: true
```

## Usage

```bash
# Run the application in the current directory
lazylint

# Specify a target file or directory
lazylint --target=/path/to/your/code

# Specify a configuration file
lazylint --config=/path/to/config.yaml

# Create a default configuration file
lazylint --create-config

# Show version information
lazylint --version
```

## Keyboard Shortcuts

| Key       | Action                |
|-----------|----------------------|
| `↑` or `k` | Move up               |
| `↓` or `j` | Move down             |
| `Enter`   | Select                |
| `Esc`     | Go back               |
| `?`       | Toggle help           |
| `q`       | Quit                  |
| `Ctrl+C`  | Quit                  |
| `t`       | Cycle through themes  |
| `1-4`     | Switch between panes  |
| `h/l`     | Navigate between panes|

In the file explorer:
| Key       | Action                |
|-----------|----------------------|
| `Space`   | Select file           |
| `Enter`   | Open file/directory   |
| `Tab`     | Toggle preview        |
| `r`       | Run linters on selected files |

## Development

### Running Tests

```bash
go test ./...
```

### Building for Different Platforms

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o lazylint-linux-amd64 cmd/lazylint/main.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o lazylint-macos-amd64 cmd/lazylint/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o lazylint-windows-amd64.exe cmd/lazylint/main.go
```

## Implementing a New Linter

LazyLint is designed to be easily extensible with new linters. Here's how to add a new linter:

1. Create a new file in the `pkg/linters` directory, e.g., `mylinter.go`
2. Implement the `Linter` interface:

```go
package linters

import (
    "context"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
)

// MyLinter implements the Linter interface
type MyLinter struct {
    path    string
    args    []string
    enabled bool
}

// NewMyLinter creates a new MyLinter
func NewMyLinter() *MyLinter {
    // Find the linter executable
    linterPath := "mylinter" // Default path

    // Check for project-specific path
    gitRoot, _ := findGitRoot()
    if gitRoot != "" {
        // Check custom locations based on your linter
        customPath := filepath.Join(gitRoot, "tools", "mylinter")
        if _, err := os.Stat(customPath); err == nil {
            linterPath = customPath
        }
    }

    return &MyLinter{
        path:    linterPath,
        args:    []string{"--default-args"},
        enabled: true,
    }
}

// Name returns the name of the linter
func (l *MyLinter) Name() string {
    return "mylinter"
}

// Description returns a short description of the linter
func (l *MyLinter) Description() string {
    return "My custom linter for X language"
}

// Run executes the linter on the given target
func (l *MyLinter) Run(ctx context.Context, target string) (*Result, error) {
    if !l.enabled {
        return &Result{
            Name:      l.Name(),
            Success:   true,
            Output:    "MyLinter is disabled in configuration",
            Timestamp: time.Now(),
        }, nil
    }

    args := append([]string{}, l.args...)
    if target != "" {
        args = append(args, target)
    }

    start := time.Now()
    cmd := exec.CommandContext(ctx, l.path, args...)

    var stdout, stderr strings.Builder
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()
    duration := time.Since(start)

    result := &Result{
        Name:      l.Name(),
        Output:    stdout.String(),
        Error:     stderr.String(),
        Duration:  duration,
        Timestamp: time.Now(),
    }

    if err != nil {
        // Check if it's a timeout
        if ctx.Err() == context.DeadlineExceeded {
            return result, fmt.Errorf("command timed out after %s", duration)
        }

        // Check if it's an exit code error (which is expected for linters when they find issues)
        if _, ok := err.(*exec.ExitError); ok {
            result.Success = false
            return result, nil
        }

        return result, fmt.Errorf("command failed: %w", err)
    }

    result.Success = true
    return result, nil
}

// IsAvailable checks if the linter is available
func (l *MyLinter) IsAvailable() bool {
    _, err := exec.LookPath(l.path)
    return err == nil
}

// FileExtensions returns the file extensions this linter can process
func (l *MyLinter) FileExtensions() []string {
    return []string{".mylang"}
}

// Configure configures the linter with the given options
func (l *MyLinter) Configure(options map[string]interface{}) error {
    if path, ok := options["path"].(string); ok {
        l.path = path
    }

    if args, ok := options["args"].([]string); ok {
        l.args = args
    } else if argsInterface, ok := options["args"].([]interface{}); ok {
        args := make([]string, len(argsInterface))
        for i, arg := range argsInterface {
            if str, ok := arg.(string); ok {
                args[i] = str
            }
        }
        l.args = args
    }

    if enabled, ok := options["enabled"].(bool); ok {
        l.enabled = enabled
    }

    return nil
}
```

3. Register your linter in `pkg/linters/registry.go`:

```go
func DefaultRegistry() *Registry {
    registry := NewRegistry()

    // Register existing linters
    registry.Register(NewPHPStan())
    registry.Register(NewPHPCS())
    registry.Register(NewGolangCI())
    registry.Register(NewESLint())

    // Register your new linter
    registry.Register(NewMyLinter())

    return registry
}
```

4. Add default configuration in `pkg/config/config.go`:

```go
func DefaultConfig() *Config {
    return &Config{
        Linters: map[string]map[string]interface{}{
            // Existing linters...

            "mylinter": {
                "path":    "mylinter",
                "args":    []string{"--default-args"},
                "enabled": true,
            },
        },
    }
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT
