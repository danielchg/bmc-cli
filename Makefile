# BMC CLI Makefile

# Variables
BINARY_NAME=bmc-cli
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe
VERSION=1.0.0
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_UNIX) .

# Build for Windows
.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_WINDOWS) .

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-windows

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Linting code..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Run 'make dev-setup' or use 'make lint-basic'" && exit 1)
	golangci-lint run --disable-all --enable=errcheck,gosimple,govet,ineffassign,staticcheck,typecheck,unused,gofmt,goimports,misspell --timeout=5m

# Basic linting without golangci-lint
.PHONY: lint-basic
lint-basic: fmt vet
	@echo "Running basic linting..."

# Vet code
.PHONY: vet
vet:
	@echo "Vetting code..."
	go vet ./...

# Run all quality checks
.PHONY: check
check: fmt vet lint-basic test

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)
	rm -f coverage.out
	rm -f coverage.html

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Install the binary to GOPATH/bin
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(LDFLAGS) .

# Generate sample config
.PHONY: config
config: build
	@echo "Generating sample configuration..."
	./$(BINARY_NAME) config generate

# Run the binary with help
.PHONY: run-help
run-help: build
	./$(BINARY_NAME) --help

# Development setup
.PHONY: dev-setup
dev-setup: deps
	@echo "Setting up development environment..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  build-linux   - Build for Linux"
	@echo "  build-windows - Build for Windows"
	@echo "  build-all     - Build for all platforms"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code (requires golangci-lint)"
	@echo "  lint-basic    - Basic linting without golangci-lint"
	@echo "  vet           - Vet code"
	@echo "  check         - Run all quality checks"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  install       - Install binary to GOPATH/bin"
	@echo "  config        - Generate sample configuration"
	@echo "  run-help      - Build and show help"
	@echo "  dev-setup     - Setup development environment"
	@echo "  help          - Show this help message" 