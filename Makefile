.PHONY: build clean test run install

# Build variables
BINARY_NAME=lazylint
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)"

# Default target
all: build

# Build the application
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) cmd/lazylint/*.go

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

# Run tests
test:
	go test -v ./...

# Run the application
run:
	go run cmd/lazylint/*.go

# Install the application
install: build
	mv $(BINARY_NAME) $(GOPATH)/bin/

# Release with goreleaser
release:
	goreleaser release --clean

# Create a snapshot release for testing
snapshot:
	goreleaser release --snapshot --clean
