# Variables
BINARY_NAME=gingest
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD)
DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-X main.Version=$(VERSION) -X main.GitCommit=$(COMMIT) -X main.BuildDate=$(DATE)

# Default target
.PHONY: all
all: test build

# Build the binary
.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) cmd/gingest/main.go

# Build for all platforms
.PHONY: build-all
build-all:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-linux-amd64 cmd/gingest/main.go
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-linux-arm64 cmd/gingest/main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-darwin-amd64 cmd/gingest/main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-darwin-arm64 cmd/gingest/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-windows-amd64.exe cmd/gingest/main.go
	GOOS=windows GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-windows-arm64.exe cmd/gingest/main.go

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run integration tests
.PHONY: test-integration
test-integration:
	go test -tags=integration -v

# Run all tests
.PHONY: test-all
test-all: test test-integration

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests with race detector
.PHONY: test-race
test-race:
	go test -race -v ./...

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Run linters
.PHONY: lint
lint:
	go vet ./...
	gofmt -s -l .

# Install staticcheck and run it
.PHONY: staticcheck
staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out coverage.html

# Install the binary
.PHONY: install
install:
	go install -ldflags "$(LDFLAGS)" cmd/gingest/main.go

# Run the binary with help
.PHONY: run-help
run-help: build
	./$(BINARY_NAME) --help

# Run the binary with version
.PHONY: run-version
run-version: build
	./$(BINARY_NAME) --version

# Build Docker image
.PHONY: docker-build
docker-build:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg DATE=$(DATE) \
		-t gingest:$(VERSION) \
		-t gingest:latest \
		.

# Run Docker container
.PHONY: docker-run
docker-run:
	docker run --rm gingest:latest

# Development setup
.PHONY: dev-setup
dev-setup:
	go mod download
	go install honnef.co/go/tools/cmd/staticcheck@latest

# Create a release (for maintainers)
.PHONY: release
release:
	@echo "Creating release $(VERSION)"
	@if [ "$(VERSION)" = "dev" ]; then \
		echo "Error: VERSION must be set (e.g., make release VERSION=v1.0.0)"; \
		exit 1; \
	fi
	git tag $(VERSION)
	git push origin $(VERSION)

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  build-all      - Build for all platforms"
	@echo "  test           - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-all       - Run all tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  test-race      - Run tests with race detector"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linters"
	@echo "  staticcheck    - Run staticcheck"
	@echo "  clean          - Clean build artifacts"
	@echo "  install        - Install the binary"
	@echo "  run-help       - Build and show help"
	@echo "  run-version    - Build and show version"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  dev-setup      - Set up development environment"
	@echo "  release        - Create a release (VERSION=vX.Y.Z)"
	@echo "  help           - Show this help" 