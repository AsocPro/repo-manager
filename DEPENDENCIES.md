# Repository Manager - Dependencies Analysis

## Overview
This document explains the key dependencies and libraries that will be used in the repository manager project, along with the rationale for each choice.

## Core Dependencies

### 1. CLI Framework

#### Option A: **github.com/spf13/cobra** (RECOMMENDED)
**Why chosen:**
- Industry standard for Go CLI applications (used by kubectl, hugo, docker CLI)
- Excellent documentation and community support
- Built-in command structure, flag parsing, and help generation
- Subcommand support fits our hybrid CLI/TUI model perfectly
- Easy to generate shell completions (bash, zsh, fish)
- Integrates well with viper for configuration

**Pros:**
- Mature and battle-tested
- Great documentation
- Active maintenance
- Automatic help generation
- Nested subcommands

**Cons:**
- Slightly more complex than simpler alternatives
- May be overkill for very simple CLIs (not our case)

#### Option B: github.com/urfave/cli/v2
**Why not chosen:**
- Simpler than cobra, but we need the complexity for our feature set
- Less commonly used in modern projects
- Good alternative if we want something lighter

**Recommendation:** Use **cobra** for robust CLI structure.

---

### 2. Configuration Management

#### **gopkg.in/yaml.v3** (RECOMMENDED)
**Why chosen:**
- Pure Go YAML parser
- Human-readable format (important for manual editing)
- Supports comments in YAML files
- Good error messages for parsing issues
- Standard library quality

**Alternative considered:** github.com/spf13/viper
- Powerful config management with multiple formats
- Can watch config files for changes
- May be overkill since we just need YAML
- Would pair well with cobra if we want advanced features

**Recommendation:** Start with **yaml.v3** for simplicity. Consider adding **viper** later if we need config watching or multiple formats.

---

### 3. Git Operations

#### Option A: **os/exec + git CLI** (RECOMMENDED)
**Why chosen:**
- Directly uses system git binary
- No compatibility issues with git features
- Simpler to implement for git-worktree operations
- Git worktree commands are complex and well-tested in git itself
- Smaller binary size
- Works with user's existing git config and credentials

**Pros:**
- Always compatible with latest git features
- Leverages user's git configuration
- Simple to implement
- Smaller dependencies

**Cons:**
- Requires git to be installed
- Need to parse git output
- Error handling requires parsing stderr

#### Option B: github.com/go-git/go-git/v5
**Why not chosen:**
- Pure Go implementation of git
- Doesn't support all git features (git-worktree support is limited)
- Larger binary size
- More complex for worktree operations
- Still developing some features

**Recommendation:** Use **os/exec + git CLI** for reliability with worktrees.

---

### 4. Git Provider APIs

#### **GitHub**: github.com/google/go-github/v57
**Why chosen:**
- Official Go client library maintained by Google
- Comprehensive API coverage
- Well documented
- Active maintenance
- Type-safe API

#### **GitLab**: github.com/xanzy/go-gitlab
**Why chosen:**
- De facto standard Go client for GitLab
- Covers both GitLab.com and self-hosted
- Actively maintained
- Comprehensive API coverage
- Supports both API v4

#### **Gitea/Forgejo**: code.gitea.io/sdk/gitea
**Why chosen:**
- Official SDK from Gitea
- Works with both Gitea and Forgejo (API compatible)
- Smaller and simpler than GitHub/GitLab clients
- Direct support from Gitea team

**Authentication**: All three support token-based auth which fits our hybrid model.

---

### 5. TUI (Terminal User Interface)

#### Option A: **github.com/charmbracelet/bubbletea** (RECOMMENDED)
**Why chosen:**
- Modern, elegant TUI framework based on Elm architecture
- Active development and great community
- Works well with other charm tools (bubbles for components, lipgloss for styling)
- Clean separation of model/view/update
- Great for complex interactive applications

**Pros:**
- Modern architecture
- Beautiful default styling
- Composable components
- Active development
- Excellent documentation

**Cons:**
- Learning curve if unfamiliar with Elm architecture
- Newer than alternatives

