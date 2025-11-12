# Repository Manager - File Structure

Complete list of all files in the project.

## Build System

- `Makefile` - Build automation with make
- `build.sh` - Build script (portable alternative to Makefile)
- `.gitignore` - Git ignore patterns

## Documentation

- `README.md` - Main documentation with features and usage
- `QUICKSTART.md` - 5-minute quick start guide
- `BUILD.md` - Complete build instructions
- `DEPENDENCIES.md` - Dependency analysis and rationale
- `PROJECT_PLAN.md` - Implementation tracking and status
- `SUMMARY.md` - Project overview and statistics
- `COMMIT_MESSAGE.md` - Git commit message template
- `FILES.md` - This file

## Source Code

### Main Entry Point
- `cmd/repo-manager/main.go` - Application entry point

### CLI Commands
- `internal/cmd/root.go` - Root command and initialization
- `internal/cmd/clone.go` - Clone command implementation
- `internal/cmd/worktree.go` - Worktree commands (create/list/delete/switch)
- `internal/cmd/sync.go` - Sync command (stub)
- `internal/cmd/status.go` - Status command (stub)
- `internal/cmd/tui.go` - TUI command (stub)

### Core Libraries
- `internal/config/config.go` - YAML configuration management
- `internal/git/git.go` - Git operations wrapper
- `internal/metadata/worktree.go` - Worktree metadata tracking

### Provider System
- `internal/provider/interface.go` - Provider interface definition
- `internal/provider/auth.go` - Authentication helper
- `internal/provider/gitlab/gitlab.go` - GitLab provider
- `internal/provider/github/github.go` - GitHub provider
- `internal/provider/gitea/gitea.go` - Gitea/Forgejo provider
- `internal/provider/generic/generic.go` - Generic Git provider

### Repository Management
- `internal/repository/manager.go` - Repository operations
- `internal/repository/provider_manager.go` - Provider factory

### Worktree Management
- `internal/worktree/manager.go` - Worktree operations

### UI Components
- `internal/ui/fzf.go` - fzf integration

## Shell Integration
- `scripts/repo-manager.bash` - Bash completion and functions
- `scripts/repo-manager.zsh` - Zsh completion and functions

## Go Module Files
- `go.mod` - Go module definition
- `go.sum` - Dependency checksums

## Build Output (Not Committed)
- `bin/repo-manager` - Compiled binary (13MB)

## Total Statistics

- **Go Source Files**: 25
- **Documentation Files**: 8
- **Shell Scripts**: 3
- **Total Lines of Code**: ~3,500+
- **Total Lines of Documentation**: ~1,500+
- **Binary Size**: 13MB (debug), ~10MB (release)

## Package Structure

```
github.com/asocpro/repo-manager/
├── cmd/
│   └── repo-manager/         # Main entry point (1 file)
├── internal/
│   ├── cmd/                  # CLI commands (6 files)
│   ├── config/               # Configuration (1 file)
│   ├── git/                  # Git wrapper (1 file)
│   ├── metadata/             # Metadata tracking (1 file)
│   ├── provider/             # Providers (6 files)
│   ├── repository/           # Repository management (2 files)
│   ├── ui/                   # UI components (1 file)
│   └── worktree/             # Worktree management (1 file)
├── scripts/                  # Shell integration (2 files)
└── pkg/                      # Public packages (empty, reserved)
```

## File Sizes (Approximate)

```
Source Code:
- Go files: ~3,500 lines
- Shell scripts: ~200 lines

Documentation:
- README.md: ~400 lines
- QUICKSTART.md: ~150 lines
- BUILD.md: ~300 lines
- DEPENDENCIES.md: ~200 lines
- PROJECT_PLAN.md: ~250 lines
- SUMMARY.md: ~200 lines

Total: ~5,200 lines across all files
```

## Dependencies (External)

See `DEPENDENCIES.md` for full details. Key dependencies:

- `github.com/spf13/cobra` - CLI framework
- `gitlab.com/gitlab-org/api/client-go` - GitLab API
- `github.com/google/go-github/v57` - GitHub API
- `code.gitea.io/sdk/gitea` - Gitea API
- `gopkg.in/yaml.v3` - YAML parsing
- `github.com/charmbracelet/bubbletea` - TUI (future)
- `github.com/charmbracelet/bubbles` - TUI components (future)
- `github.com/charmbracelet/lipgloss` - TUI styling (future)

Total: 15+ direct dependencies, 50+ transitive dependencies
