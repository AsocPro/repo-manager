# Repository Manager

> WARNING At this point this is a bunch of vibe coded shenanigans that is being used to prototype how I would best want this to work. Its far from robustly implemneted and will not be stable for the current forseeable future.

A powerful CLI tool for managing multiple git repositories and their worktrees, organized by branch names. Perfect for working on the same feature across multiple repositories simultaneously.

## Features

### Core Capabilities

- 🌳 **Branch-Based Worktrees** - Organize worktrees by branch name for multi-repo development
- 📦 **Batch Operations** - Clone entire groups and create worktrees across multiple repos at once
- 🔍 **Interactive Selection** - fzf integration for easy repository and worktree selection
- 🏷️ **Tag System** - Organize and filter repositories with custom tags
- 🔄 **Auto-Branch Creation** - Automatically creates branches when creating worktrees
- 🚀 **Multi-Provider Support** - Works with GitLab, GitHub, Gitea, Forgejo, and generic Git

### Provider Support

| Provider | Features | Notes |
|----------|----------|-------|
| **GitLab** | ✅ Groups & subgroups<br>✅ Project listing<br>✅ Metadata filtering<br>✅ Self-hosted support | Supports both GitLab.com and self-hosted instances |
| **GitHub** | ✅ Organizations<br>✅ Repository listing<br>✅ Metadata filtering<br>✅ Enterprise support | Supports both GitHub.com and GitHub Enterprise |
| **Gitea** | ✅ Organizations<br>✅ Repository listing<br>✅ Metadata filtering | Works with both Gitea and Forgejo |
| **Generic** | ✅ Any Git URL<br>✅ No API required | Works with any Git repository |

## Installation

### Prerequisites

