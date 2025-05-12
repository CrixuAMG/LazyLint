# Publishing LazyLint on Homebrew

This document explains how to publish LazyLint on Homebrew using GoReleaser.

## Prerequisites

1. A GitHub repository for LazyLint
2. A GitHub repository for the Homebrew tap (e.g., `homebrew-lazylint`)
3. [GoReleaser](https://goreleaser.com/) installed locally
4. A GitHub token with repo scope

## Setup

### 1. Create a Homebrew Tap Repository

Create a new GitHub repository named `homebrew-lazylint` to host your Homebrew formula.

```bash
# Clone the main repository
git clone https://github.com/crixuamg/lazylint.git
cd lazylint

# Create the tap repository
mkdir -p ../homebrew-lazylint
cd ../homebrew-lazylint
git init
echo "# LazyLint Homebrew Tap" > README.md
git add README.md
git commit -m "Initial commit"
git branch -M main
git remote add origin https://github.com/crixuamg/homebrew-lazylint.git
git push -u origin main
```

### 2. Configure GoReleaser

The `.goreleaser.yml` file in the main repository is already configured for Homebrew. It includes:

- Build configuration for multiple platforms
- Archive settings
- Homebrew tap configuration

### 3. Create a GitHub Release

To publish a new version:

1. Tag a new version:

```bash
cd lazylint
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
```

2. Run GoReleaser:

```bash
# Set your GitHub token
export GITHUB_TOKEN=your_github_token

# Run GoReleaser
goreleaser release --rm-dist
```

This will:
- Build binaries for all platforms
- Create archives
- Generate a Homebrew formula
- Push the formula to your tap repository
- Create a GitHub release with the binaries

### 4. Install from Homebrew

Users can now install LazyLint using:

```bash
# Add the tap
brew tap crixuamg/lazylint

# Install LazyLint
brew install lazylint
```

## Updating

To release a new version:

1. Tag a new version:

```bash
git tag -a v0.1.1 -m "Bug fixes and improvements"
git push origin v0.1.1
```

2. Run GoReleaser as before:

```bash
export GITHUB_TOKEN=your_github_token
goreleaser release --rm-dist
```

The Homebrew formula will be automatically updated in your tap repository.

## Troubleshooting

- If you get permission errors, make sure your GitHub token has the correct permissions.
- If the formula doesn't update, check the GoReleaser logs for errors.
- If users have issues installing, make sure the tap repository is public.

## Resources

- [GoReleaser Documentation](https://goreleaser.com/intro/)
- [Homebrew Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)
- [GitHub Actions for GoReleaser](https://goreleaser.com/ci/actions/)
