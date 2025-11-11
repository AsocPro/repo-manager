package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Client wraps git operations
type Client struct{}

// New creates a new Git client
func New() *Client {
	return &Client{}
}

// Clone clones a repository to the specified path
func (c *Client) Clone(url, destPath string) error {
	cmd := exec.Command("git", "clone", url, destPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clone repository: %w\nOutput: %s", err, output)
	}
	return nil
}

// CloneBare clones a bare repository to the specified path
func (c *Client) CloneBare(url, destPath string) error {
	cmd := exec.Command("git", "clone", "--bare", url, destPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to clone bare repository: %w\nOutput: %s", err, output)
	}
	return nil
}

// WorktreeAdd creates a new worktree at the specified path for the given branch
func (c *Client) WorktreeAdd(repoPath, worktreePath, branch string, createBranch bool) error {
	args := []string{"-C", repoPath, "worktree", "add"}

	if createBranch {
		args = append(args, "-b", branch)
	}

	args = append(args, worktreePath)

	if !createBranch {
		args = append(args, branch)
	}

	cmd := exec.Command("git", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to add worktree: %w\nOutput: %s", err, output)
	}
	return nil
}

// WorktreeList lists all worktrees for a repository
func (c *Client) WorktreeList(repoPath string) ([]Worktree, error) {
	cmd := exec.Command("git", "-C", repoPath, "worktree", "list", "--porcelain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w\nOutput: %s", err, output)
	}

	return parseWorktreeList(string(output)), nil
}

// Worktree represents a git worktree
type Worktree struct {
	Path   string
	Branch string
	Commit string
	Bare   bool
}

// parseWorktreeList parses the output of git worktree list --porcelain
func parseWorktreeList(output string) []Worktree {
	var worktrees []Worktree
	var current Worktree

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if current.Path != "" {
				worktrees = append(worktrees, current)
				current = Worktree{}
			}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			current.Path = strings.TrimPrefix(line, "worktree ")
		} else if strings.HasPrefix(line, "branch ") {
			current.Branch = strings.TrimPrefix(line, "branch ")
			current.Branch = strings.TrimPrefix(current.Branch, "refs/heads/")
		} else if strings.HasPrefix(line, "HEAD ") {
			current.Commit = strings.TrimPrefix(line, "HEAD ")
		} else if line == "bare" {
			current.Bare = true
		}
	}

	// Add last worktree if any
	if current.Path != "" {
		worktrees = append(worktrees, current)
	}

	return worktrees
}

// WorktreeRemove removes a worktree
func (c *Client) WorktreeRemove(repoPath, worktreePath string, force bool) error {
	args := []string{"-C", repoPath, "worktree", "remove"}
	if force {
		args = append(args, "--force")
	}
	args = append(args, worktreePath)

	cmd := exec.Command("git", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove worktree: %w\nOutput: %s", err, output)
	}
	return nil
}

// Fetch fetches updates from the remote repository
func (c *Client) Fetch(repoPath string, prune bool) error {
	args := []string{"-C", repoPath, "fetch"}
	if prune {
		args = append(args, "--prune")
	}

	cmd := exec.Command("git", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to fetch: %w\nOutput: %s", err, output)
	}
	return nil
}

// Pull pulls updates from the remote repository
func (c *Client) Pull(repoPath string) error {
	cmd := exec.Command("git", "-C", repoPath, "pull")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to pull: %w\nOutput: %s", err, output)
	}
	return nil
}

// Status returns the status of the repository
func (c *Client) Status(repoPath string) (*Status, error) {
	cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain", "--branch")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w\nOutput: %s", err, output)
	}

	return parseStatus(string(output)), nil
}

// Status represents repository status
type Status struct {
	Branch      string
	Tracking    string
	Ahead       int
	Behind      int
	Clean       bool
	Modified    []string
	Untracked   []string
}

// parseStatus parses git status --porcelain --branch output
func parseStatus(output string) *Status {
	status := &Status{
		Clean: true,
	}

	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if i == 0 && strings.HasPrefix(line, "##") {
			// Parse branch line
			branchInfo := strings.TrimPrefix(line, "## ")
			parts := strings.Split(branchInfo, "...")
			if len(parts) > 0 {
				status.Branch = parts[0]
			}
			if len(parts) > 1 {
				tracking := parts[1]
				if strings.Contains(tracking, "[") {
					// Parse ahead/behind info
					status.Tracking = strings.Split(tracking, " [")[0]
					// TODO: Parse ahead/behind numbers
				} else {
					status.Tracking = tracking
				}
			}
			continue
		}

		if len(line) < 3 {
			continue
		}

		status.Clean = false
		statusCode := line[0:2]
		filePath := line[3:]

		if strings.Contains(statusCode, "?") {
			status.Untracked = append(status.Untracked, filePath)
		} else {
			status.Modified = append(status.Modified, filePath)
		}
	}

	return status
}

// GetDefaultBranch gets the default branch of the repository
func (c *Client) GetDefaultBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "symbolic-ref", "refs/remotes/origin/HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Fallback: try to detect main or master
		branches := []string{"main", "master"}
		for _, branch := range branches {
			if c.BranchExists(repoPath, branch) {
				return branch, nil
			}
		}
		return "", fmt.Errorf("failed to determine default branch: %w", err)
	}

	// Parse refs/remotes/origin/HEAD -> refs/remotes/origin/main
	ref := strings.TrimSpace(string(output))
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1], nil
	}

	return "", fmt.Errorf("failed to parse default branch")
}

// BranchExists checks if a branch exists
func (c *Client) BranchExists(repoPath, branch string) bool {
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--verify", fmt.Sprintf("refs/heads/%s", branch))
	err := cmd.Run()
	return err == nil
}

// GetCurrentBranch returns the current branch name
func (c *Client) GetCurrentBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "branch", "--show-current")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(out.String()), nil
}

// GetRemoteURL returns the remote URL of the repository
func (c *Client) GetRemoteURL(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "remote", "get-url", "origin")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get remote URL: %w", err)
	}
	return strings.TrimSpace(out.String()), nil
}

// EnsureDir ensures a directory exists, creating it if necessary
func EnsureDir(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	return exec.Command("mkdir", "-p", absPath).Run()
}
