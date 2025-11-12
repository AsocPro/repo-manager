#!/bin/bash
# Repository Manager - Bash Integration
# Add this to your ~/.bashrc:
#   source /path/to/repo-manager/scripts/repo-manager.bash

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

# Bash completion for repo-manager (basic)
_repo_manager_complete() {
  local cur prev commands
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"
  commands="clone worktree sync status tui"

  if [ $COMP_CWORD -eq 1 ]; then
    COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
    return 0
  fi

  case "${prev}" in
    worktree)
      COMPREPLY=( $(compgen -W "create list delete switch" -- ${cur}) )
      return 0
      ;;
    --provider|-p)
      COMPREPLY=( $(compgen -W "gitlab github gitea generic" -- ${cur}) )
      return 0
      ;;
  esac
}

complete -F _repo_manager_complete repo-manager
