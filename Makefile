.PHONY: build install clean test help dev release

# Binary name
BINARY_NAME=repo-manager
BINARY_PATH=bin/$(BINARY_NAME)

# Build variables
BUILD_DIR=bin
CMD_DIR=cmd/repo-manager
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Installation directory
INSTALL_DIR=/usr/local/bin

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Default target
all: build

## help: Display this help message
help:
	@echo "Repository Manager - Build System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' Makefile | sed 's/## /  /'

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) ./$(CMD_DIR)
	@echo "✓ Built: $(BINARY_PATH)"
	@ls -lh $(BINARY_PATH)

## install: Install the binary to $(INSTALL_DIR)
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@sudo cp $(BINARY_PATH) $(INSTALL_DIR)/
	@echo "✓ Installed: $(INSTALL_DIR)/$(BINARY_NAME)"
	@which $(BINARY_NAME)

## dev: Build with race detector (for development)
dev:
	@echo "Building $(BINARY_NAME) with race detector..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -race $(LDFLAGS) -o $(BINARY_PATH) ./$(CMD_DIR)
	@echo "✓ Built (dev): $(BINARY_PATH)"

## release: Build optimized release binary
release:
	@echo "Building release version..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -a -installsuffix cgo -o $(BINARY_PATH) ./$(CMD_DIR)
	@echo "✓ Built (release): $(BINARY_PATH)"
	@ls -lh $(BINARY_PATH)

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "✓ Cleaned"

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "✓ Formatted"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	@echo "✓ Dependencies downloaded"

## tidy: Tidy and verify dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy
	$(GOMOD) verify
	@echo "✓ Dependencies tidied"

## version: Show version information
version:
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Built:   $(BUILD_TIME)"

## run: Build and run the binary
run: build
	@$(BINARY_PATH) --help

## shell-bash: Install bash shell integration
shell-bash:
	@echo "Installing bash integration..."
	@echo "Add this to your ~/.bashrc:"
	@echo ""
	@echo "source $(PWD)/scripts/repo-manager.bash"
	@echo ""

## shell-zsh: Install zsh shell integration
shell-zsh:
	@echo "Installing zsh integration..."
	@echo "Add this to your ~/.zshrc:"
	@echo ""
	@echo "source $(PWD)/scripts/repo-manager.zsh"
	@echo ""

## uninstall: Remove installed binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ Uninstalled"
