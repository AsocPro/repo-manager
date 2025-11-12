# Quick Start Guide

Get started with Repository Manager in 5 minutes!

## 1. Build & Install

```bash
cd repo-manager

# Using build script (recommended)
./build.sh install

# Or using Makefile (if make is installed)
make install

# Or manually
go build -o bin/repo-manager ./cmd/repo-manager
sudo cp bin/repo-manager /usr/local/bin/
```

## 2. Install fzf (Optional but Recommended)

```bash
# Fedora
sudo dnf install fzf

# Ubuntu/Debian
sudo apt install fzf

# macOS
brew install fzf
```

## 3. Configure Your First Provider

Run the tool once to generate the config:

```bash
repo-manager --version
```

Edit `~/.config/repo-manager/config.yaml`:

### For GitLab

```yaml
providers:
  gitlab:
    type: gitlab
    url: https://gitlab.com
    token: "glpat-your-token-here"  # Get from https://gitlab.com/-/profile/personal_access_tokens
    base_dir: /home/youruser/git/bases/gitlab
    enabled: true
```

### For GitHub

```yaml
providers:
  github:
    type: github
    url: https://github.com
    token: "ghp_your-token-here"  # Get from https://github.com/settings/tokens
    base_dir: /home/youruser/git/bases/github
    enabled: true
```

## 4. Set Up Git Credentials

```bash
# Store credentials (recommended for private repos)
git config --global credential.helper store

# Or use SSH (add your key to GitHub/GitLab first)
ssh-keygen -t ed25519 -C "your_email@example.com"
ssh-add ~/.ssh/id_ed25519
```

## 5. Clone Your First Group

```bash
# Clone a GitLab group (interactive selection)
repo-manager clone --provider gitlab --group mygroup --interactive

# Or GitHub organization
repo-manager clone --provider github --group myorg --interactive
```

## 6. Create Worktrees

```bash
# Create worktrees for a new feature across all repos
repo-manager worktree create --branch feature-awesome --all

# Or interactively select which repos
repo-manager worktree create --branch feature-awesome --interactive
```

## 7. Set Up Shell Integration

Add to your `~/.bashrc` or `~/.zshrc`:

```bash
# Quick worktree switch
wt() {
  eval $(repo-manager worktree switch "$@")
}
```

Reload your shell:
```bash
source ~/.bashrc  # or ~/.zshrc
```

## 8. Use Your Worktrees!

```bash
# Switch to any worktree interactively
wt

# Or filter by branch
wt --branch feature-awesome

# List all worktrees
repo-manager worktree list

# List with details
repo-manager worktree list --verbose
```

## Common Workflows

### Single Feature Across Multiple Repos

```bash
# 1. Create worktrees
repo-manager worktree create --branch feature-x --all

# 2. Switch to first repo
wt --branch feature-x
# Make changes, commit, push

# 3. Switch to next repo
wt --branch feature-x
# Make changes, commit, push

# 4. Clean up when done
repo-manager worktree delete --branch feature-x
```

### Multiple Features Simultaneously

```bash
# Create worktrees for different features
repo-manager worktree create --branch feature-a --all
repo-manager worktree create --branch feature-b --all
repo-manager worktree create --branch bugfix-c --all

# Switch between them easily
wt --branch feature-a
wt --branch feature-b
wt --branch bugfix-c

# List all to see what you're working on
repo-manager worktree list
```

### Working with Subsets of Repos

```bash
# Clone and tag frontend repos
repo-manager clone --provider gitlab --group mygroup \
  --filter "frontend-*" \
  --tags frontend

# Clone and tag backend repos
repo-manager clone --provider gitlab --group mygroup \
  --filter "backend-*" \
  --tags backend

# Work on frontend feature
repo-manager worktree create --branch ui-update --tags frontend

# Work on backend feature
repo-manager worktree create --branch api-update --tags backend
```

## Helpful Aliases

Add these to your `~/.bashrc` or `~/.zshrc`:

```bash
# Short alias
alias rm='repo-manager'

# Quick worktree commands
alias w='wt'
alias wl='repo-manager worktree list'
alias wc='repo-manager worktree create'
alias wd='repo-manager worktree delete'

# Create worktree for current branch
wtc() {
  repo-manager worktree create --branch $(git branch --show-current) --all
}
```

## Next Steps

- Read the full [README.md](README.md) for all features
- Check [DEPENDENCIES.md](DEPENDENCIES.md) for library details
- View [PROJECT_PLAN.md](PROJECT_PLAN.md) for implementation status

## Getting Help

```bash
# General help
repo-manager --help

# Command-specific help
repo-manager clone --help
repo-manager worktree --help
repo-manager worktree create --help
```

## Troubleshooting

### "fzf is not installed"
Install fzf or use non-interactive mode:
```bash
repo-manager clone --provider gitlab --group mygroup --repos repo1,repo2
```

### Authentication errors
1. Check your token in `~/.config/repo-manager/config.yaml`
2. Set up git credential helper: `git config --global credential.helper store`
3. Or use SSH: `repo-manager clone --provider gitlab --group mygroup --ssh`

### "provider X is not enabled"
Edit config and set `enabled: true` for the provider.

---

Happy coding! 🚀
