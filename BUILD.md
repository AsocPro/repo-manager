# Building Repository Manager

Complete guide for building, installing, and developing Repository Manager.

## Quick Build

```bash
# Clone the repository
git clone https://github.com/asocpro/repo-manager
cd repo-manager

# Build and install (choose one)
./build.sh install    # Using build script
make install          # Using Makefile
```

## Prerequisites

### Required
- **Go 1.21+** - [Download](https://go.dev/doc/install)
- **Git** - Usually pre-installed

### Optional but Recommended
- **fzf** - For interactive selection
- **make** - For using Makefile (alternative: use build.sh)

### Installing Prerequisites

#### Go
```bash
# Fedora
sudo dnf install golang

# Ubuntu/Debian
sudo apt install golang-go

# macOS
brew install go

# Or download from https://go.dev/doc/install
```

#### fzf
```bash
# Fedora
sudo dnf install fzf

# Ubuntu/Debian
sudo apt install fzf

# macOS
brew install fzf

# From source
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf
~/.fzf/install
```

## Build Methods

### Method 1: Build Script (Recommended)

The build script provides a simple, portable way to build the project.

```bash
# Show all commands
./build.sh help

# Basic build
./build.sh build

# Build and install to /usr/local/bin
./build.sh install

# Build optimized release version
./build.sh release

# Build with race detector (development)
./build.sh dev

# Clean build artifacts
./build.sh clean

# Run tests
./build.sh test

# Tidy dependencies
./build.sh tidy
```

**Chaining commands:**
```bash
./build.sh clean build    # Clean then build
./build.sh tidy test build  # Tidy, test, then build
```

### Method 2: Makefile

If you have `make` installed, you can use the Makefile.

```bash
# Show all targets
make help

# Build
make build

# Build and install
make install

# Build optimized release
make release

# Build with race detector
make dev

# Clean
make clean

# Run tests
make test

# Format code
make fmt

# Tidy dependencies
make tidy
```

### Method 3: Manual Go Commands

```bash
# Simple build
go build -o bin/repo-manager ./cmd/repo-manager

# Build with version info
VERSION=$(git describe --tags --always --dirty)
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

go build \
  -ldflags "-X main.Version=$VERSION -X main.Commit=$COMMIT -X main.BuildTime=$BUILD_TIME" \
  -o bin/repo-manager \
  ./cmd/repo-manager

# Optimized release build (smaller binary, no debug info)
CGO_ENABLED=0 go build \
  -ldflags "-s -w" \
  -a -installsuffix cgo \
  -o bin/repo-manager \
  ./cmd/repo-manager

# Build with race detector (development)
go build -race -o bin/repo-manager ./cmd/repo-manager
```

## Installation

### System-Wide Installation

```bash
# After building, install to /usr/local/bin
sudo cp bin/repo-manager /usr/local/bin/

# Or use the build script
./build.sh install

# Verify installation
which repo-manager
repo-manager --version
```

### Local Installation (User-Only)

```bash
# Install to ~/bin (add ~/bin to PATH if needed)
mkdir -p ~/bin
cp bin/repo-manager ~/bin/

# Add to PATH in ~/.bashrc or ~/.zshrc
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc

# Reload shell
source ~/.bashrc
```

### Install Shell Integration

```bash
# For Bash
echo 'source /path/to/repo-manager/scripts/repo-manager.bash' >> ~/.bashrc
source ~/.bashrc

# For Zsh
echo 'source /path/to/repo-manager/scripts/repo-manager.zsh' >> ~/.zshrc
source ~/.zshrc
```

## Development

### Setting Up Development Environment

```bash
# Clone repository
git clone https://github.com/asocpro/repo-manager
cd repo-manager

# Download dependencies
go mod download

# Or using build tools
./build.sh deps
# or
make deps
```

### Building for Development

```bash
# Build with race detector (detects race conditions)
./build.sh dev
# or
make dev
# or
go build -race -o bin/repo-manager ./cmd/repo-manager
```

### Running Tests

```bash
# Run all tests
./build.sh test
# or
make test
# or
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

### Code Formatting

```bash
# Format code
make fmt
# or
go fmt ./...

# Check formatting
gofmt -l .
```

### Dependency Management

```bash
# Download dependencies
go mod download

# Add a new dependency
go get github.com/example/package

# Remove unused dependencies
go mod tidy

# Verify dependencies
go mod verify

# View dependency graph
go mod graph
```

## Build Configurations

### Debug Build (Default)
- Includes debug symbols
- Larger binary size (~13MB)
- Useful for development

```bash
go build -o bin/repo-manager ./cmd/repo-manager
```

### Release Build
- Stripped symbols (`-s -w`)
- Smaller binary size
- Optimized for production

```bash
./build.sh release
# or
make release
```

### Development Build
- Includes race detector
- Slower execution
- Detects race conditions

```bash
./build.sh dev
# or
make dev
```

## Cross-Compilation

Build for different platforms:

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o bin/repo-manager-linux-amd64 ./cmd/repo-manager

# Linux ARM64
GOOS=linux GOARCH=arm64 go build -o bin/repo-manager-linux-arm64 ./cmd/repo-manager

# macOS AMD64
GOOS=darwin GOARCH=amd64 go build -o bin/repo-manager-darwin-amd64 ./cmd/repo-manager

# macOS ARM64 (M1/M2)
GOOS=darwin GOARCH=arm64 go build -o bin/repo-manager-darwin-arm64 ./cmd/repo-manager

# Windows AMD64
GOOS=windows GOARCH=amd64 go build -o bin/repo-manager-windows-amd64.exe ./cmd/repo-manager
```

## Troubleshooting

### Go Not Found
```bash
# Check if Go is installed
go version

# If not installed, see Prerequisites section above
```

### Permission Denied During Install
```bash
# Use sudo for system-wide installation
sudo cp bin/repo-manager /usr/local/bin/

# Or install to user directory
cp bin/repo-manager ~/bin/
```

### Build Fails with Missing Dependencies
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Try building again
./build.sh clean build
```

### Binary Too Large
```bash
# Build release version (stripped symbols)
./build.sh release

# Check size
ls -lh bin/repo-manager
```

## Build Output

After building, you'll have:

```
repo-manager/
├── bin/
│   └── repo-manager    # Binary (executable)
└── ...
```

Binary size:
- Debug build: ~13 MB
- Release build: ~10 MB

## Uninstalling

```bash
# Remove installed binary
./build.sh uninstall
# or
make uninstall
# or
sudo rm /usr/local/bin/repo-manager

# Remove shell integration (edit ~/.bashrc or ~/.zshrc and remove the source line)

# Remove config (optional)
rm -rf ~/.config/repo-manager
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Build

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Build
        run: ./build.sh build
      - name: Test
        run: ./build.sh test
```

### GitLab CI Example

```yaml
build:
  image: golang:1.21
  script:
    - ./build.sh build
    - ./build.sh test
```

## Additional Resources

- [README.md](README.md) - Full documentation
- [QUICKSTART.md](QUICKSTART.md) - Quick start guide
- [DEPENDENCIES.md](DEPENDENCIES.md) - Dependency details
- [Go Documentation](https://go.dev/doc/)
