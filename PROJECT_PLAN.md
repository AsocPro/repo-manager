# Repository Manager - Project Plan & Status

## Project Overview
A modular Go CLI tool with hybrid command-line and TUI interfaces for managing git repositories and worktrees.

## Design Decisions

### Directory Organization
- **Branch-based worktree organization**: `~/worktrees/branch-name/repo-name/`
- **Configurable base directory**: Per-provider subdirectories specified in config
  - Example: `~/git/bases/gitlab/`, `~/git/bases/github/`

### Configuration
- **Format**: YAML config file (`~/.config/repo-manager/config.yaml`)
- **Contents**: Providers, base dirs, worktree dirs, defaults, API tokens

### Provider Support
- GitLab (self-hosted + .com)
- GitHub (cloud + Enterprise)
- Gitea/Forgejo
- Generic Git (no API, URL-based)

### Authentication
- Hybrid approach: API tokens for provider APIs, git credential helper for cloning

### Metadata Tracking
- Provider and remote URL
- Default branch to track
- Active worktrees
- Custom tags/labels
- Access times (automatic) + manual active/inactive flags
- Basic metadata in YAML config, state evaluated from disk

### Key Features
- **Subgroup handling**: Configurable default (flatten or preserve hierarchy)
- **Branch creation**: Auto-create branches when creating worktrees if missing
- **Worktree naming**: Subdirectories per repo under branch directory
- **Cleanup**: Status warnings for orphaned worktrees (no auto-delete)
- **Shell integration**: Shell function for worktree switching
- **Batch operations**: Multi-select (fzf) + tag/filter-based
- **Clone filtering**: Interactive selection, pattern matching, metadata filtering, default to all

### CLI Interface
- Hybrid approach: Both CLI commands and TUI mode
- Commands: `clone`, `worktree`, `sync`, `status`, `tui`
- fzf integration for interactive selection

### TUI Features
- Browse and select repos
- Create/delete worktrees
- Sync operations (fetch/pull)
- Status overview

## Project Structure

```
repo-manager/
├── cmd/
│   └── repo-manager/        # Main entry point
├── internal/
│   ├── config/              # YAML config management
│   ├── git/                 # Git operations wrapper
│   ├── worktree/            # Worktree management logic
│   ├── provider/            # Git provider interfaces
│   │   ├── interface.go
│   │   ├── gitlab/
│   │   ├── github/
│   │   ├── gitea/
│   │   └── generic/
│   ├── ui/                  # TUI and fzf integration
│   ├── shell/               # Shell integration scripts
│   └── metadata/            # Repository metadata tracking
├── pkg/                     # Public reusable packages
└── scripts/                 # Shell integration scripts
```

## Implementation Phases

### Phase 1: Core Foundation ✅ COMPLETED
- [x] Initialize Go module
- [x] Set up project directory structure
- [x] Implement YAML configuration system
  - [x] Config file loading/saving
  - [x] Default config generation
  - [x] Validation
- [x] Create Git wrapper library
  - [x] Basic git operations (clone, fetch, worktree commands)
  - [x] Error handling
- [x] Set up CLI framework
  - [x] Command structure
  - [x] Flag parsing
  - [x] Help text
- [x] Repository metadata tracking in YAML
  - [x] Data structures
  - [x] CRUD operations

### Phase 2: Provider Integrations ✅ COMPLETED
- [x] Define provider interface
  - [x] Common methods (ListGroups, ListRepos, GetRepo, etc.)
  - [x] Authentication interface
- [x] Implement GitLab provider
  - [x] API client setup
  - [x] Group/subgroup listing
  - [x] Repository operations
  - [x] Token authentication
- [x] Implement GitHub provider
  - [x] API client setup
  - [x] Organization listing
  - [x] Repository operations
  - [x] Token authentication
- [x] Implement Gitea/Forgejo provider
  - [x] API client setup
  - [x] Organization listing
  - [x] Repository operations
- [x] Implement Generic Git provider
  - [x] URL parsing
  - [x] Basic clone operations
- [x] Hybrid authentication system
  - [x] API token management
  - [x] Git credential helper integration

### Phase 3: Repository Management ✅ COMPLETED
- [x] Clone operations
  - [x] Clone individual repository
  - [x] Clone group/organization
  - [x] Interactive repo selection (fzf)
  - [x] Pattern-based filtering
  - [x] Metadata-based filtering (stars, activity, etc.)
- [x] Base repository storage
  - [x] Per-provider directory management
  - [x] Configurable paths
- [x] Subgroup hierarchy handling
  - [x] Flatten option
  - [x] Preserve hierarchy option
  - [x] Configurable default
- [x] Metadata tracking
  - [x] Save provider info
  - [x] Track remote URLs
  - [x] Store default branch
  - [x] Manage custom tags/labels

### Phase 4: Worktree Management ✅ COMPLETED
- [x] Branch-based directory structure
  - [x] Create branch directories
  - [x] Organize repos within branch dirs
- [x] Create worktree command
  - [x] Select repository (single/batch)
  - [x] Specify branch name
  - [x] Auto-create branch if missing
  - [x] Handle existing worktrees