- **Go 1.21 or later** - [Install Go](https://go.dev/doc/install)
- **Git** - Should be already installed
- **fzf** (optional, for interactive selection) - Highly recommended
- **make** (optional) - For using Makefile

### Quick Install

```bash
# Clone the repository
git clone https://github.com/asocpro/repo-manager
cd repo-manager

# Option 1: Using Makefile (recommended)
make install

# Option 2: Using build script
./build.sh install

# Option 3: Manual build
go build -o bin/repo-manager ./cmd/repo-manager
sudo cp bin/repo-manager /usr/local/bin/
```

### Build Options

#### Using Makefile

```bash
# Build only
make build

# Build and install
make install

# Build optimized release version
make release

# Build with race detector (for development)
make dev

# Clean build artifacts
make clean

# Show all available commands
make help
```

#### Using Build Script

```bash
# Build only
./build.sh build

# Build and install
./build.sh install

# Build optimized release
./build.sh release

# Clean build artifacts
./build.sh clean

# Show all available commands
./build.sh help
```

#### Manual Build

```bash
# Simple build
go build -o bin/repo-manager ./cmd/repo-manager

# Build with version information
VERSION=$(git describe --tags --always --dirty)
go build -ldflags "-X main.Version=$VERSION" -o bin/repo-manager ./cmd/repo-manager

# Optimized release build
CGO_ENABLED=0 go build -ldflags "-s -w" -a -installsuffix cgo -o bin/repo-manager ./cmd/repo-manager
```

### Install fzf (Optional but Recommended)

fzf enables interactive selection for repositories and worktrees:

```bash
# Fedora
sudo dnf install fzf

# Ubuntu/Debian
sudo apt install fzf

# macOS
brew install fzf

# Or install from source
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf
~/.fzf/install
```

### Verify Installation

```bash
# Check version
repo-manager --version

# Show help
repo-manager --help

# Test fzf availability
fzf --version
```

## Configuration

The configuration file is automatically created at `~/.config/repo-manager/config.yaml` on first run.

### Basic Configuration

```yaml
version: "1"
providers:
  gitlab:
    type: gitlab
    url: https://gitlab.com
    token: "glpat-xxxxxxxxxxxx"  # Optional: for private repos
    base_dir: /home/user/git/bases/gitlab
    enabled: true

  github:
    type: github
    url: https://github.com
    token: "ghp_xxxxxxxxxxxx"  # Optional: for private repos
    base_dir: /home/user/git/bases/github
    enabled: true

  gitea:
    type: gitea
    url: https://git.example.com  # Required for Gitea
    token: "xxxxxxxxxxxx"
    base_dir: /home/user/git/bases/gitea
    enabled: false

worktree:
  base_dir: /home/user/worktrees

defaults:
  flatten_subgroups: false
  default_branch: main
```

### Getting API Tokens

#### GitLab
1. Go to https://gitlab.com/-/profile/personal_access_tokens
2. Create token with `read_api` and `read_repository` scopes
3. Add to config: `token: "glpat-xxxxxxxxxxxx"`

#### GitHub
1. Go to https://github.com/settings/tokens
2. Generate new token (classic) with `repo` scope
3. Add to config: `token: "ghp_xxxxxxxxxxxx"`

#### Gitea/Forgejo
1. Go to your instance: `https://your-gitea.com/user/settings/applications`
2. Generate new token
3. Add to config: `token: "xxxxxxxxxxxx"`

### Git Credentials

For cloning, the tool uses Git's credential helper:

```bash
# Store credentials (they'll be saved on disk)
git config --global credential.helper store

# Or cache credentials temporarily (1 hour)
git config --global credential.helper "cache --timeout=3600"

# For SSH, set up keys as usual
ssh-keygen -t ed25519 -C "your_email@example.com"
ssh-add ~/.ssh/id_ed25519
# Add public key to your provider
```

## Usage

### Clone Repositories

#### Clone a Single Repository

```bash
# Clone from GitLab
repo-manager clone --provider gitlab --repo mygroup/myrepo

# Clone from GitHub
repo-manager clone --provider github --repo myorg/myrepo

# Use SSH instead of HTTPS
repo-manager clone --provider gitlab --repo mygroup/myrepo --ssh
```

#### Clone an Entire Group/Organization

```bash
# Clone all repos in a GitLab group
repo-manager clone --provider gitlab --group mygroup

# Clone with interactive selection
repo-manager clone --provider github --group myorg --interactive

# Clone with filtering
repo-manager clone --provider gitlab --group mygroup \
  --filter "backend-*" \
  --min-stars 5 \
  --min-activity 2024-01-01 \
  --tags backend,production

# Flatten subgroups (ignore hierarchy)
repo-manager clone --provider gitlab --group mygroup --flatten
```

### Manage Worktrees

#### Create Worktrees

```bash
# Create worktree for all repos on branch feature-x
repo-manager worktree create --branch feature-x --all

# Create for specific repos
repo-manager worktree create --branch bugfix-123 --repos repo1,repo2

# Create with interactive selection
repo-manager worktree create --branch feature-y --interactive

# Create for repos with specific tags
repo-manager worktree create --branch dev --tags frontend,active
```

**Result:** Worktrees are created in branch-based directories:
```
~/worktrees/feature-x/
├── repo1/
├── repo2/
└── repo3/
```

#### List Worktrees

```bash
# List all worktrees
repo-manager worktree list

# List worktrees for a specific branch
repo-manager worktree list --branch feature-x

# List with detailed information
repo-manager worktree list --verbose

# Filter by specific repos
repo-manager worktree list --repos repo1,repo2
```

#### Delete Worktrees

```bash
# Delete all worktrees for a branch
repo-manager worktree delete --branch feature-x

# Delete with interactive selection
repo-manager worktree delete --interactive

# Force delete (skip safety checks)
repo-manager worktree delete --branch old-feature --force

# Delete specific repos' worktrees
repo-manager worktree delete --branch feature-x --repos repo1,repo2
```

#### Switch Between Worktrees

First, add this function to your `~/.bashrc` or `~/.zshrc`:

```bash
wt() {
  eval $(repo-manager worktree switch "$@")
}
```

Then use it:

```bash
# Interactively switch to any worktree
wt

# Switch to a worktree on a specific branch
wt --branch feature-x
```

## Workflows

### Multi-Repo Feature Development

```bash
# 1. Clone all repos in your organization
repo-manager clone --provider github --group myorg --interactive

# 2. Create worktrees for your feature across all repos
repo-manager worktree create --branch feature-awesome --all

# 3. Switch to a specific repo's worktree
wt --branch feature-awesome

# 4. When done, delete all worktrees for the branch
repo-manager worktree delete --branch feature-awesome
```

### Working with Multiple Branches

```bash
# Create worktrees for different branches
repo-manager worktree create --branch feature-x --all
repo-manager worktree create --branch bugfix-y --all
repo-manager worktree create --branch experiment-z --all

# Switch between branches easily
wt --branch feature-x
wt --branch bugfix-y

# List all your active branches
repo-manager worktree list
```

### Organizing with Tags

```bash
# Clone repos and tag them
repo-manager clone --provider gitlab --group mygroup \
  --filter "frontend-*" \
  --tags frontend,active

repo-manager clone --provider gitlab --group mygroup \
  --filter "backend-*" \
  --tags backend,active

# Create worktrees only for frontend repos
repo-manager worktree create --branch ui-redesign --tags frontend

# Create worktrees only for backend repos
repo-manager worktree create --branch api-v2 --tags backend
```

## Advanced Features

### Filtering Options

When cloning groups, you can filter repositories:

- `--filter "pattern"` - Glob pattern for repo names (e.g., `backend-*`)
- `--min-stars N` - Minimum number of stars
- `--min-activity YYYY-MM-DD` - Minimum last activity date
- `--archived` - Include archived repositories (excluded by default)
- `--interactive` - Interactively select repos with fzf

### Subgroup Hierarchy

GitLab groups can have subgroups. Control how they're organized:

```bash
# Preserve hierarchy (default)
repo-manager clone --provider gitlab --group mygroup
# Result: bases/gitlab/mygroup/subgroup/repo

# Flatten (all repos at top level)
repo-manager clone --provider gitlab --group mygroup --flatten
# Result: bases/gitlab/repo
```

Set default behavior in config:
```yaml
defaults:
  flatten_subgroups: false
```

### Access Time Tracking

The tool automatically tracks when you last accessed each worktree:

```bash
repo-manager worktree list --verbose
```

Output shows:
- Last accessed: "2 hours ago", "3 days ago", etc.
- Orphaned worktree warnings (if branch deleted remotely)

## Directory Structure

### Base Repositories

Base repositories (the source of truth) are stored per-provider:

```
~/git/bases/
├── gitlab/
│   ├── repo1/
│   └── mygroup/
│       └── repo2/
├── github/
│   └── repo3/
└── gitea/
    └── repo4/
```

### Worktrees

Worktrees are organized by branch name:

```
~/worktrees/
├── feature-x/
│   ├── repo1/
│   ├── repo2/
│   └── repo3/
├── bugfix-y/
│   ├── repo1/
│   └── repo2/
└── experiment-z/
    └── repo1/
```

## Configuration Reference

### Provider Types

- `gitlab` - GitLab.com or self-hosted
- `github` - GitHub.com or Enterprise
- `gitea` - Gitea or Forgejo
- `generic` - Any Git URL (no API)

### Provider Settings

```yaml
providers:
  name:
    type: gitlab|github|gitea|generic
    url: https://...           # Optional for GitHub/GitLab.com, required for others
    token: "..."               # Optional: API token for private repos
    base_dir: /path/to/bases   # Where to store base repositories
    enabled: true|false        # Enable/disable this provider
```

### Worktree Settings

```yaml
worktree:
  base_dir: /path/to/worktrees  # Where to create branch-based worktrees
```

### Defaults

```yaml
defaults:
  flatten_subgroups: false  # Preserve or flatten group hierarchy
  default_branch: main      # Default branch name
```

## Tips & Tricks

### Quick Navigation

Add these aliases to your shell config:

```bash
# Quick worktree switch
alias w='wt'

# Create worktree for current branch name
wtc() {
  repo-manager worktree create --branch $(git branch --show-current) --all
}

# List worktrees for current branch
wtl() {
  repo-manager worktree list --branch $(git branch --show-current)
}
```

### Clean Up Old Worktrees

```bash
# List all worktrees with verbose info to see old ones
repo-manager worktree list --verbose

# Delete old worktrees interactively
repo-manager worktree delete --interactive
```

### Working with Specific Repos

```bash
# Create config with tags
repo-manager clone --provider gitlab --group mygroup --tags active

# Only work with active repos
repo-manager worktree create --branch feature-x --tags active
```

## Troubleshooting

### fzf not installed

If you see "fzf is not installed":

```bash
# Install fzf
sudo dnf install fzf  # Fedora
sudo apt install fzf  # Ubuntu/Debian
brew install fzf      # macOS
```

Or use non-interactive commands:
```bash
repo-manager clone --provider gitlab --group mygroup --repos repo1,repo2
```

### Authentication Issues

If cloning fails with authentication errors:

1. Check your API token in `~/.config/repo-manager/config.yaml`
2. Set up Git credential helper:
   ```bash
   git config --global credential.helper store
   ```
3. Or use SSH:
   ```bash
   repo-manager clone --provider gitlab --group mygroup --ssh
   ```

### Provider Not Enabled

If you see "provider X is not enabled":

Edit `~/.config/repo-manager/config.yaml`:
```yaml
providers:
  gitlab:
    enabled: true  # Change to true
```

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

[Add your license here]

## Author

AsocPro (asocpro@gmail.com)
