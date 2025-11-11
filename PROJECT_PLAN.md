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

### Phase 4: Worktree Management
- [ ] Branch-based directory structure
  - [ ] Create branch directories
  - [ ] Organize repos within branch dirs
- [ ] Create worktree command
  - [ ] Select repository (single/batch)
  - [ ] Specify branch name
  - [ ] Auto-create branch if missing
  - [ ] Handle existing worktrees
- [ ] Delete worktree command
  - [ ] Safe deletion with checks
  - [ ] Clean up empty branch directories
- [ ] Track worktree state
  - [ ] List active worktrees
  - [ ] Track access times (automatic)
  - [ ] Manual active/inactive flags
- [ ] Orphaned worktree detection
  - [ ] Status warnings
  - [ ] Report in status view
- [ ] Batch operations
  - [ ] Multi-select with fzf
  - [ ] Tag-based filtering
  - [ ] Apply operations to multiple repos

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
**Phase**: Phase 3 Complete ✅
**Next Steps**: Begin Phase 4 - Worktree Management

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

## Notes
- Configuration file location: `~/.config/repo-manager/config.yaml`
- Base repositories stored per-provider with configurable paths
- Worktrees organized by branch name for easy multi-repo development
- All metadata tracked in YAML for simplicity and human editability
