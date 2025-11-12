#!/bin/bash
# Build script for Repository Manager

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="repo-manager"
BUILD_DIR="bin"
CMD_DIR="cmd/repo-manager"
INSTALL_DIR="/usr/local/bin"

# Version information
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

# Functions
print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${YELLOW}➜${NC} $1"
}

show_help() {
    cat << EOF
Repository Manager - Build Script

Usage: $0 [command]

Commands:
    build       Build the binary (default)
    install     Build and install to $INSTALL_DIR
    dev         Build with race detector
    release     Build optimized release binary
    clean       Remove build artifacts
    test        Run tests
    deps        Download dependencies
    tidy        Tidy dependencies
    run         Build and run with --help
    uninstall   Remove installed binary
    help        Show this help message

Examples:
    $0              # Build the binary
    $0 install      # Build and install
    $0 clean build  # Clean then build

Version: $VERSION
Commit:  $COMMIT
EOF
}

build() {
    print_info "Building $BINARY_NAME..."
    mkdir -p "$BUILD_DIR"

    go build \
        -ldflags "-X main.Version=$VERSION -X main.Commit=$COMMIT -X main.BuildTime=$BUILD_TIME" \
        -o "$BUILD_DIR/$BINARY_NAME" \
        ./"$CMD_DIR"

    print_success "Built: $BUILD_DIR/$BINARY_NAME"
    ls -lh "$BUILD_DIR/$BINARY_NAME"
}

install() {
    build
    print_info "Installing $BINARY_NAME to $INSTALL_DIR..."
    sudo cp "$BUILD_DIR/$BINARY_NAME" "$INSTALL_DIR/"
    print_success "Installed: $INSTALL_DIR/$BINARY_NAME"
    which "$BINARY_NAME"
}

build_dev() {
    print_info "Building $BINARY_NAME with race detector..."
    mkdir -p "$BUILD_DIR"

    go build \
        -race \
        -ldflags "-X main.Version=$VERSION -X main.Commit=$COMMIT -X main.BuildTime=$BUILD_TIME" \
        -o "$BUILD_DIR/$BINARY_NAME" \
        ./"$CMD_DIR"

    print_success "Built (dev): $BUILD_DIR/$BINARY_NAME"
}

build_release() {
    print_info "Building release version..."
    mkdir -p "$BUILD_DIR"

    CGO_ENABLED=0 go build \
        -ldflags "-X main.Version=$VERSION -X main.Commit=$COMMIT -X main.BuildTime=$BUILD_TIME -s -w" \
        -a -installsuffix cgo \
        -o "$BUILD_DIR/$BINARY_NAME" \
        ./"$CMD_DIR"

    print_success "Built (release): $BUILD_DIR/$BINARY_NAME"
    ls -lh "$BUILD_DIR/$BINARY_NAME"
}

clean() {
    print_info "Cleaning..."
    go clean
    rm -rf "$BUILD_DIR"
    print_success "Cleaned"
}

run_tests() {
    print_info "Running tests..."
    go test -v ./...
    print_success "Tests completed"
}

download_deps() {
    print_info "Downloading dependencies..."
    go mod download
    print_success "Dependencies downloaded"
}

tidy_deps() {
    print_info "Tidying dependencies..."
    go mod tidy
    go mod verify
    print_success "Dependencies tidied"
}

run_binary() {
    build
    print_info "Running $BINARY_NAME..."
    "./$BUILD_DIR/$BINARY_NAME" --help
}

uninstall() {
    print_info "Uninstalling $BINARY_NAME..."
    sudo rm -f "$INSTALL_DIR/$BINARY_NAME"
    print_success "Uninstalled"
}

check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21 or later."
        exit 1
    fi

    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_info "Using Go version: $GO_VERSION"
}

# Main
main() {
    check_go

    # If no arguments, default to build
    if [ $# -eq 0 ]; then
        build
        exit 0
    fi

    # Process each command
    while [ $# -gt 0 ]; do
        case "$1" in
            build)
                build
                ;;
            install)
                install
                ;;
            dev)
                build_dev
                ;;
            release)
                build_release
                ;;
            clean)
                clean
                ;;
            test)
                run_tests
                ;;
            deps)
                download_deps
                ;;
            tidy)
                tidy_deps
                ;;
            run)
                run_binary
                ;;
            uninstall)
                uninstall
                ;;
            help|--help|-h)
                show_help
                exit 0
                ;;
            *)
                print_error "Unknown command: $1"
                echo ""
                show_help
                exit 1
                ;;
        esac
        shift
    done
}

main "$@"
