# Future Render — Build Targets
#
# Usage:
#   make          — run the default CI pipeline (fmt, vet, lint, test, build)
#   make ci       — same as default, explicit name for CI systems
#   make test     — run tests
#   make lint     — run golangci-lint
#   make fix      — auto-fix lint and formatting issues
#   make bench    — run benchmarks
#   make clean    — remove build artifacts
#
# Prerequisites:
#   go 1.24+
#   golangci-lint (https://golangci-lint.run/welcome/install/)

.PHONY: all ci fmt vet lint test test-race bench build clean fix check-tools

# Default target runs the full CI pipeline
all: ci

# --- CI Pipeline (order matters) ---

ci: fmt vet lint test build
	@echo "CI pipeline passed."

# --- Individual Targets ---

# Check formatting (fails if files need formatting)
fmt:
	@echo "==> Checking formatting..."
	@test -z "$$(gofmt -l .)" || { echo "Files need formatting:"; gofmt -l .; exit 1; }

# Go vet
vet:
	@echo "==> Running go vet..."
	go vet ./...

# Lint with golangci-lint
lint: check-lint
	@echo "==> Running golangci-lint..."
	golangci-lint run ./...

# Run all tests
test:
	@echo "==> Running tests..."
	go test ./...

# Run tests with race detector
test-race:
	@echo "==> Running tests with race detector..."
	go test -race ./...

# Run benchmarks
bench:
	@echo "==> Running benchmarks..."
	go test -bench=. -benchmem ./math/ ./internal/batch/

# Build all packages
build:
	@echo "==> Building..."
	go build ./...

# Auto-fix formatting and lint issues
fix: check-lint
	@echo "==> Fixing formatting..."
	gofmt -w .
	@echo "==> Fixing lint issues..."
	golangci-lint run --fix ./...

# Remove build artifacts
clean:
	@echo "==> Cleaning..."
	go clean ./...

# --- Tool Checks ---

check-lint:
	@which golangci-lint > /dev/null 2>&1 || { \
		echo "golangci-lint not found. Install: https://golangci-lint.run/welcome/install/"; \
		exit 1; \
	}
