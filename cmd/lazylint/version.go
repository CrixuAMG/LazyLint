package main

// Version information set by build flags
var (
	// Version is the current version of LazyLint
	Version = "dev"
	// Commit is the git commit hash
	Commit = "none"
	// Date is the build date
	Date = "unknown"
)

// GetVersion returns the full version string
func GetVersion() string {
	return Version + " (" + Commit + ") built on " + Date
}
