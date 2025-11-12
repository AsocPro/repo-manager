#!/bin/zsh
# Repository Manager - Zsh Integration
# Add this to your ~/.zshrc:
#   source /path/to/repo-manager/scripts/repo-manager.zsh

# Quick worktree switch
wt() {
  eval $(repo-manager worktree switch "$@")
}

# Aliases for common commands
alias rm-clone='repo-manager clone'
alias rm-wt='repo-manager worktree'
alias rm-wt-create='repo-manager worktree create'
alias rm-wt-list='repo-manager worktree list'
alias rm-wt-delete='repo-manager worktree delete'

# Create worktree for current branch across all repos
wt-current() {
  local current_branch=$(git branch --show-current 2>/dev/null)
  if [ -z "$current_branch" ]; then
    echo "Error: Not in a git repository or no branch checked out"
    return 1
  fi
  repo-manager worktree create --branch "$current_branch" "$@"
}

# List worktrees for current branch
wt-current-list() {
  local current_branch=$(git branch --show-current 2>/dev/null)
  if [ -z "$current_branch" ]; then
    echo "Error: Not in a git repository or no branch checked out"
    return 1
  fi
  repo-manager worktree list --branch "$current_branch" "$@"
}

# Delete worktrees for current branch
wt-current-delete() {
  local current_branch=$(git branch --show-current 2>/dev/null)
  if [ -z "$current_branch" ]; then
    echo "Error: Not in a git repository or no branch checked out"
    return 1
  fi
  repo-manager worktree delete --branch "$current_branch" "$@"
}

# Zsh completion for repo-manager
_repo_manager() {
  local -a commands
  commands=(
    'clone:Clone a group or individual repository'
    'worktree:Manage git worktrees'
    'sync:Sync repositories and worktrees'
    'status:Show status of repositories and worktrees'
    'tui:Launch the interactive TUI'
  )

  local -a worktree_commands
  worktree_commands=(
    'create:Create a new worktree'
    'list:List all worktrees'
    'delete:Delete a worktree'
    'switch:Switch to a worktree'
  )

  local -a providers
  providers=(
    'gitlab:GitLab provider'
    'github:GitHub provider'
    'gitea:Gitea/Forgejo provider'
    'generic:Generic Git provider'
  )

  if (( CURRENT == 2 )); then
    _describe -t commands 'repo-manager commands' commands
  elif (( CURRENT == 3 )) && [[ ${words[2]} == "worktree" ]]; then
    _describe -t worktree-commands 'worktree commands' worktree_commands
  elif [[ ${words[CURRENT-1]} == "--provider" ]] || [[ ${words[CURRENT-1]} == "-p" ]]; then
    _describe -t providers 'providers' providers
  else
    _arguments \
      '--config[Config file]:file:_files' \
      '--help[Show help]' \
      '--version[Show version]'
  fi
}

compdef _repo_manager repo-manager