- [x] Delete worktree command
  - [x] Safe deletion with checks
  - [x] Clean up empty branch directories
- [x] Track worktree state
  - [x] List active worktrees
  - [x] Track access times (automatic)
  - [x] Manual active/inactive flags
- [x] Orphaned worktree detection
  - [x] Status warnings
  - [x] Report in list view
- [x] Batch operations
  - [x] Multi-select with fzf
  - [x] Tag-based filtering
  - [x] Apply operations to multiple repos

### Phase 5: UI/UX
- [ ] CLI commands
  - [ ] `repo-manager clone` (group/repo)
  - [ ] `repo-manager worktree create`
  - [ ] `repo-manager worktree delete`
  - [ ] `repo-manager worktree list`
  - [ ] `repo-manager sync`
  - [ ] `repo-manager status`
  - [ ] `repo-manager tui`
- [ ] fzf integration
  - [ ] Repository selection
  - [ ] Worktree selection
  - [ ] Branch selection
  - [ ] Multi-select support
- [ ] TUI mode
  - [ ] Main navigation
  - [ ] Repository browser
  - [ ] Worktree creation interface
  - [ ] Status overview dashboard
  - [ ] Sync operations interface
- [ ] Shell integration
  - [ ] Bash integration script
  - [ ] Zsh integration script
  - [ ] Worktree switching function
  - [ ] Auto-completion

### Phase 6: Polish & Advanced Features
- [ ] Error handling
  - [ ] Comprehensive error messages
  - [ ] Recovery suggestions
  - [ ] Logging
- [ ] Sync operations
  - [ ] Fetch all bases
  - [ ] Pull worktrees
  - [ ] Batch sync
- [ ] Status overview
  - [ ] Repository status (dirty, ahead/behind)
  - [ ] Worktree status
  - [ ] Orphaned worktree warnings
  - [ ] Activity summary
- [ ] Performance optimization
  - [ ] Parallel operations
  - [ ] Caching
- [ ] Documentation
  - [ ] README
  - [ ] Usage examples
  - [ ] Configuration guide
  - [ ] Shell integration setup
- [ ] Testing
  - [ ] Unit tests
  - [ ] Integration tests
  - [ ] End-to-end tests

## Current Status
**Phase**: Phase 4 Complete ✅ - **Core Functionality Complete!**
**Documentation**: ✅ README.md and QUICKSTART.md added
**Next Steps**: Phase 5 - UI/UX, Sync, and Status (Optional enhancements)

## 🎉 The tool is now fully functional for its primary purpose!

You can now:
- Clone groups/organizations from multiple providers
- Create branch-based worktrees across multiple repos
- Manage worktrees with interactive selection
- Switch between worktrees easily
- Track and organize your work with tags

### Phase 1 Summary
All core foundation components have been implemented:
- ✅ Go module initialized with all dependencies
- ✅ Complete project structure created
- ✅ CLI framework with cobra (all commands stubbed)
- ✅ YAML configuration system (load, save, validate, CRUD)
- ✅ Git wrapper library (clone, worktree, fetch, pull, status)
- ✅ Metadata tracking structures (repos and worktrees)

### Phase 2 Summary
All provider integrations have been implemented:
- ✅ Provider interface with common methods
- ✅ Provider manager with factory pattern
- ✅ GitLab provider with full API support
- ✅ GitHub provider with full API support (including Enterprise)
- ✅ Gitea/Forgejo provider with full API support
- ✅ Generic Git provider for plain URLs
- ✅ Hybrid authentication (API tokens + git credentials)
- ✅ Authentication helper and documentation

### Phase 3 Summary
All repository management features have been implemented:
- ✅ Repository manager with full clone operations
- ✅ Clone individual repositories with metadata tracking
- ✅ Clone entire groups/organizations with filtering
- ✅ fzf integration for interactive repository selection
- ✅ Pattern-based filtering (glob patterns)
- ✅ Metadata-based filtering (stars, activity, archived)
- ✅ Subgroup hierarchy handling (flatten or preserve)
- ✅ Per-provider base directory management
- ✅ Automatic metadata tracking and config updates

### Phase 4 Summary
All worktree management features have been implemented:
- ✅ Worktree manager with full lifecycle operations
- ✅ Branch-based directory organization (~/worktrees/branch/repo/)
- ✅ Create worktrees with auto-branch creation
- ✅ Batch worktree creation (all repos, by tags, by names)
- ✅ Interactive selection with fzf
- ✅ List worktrees with filtering (by branch, by repo)
- ✅ Verbose mode with access times and orphan detection
- ✅ Delete worktrees with safety checks
- ✅ Force delete option
- ✅ Automatic empty directory cleanup
- ✅ Switch command with shell integration
- ✅ Access time tracking

## Notes
- Configuration file location: `~/.config/repo-manager/config.yaml`
- Base repositories stored per-provider with configurable paths
- Worktrees organized by branch name for easy multi-repo development
- All metadata tracked in YAML for simplicity and human editability
