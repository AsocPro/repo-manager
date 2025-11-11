package metadata

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/asocpro/repo-manager/internal/git"
)

// WorktreeInfo represents metadata about a worktree
type WorktreeInfo struct {
	Path         string    `yaml:"path"`
	RepoName     string    `yaml:"repo_name"`
	Branch       string    `yaml:"branch"`
	LastAccessed time.Time `yaml:"last_accessed"`
	Active       bool      `yaml:"active"`
	Orphaned     bool      `yaml:"orphaned"`
}

// WorktreeTracker tracks worktrees and their metadata
type WorktreeTracker struct {
	gitClient *git.Client
}

// NewWorktreeTracker creates a new worktree tracker
func NewWorktreeTracker() *WorktreeTracker {
	return &WorktreeTracker{
		gitClient: git.New(),
	}
}

// DiscoverWorktrees discovers all worktrees in the worktree base directory
func (t *WorktreeTracker) DiscoverWorktrees(worktreeBaseDir string) (map[string][]WorktreeInfo, error) {
	// Map of branch -> []WorktreeInfo
	worktreesByBranch := make(map[string][]WorktreeInfo)

	// List all branch directories
	entries, err := os.ReadDir(worktreeBaseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return worktreesByBranch, nil
		}
		return nil, fmt.Errorf("failed to read worktree directory: %w", err)
	}

	for _, branchEntry := range entries {
		if !branchEntry.IsDir() {
			continue
		}

		branchName := branchEntry.Name()
		branchPath := filepath.Join(worktreeBaseDir, branchName)

		// List all repo directories in this branch
		repoEntries, err := os.ReadDir(branchPath)
		if err != nil {
			continue
		}

		for _, repoEntry := range repoEntries {
			if !repoEntry.IsDir() {
				continue
			}

			repoName := repoEntry.Name()
			worktreePath := filepath.Join(branchPath, repoName)

			// Check if this is actually a git worktree
			if !t.isGitWorktree(worktreePath) {
				continue
			}

			info := WorktreeInfo{
				Path:         worktreePath,
				RepoName:     repoName,
				Branch:       branchName,
				LastAccessed: t.getLastAccessTime(worktreePath),
				Active:       true,
			}

			// Check if the worktree is orphaned (branch doesn't exist remotely)
			info.Orphaned = t.isOrphaned(worktreePath)

			worktreesByBranch[branchName] = append(worktreesByBranch[branchName], info)
		}
	}

	return worktreesByBranch, nil
}

// isGitWorktree checks if a directory is a git worktree
func (t *WorktreeTracker) isGitWorktree(path string) bool {
	gitFile := filepath.Join(path, ".git")
	info, err := os.Stat(gitFile)
	if err != nil {
		return false
	}

	// In a worktree, .git is a file (not a directory) containing a reference
	return !info.IsDir()
}

// getLastAccessTime gets the last access time of a directory
func (t *WorktreeTracker) getLastAccessTime(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}

	// Use modification time as a proxy for last access
	// (atime is often disabled on modern filesystems)
	return info.ModTime()
}

// isOrphaned checks if a worktree's branch has been deleted remotely
func (t *WorktreeTracker) isOrphaned(worktreePath string) bool {
	// Get the current branch
	branch, err := t.gitClient.GetCurrentBranch(worktreePath)
	if err != nil {
		return false
	}

	// Check if the branch exists in remote
	// This is a simple check; a more thorough check would fetch from remote
	status, err := t.gitClient.Status(worktreePath)
	if err != nil {
		return false
	}

	// If there's no tracking branch, it might be orphaned
	return status.Tracking == "" && branch != ""
}

// GetBranchPath returns the path for a branch's worktree directory
func GetBranchPath(worktreeBaseDir, branch string) string {
	return filepath.Join(worktreeBaseDir, branch)
}

// GetWorktreePath returns the path for a specific worktree
func GetWorktreePath(worktreeBaseDir, branch, repoName string) string {
	return filepath.Join(worktreeBaseDir, branch, repoName)
}

// UpdateAccessTime updates the access time for a worktree
func UpdateAccessTime(worktreePath string) error {
	now := time.Now()
	return os.Chtimes(worktreePath, now, now)
}