#### Option B: github.com/rivo/tview
**Why not chosen:**
- More traditional widget-based approach
- Mature and stable
- Less modern feeling
- Heavier API

**Companion Libraries:**
- **github.com/charmbracelet/bubbles**: Pre-built components (list, table, spinner, etc.)
- **github.com/charmbracelet/lipgloss**: Styling and layout

**Recommendation:** Use **bubbletea + bubbles + lipgloss** for modern, maintainable TUI.

---

### 6. fzf Integration

#### Option A: **os/exec + fzf CLI** (RECOMMENDED)
**Why chosen:**
- Users likely already have fzf installed
- Leverages native fzf with all features
- Simple integration via stdin/stdout
- Full fzf feature set available

**Pros:**
- Uses native fzf installation
- All fzf features available
- Simple implementation
- Users familiar with fzf behavior

**Cons:**
- Requires fzf to be installed
- External dependency

#### Option B: github.com/ktr0731/go-fuzzyfinder
**Why not chosen:**
- Pure Go fuzzy finder
- No external dependency
- Less feature-complete than fzf
- Different keybindings than standard fzf

**Recommendation:** Use **os/exec + fzf** for familiarity. Gracefully detect if fzf is missing and suggest installation.

---

## Utility Libraries

### 7. **github.com/mitchellh/go-homedir**
**Purpose:** Cross-platform home directory detection
**Why:** Simpler than os.UserHomeDir for our use case, handles edge cases

**Alternative:** Standard library `os.UserHomeDir()` works fine too

---

### 8. **github.com/fatih/color** (Optional)
**Purpose:** Colored terminal output
**Why:** Nice for status messages, warnings, errors
**Note:** bubbletea/lipgloss covers this in TUI mode

---

### 9. **github.com/spf13/afero** (For Testing)
**Purpose:** Filesystem abstraction for testing
**Why:** Makes testing file operations much easier without touching real filesystem

---

## Summary of Recommended Stack

```go
require (
    // CLI Framework
    github.com/spf13/cobra v1.8.0

    // Configuration
    gopkg.in/yaml.v3 v3.0.1

    // Git Provider APIs
    github.com/google/go-github/v57 v57.0.0
    github.com/xanzy/go-gitlab v0.95.0
    code.gitea.io/sdk/gitea v0.17.0

    // TUI
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/bubbles v0.18.0
    github.com/charmbracelet/lipgloss v0.9.1

    // Utilities
    github.com/mitchellh/go-homedir v1.1.0
    github.com/fatih/color v1.16.0

    // Testing
    github.com/spf13/afero v1.11.0
    github.com/stretchr/testify v1.8.4
)
```

## External Tools Required

1. **git** - Required for all git operations
2. **fzf** - Recommended for interactive selection (we'll detect and suggest installation)

## Decision Matrix

| Need | Library | Alternative | Reason |
|------|---------|-------------|---------|
| CLI | cobra | urfave/cli | More features, better docs, industry standard |
| Config | yaml.v3 | viper | Simpler, sufficient for our needs |
| Git Ops | os/exec | go-git | Better worktree support, uses user's git config |
| TUI | bubbletea | tview | Modern, better UX, active development |
| fzf | os/exec | go-fuzzyfinder | Users expect standard fzf behavior |
| GitHub | go-github | - | Official, comprehensive |
| GitLab | go-gitlab | - | De facto standard |
| Gitea | gitea SDK | - | Official SDK |

## Design Philosophy

1. **Leverage existing tools**: Use git and fzf binaries rather than reimplementing
2. **Modern UX**: Choose libraries that provide great user experience (bubbletea)
3. **Reliability**: Prefer battle-tested libraries (cobra, provider SDKs)
4. **Maintainability**: Choose actively maintained projects
5. **Type safety**: Prefer libraries with good Go idioms and type safety

## Next Steps

Once you approve this dependency list, we'll:
1. Initialize go.mod with these dependencies
2. Set up the basic project structure
3. Begin Phase 1 implementation
