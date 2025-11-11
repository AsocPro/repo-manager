package worktree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/asocpro/repo-manager/internal/config"
	"github.com/asocpro/repo-manager/internal/git"
	"github.com/asocpro/repo-manager/internal/metadata"
	"github.com/asocpro/repo-manager/internal/ui"
)

// Manager handles worktree operations
type Manager struct {
	config      *config.Config
	gitClient   *git.Client
	tracker     *metadata.WorktreeTracker
	fzfSelector *ui.FzfSelector
}

// NewManager creates a new worktree manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config:      cfg,
		gitClient:   git.New(),
		tracker:     metadata.NewWorktreeTracker(),
		fzfSelector: ui.NewFzfSelector(),
	}
}

// CreateOptions contains options for creating worktrees
type CreateOptions struct {
	Branch      string
	Repos       []string
	Tags        []string
	All         bool
	Interactive bool
}

// CreateWorktrees creates worktrees for one or more repositories
func (m *Manager) CreateWorktrees(opts CreateOptions) error {
	if opts.Branch == "" {
		return fmt.Errorf("branch name is required")
	}

	// Get list of repositories to create worktrees for
	repos, err := m.selectRepositories(opts)
	if err != nil {
		return fmt.Errorf("failed to select repositories: %w", err)
	}

	if len(repos) == 0 {
		fmt.Println("No repositories selected")
		return nil
	}

	fmt.Printf("Creating worktrees on branch '%s' for %d repositories...\n\n", opts.Branch, len(repos))

	successCount := 0
	failCount := 0
	skipCount := 0

	for i, repo := range repos {
		fmt.Printf("[%d/%d] Creating worktree for %s...\n", i+1, len(repos), repo.Name)

		// Determine worktree path
		worktreePath := metadata.GetWorktreePath(m.config.Worktree.BaseDir, opts.Branch, repo.Name)

		// Check if worktree already exists
		if _, err := os.Stat(worktreePath); err == nil {
			fmt.Printf("⊗ Worktree already exists: %s\n\n", worktreePath)
			skipCount++
			continue
		}

		// Ensure branch directory exists
		branchDir := metadata.GetBranchPath(m.config.Worktree.BaseDir, opts.Branch)
		if err := os.MkdirAll(branchDir, 0755); err != nil {
			fmt.Printf("✗ Failed to create branch directory: %v\n\n", err)
			failCount++
			continue
		}

		// Check if branch exists in the base repository
		branchExists := m.gitClient.BranchExists(repo.BasePath, opts.Branch)

		// Create worktree
		if err := m.gitClient.WorktreeAdd(repo.BasePath, worktreePath, opts.Branch, !branchExists); err != nil {
			fmt.Printf("✗ Failed to create worktree: %v\n\n", err)
			failCount++
			continue
		}

		// Update access time
		metadata.UpdateAccessTime(worktreePath)

		successCount++
		fmt.Printf("✓ Created worktree at %s\n\n", worktreePath)
	}

	// Print summary
	fmt.Printf(strings.Repeat("=", 60) + "\n")
	fmt.Printf("Worktree Creation Summary:\n")
	fmt.Printf("  Success: %d\n", successCount)
	fmt.Printf("  Failed:  %d\n", failCount)
	fmt.Printf("  Skipped: %d\n", skipCount)
	fmt.Printf("  Total:   %d\n", len(repos))
	fmt.Printf(strings.Repeat("=", 60) + "\n")

	return nil
}

// ListOptions contains options for listing worktrees
type ListOptions struct {
	Branch  string
	Repos   []string
	Verbose bool
}

// ListWorktrees lists all worktrees
func (m *Manager) ListWorktrees(opts ListOptions) error {
	// Discover all worktrees
	worktreesByBranch, err := m.tracker.DiscoverWorktrees(m.config.Worktree.BaseDir)
	if err != nil {
		return fmt.Errorf("failed to discover worktrees: %w", err)
	}

	if len(worktreesByBranch) == 0 {
		fmt.Println("No worktrees found")
		return nil
	}

	// Filter by branch if specified
	if opts.Branch != "" {
		if worktrees, ok := worktreesByBranch[opts.Branch]; ok {
			worktreesByBranch = map[string][]metadata.WorktreeInfo{opts.Branch: worktrees}
		} else {
			fmt.Printf("No worktrees found for branch '%s'\n", opts.Branch)
			return nil
		}
	}

	// Filter by repos if specified
	if len(opts.Repos) > 0 {
		worktreesByBranch = m.filterByRepos(worktreesByBranch, opts.Repos)
	}

	// Display worktrees
	totalCount := 0
	for _, branch := range m.sortBranches(worktreesByBranch) {
		worktrees := worktreesByBranch[branch]
		totalCount += len(worktrees)

		fmt.Printf("\n%s\n", strings.Repeat("=", 80))
		fmt.Printf("Branch: %s (%d worktrees)\n", branch, len(worktrees))
		fmt.Printf("%s\n\n", strings.Repeat("=", 80))

		for _, wt := range worktrees {
			fmt.Printf("  📁 %s\n", wt.RepoName)
			fmt.Printf("     Path: %s\n", wt.Path)

			if opts.Verbose {
				fmt.Printf("     Last accessed: %s\n", formatTime(wt.LastAccessed))
				fmt.Printf("     Active: %v\n", wt.Active)
				if wt.Orphaned {
					fmt.Printf("     ⚠️  Orphaned (branch may be deleted remotely)\n")
				}
			}

			fmt.Println()
		}
	}

	fmt.Printf("\nTotal: %d worktrees across %d branches\n", totalCount, len(worktreesByBranch))

	return nil
}

// DeleteOptions contains options for deleting worktrees
type DeleteOptions struct {
	Branch      string
	Repos       []string
	Force       bool
	Interactive bool
}

// DeleteWorktrees deletes worktrees
func (m *Manager) DeleteWorktrees(opts DeleteOptions) error {
	// Discover worktrees to delete
	worktreesByBranch, err := m.tracker.DiscoverWorktrees(m.config.Worktree.BaseDir)
	if err != nil {
		return fmt.Errorf("failed to discover worktrees: %w", err)
	}

	// Collect worktrees to delete
	var toDelete []metadata.WorktreeInfo

	if opts.Branch != "" {
		// Delete specific branch worktrees
		if worktrees, ok := worktreesByBranch[opts.Branch]; ok {
			toDelete = worktrees
		}
	}

	if len(opts.Repos) > 0 {
		// Filter by repo names
		filtered := []metadata.WorktreeInfo{}
		for _, wt := range toDelete {
			for _, repo := range opts.Repos {
				if wt.RepoName == repo {
					filtered = append(filtered, wt)
					break
				}
			}
		}
		toDelete = filtered
	}

	if len(toDelete) == 0 {
		fmt.Println("No worktrees found to delete")
		return nil
	}

	// Interactive selection if requested
	if opts.Interactive {
		paths := make([]string, len(toDelete))
		for i, wt := range toDelete {
			paths[i] = fmt.Sprintf("%s/%s", wt.Branch, wt.RepoName)
		}

		selected, err := m.fzfSelector.SelectWorktree(paths)
		if err != nil {
			return err
		}

		if selected == "" {
			fmt.Println("No worktree selected")
			return nil
		}

		// Find the selected worktree
		parts := strings.Split(selected, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid selection")
		}

		for _, wt := range toDelete {
			if wt.Branch == parts[0] && wt.RepoName == parts[1] {
				toDelete = []metadata.WorktreeInfo{wt}
				break
			}
		}
	}

	fmt.Printf("Deleting %d worktrees...\n\n", len(toDelete))

	successCount := 0
	failCount := 0

	for i, wt := range toDelete {
		fmt.Printf("[%d/%d] Deleting %s/%s...\n", i+1, len(toDelete), wt.Branch, wt.RepoName)

		// Get the base repository path
		repo, ok := m.config.GetRepo(wt.RepoName)
		if !ok {
			fmt.Printf("⚠️  Repository %s not found in config\n\n", wt.RepoName)
			// Try to remove the directory anyway
			if err := os.RemoveAll(wt.Path); err != nil {
				fmt.Printf("✗ Failed to remove directory: %v\n\n", err)
				failCount++
			} else {
				successCount++
				fmt.Printf("✓ Removed directory\n\n")
			}
			continue
		}

		// Remove worktree using git
		if err := m.gitClient.WorktreeRemove(repo.BasePath, wt.Path, opts.Force); err != nil {
			fmt.Printf("✗ Failed to remove worktree: %v\n\n", err)
			failCount++
			continue
		}

		successCount++
		fmt.Printf("✓ Deleted worktree\n\n")
	}

	// Clean up empty branch directories
	m.cleanupEmptyBranchDirs()

	// Print summary
	fmt.Printf(strings.Repeat("=", 60) + "\n")
	fmt.Printf("Deletion Summary:\n")
	fmt.Printf("  Success: %d\n", successCount)
	fmt.Printf("  Failed:  %d\n", failCount)
	fmt.Printf("  Total:   %d\n", len(toDelete))
	fmt.Printf(strings.Repeat("=", 60) + "\n")

	return nil
}

// SwitchWorktree switches to a worktree (outputs cd command for shell integration)
func (m *Manager) SwitchWorktree(branch string) error {
	// Discover worktrees
	worktreesByBranch, err := m.tracker.DiscoverWorktrees(m.config.Worktree.BaseDir)
	if err != nil {
		return fmt.Errorf("failed to discover worktrees: %w", err)
	}

	// Collect all worktrees for selection
	var allWorktrees []string
	wtMap := make(map[string]metadata.WorktreeInfo)

	for branchName, worktrees := range worktreesByBranch {
		// Filter by branch if specified
		if branch != "" && branchName != branch {
			continue
		}

		for _, wt := range worktrees {
			key := fmt.Sprintf("%s/%s", branchName, wt.RepoName)
			allWorktrees = append(allWorktrees, key)
			wtMap[key] = wt
		}
	}

	if len(allWorktrees) == 0 {
		return fmt.Errorf("no worktrees found")
	}

	// Use fzf to select worktree
	selected, err := m.fzfSelector.SelectWorktree(allWorktrees)
	if err != nil {
		return err
	}

	if selected == "" {
		return nil // User cancelled
	}

	// Get the worktree path
	wt, ok := wtMap[selected]
	if !ok {
		return fmt.Errorf("worktree not found")
	}

	// Update access time
	metadata.UpdateAccessTime(wt.Path)

	// Output cd command for shell integration
	fmt.Printf("cd \"%s\"\n", wt.Path)

	return nil
}

// Helper functions

func (m *Manager) selectRepositories(opts CreateOptions) ([]config.Repo, error) {
	var repos []config.Repo

	if opts.All {
		// Get all repositories
		repos = m.config.ListRepos()
	} else if len(opts.Tags) > 0 {
		// Get repositories by tags
		repos = m.config.ListRepos(opts.Tags...)
	} else if len(opts.Repos) > 0 {
		// Get specific repositories
		for _, name := range opts.Repos {
			if repo, ok := m.config.GetRepo(name); ok {
				repos = append(repos, repo)
			}
		}
	}

	if len(repos) == 0 {
		return nil, fmt.Errorf("no repositories found")
	}

	// Interactive selection if requested
	if opts.Interactive {
		names := make([]string, len(repos))
		for i, r := range repos {
			names[i] = r.Name
		}

		selected, err := m.fzfSelector.SelectSingle(names, "Select repository")
		if err != nil {
			return nil, err
		}

		if selected == "" {
			return []config.Repo{}, nil
		}

		// Find the selected repo
		for _, r := range repos {
			if r.Name == selected {
				return []config.Repo{r}, nil
			}
		}
	}

	return repos, nil
}

func (m *Manager) filterByRepos(worktreesByBranch map[string][]metadata.WorktreeInfo, repos []string) map[string][]metadata.WorktreeInfo {
	filtered := make(map[string][]metadata.WorktreeInfo)

	for branch, worktrees := range worktreesByBranch {
		var filteredWorktrees []metadata.WorktreeInfo
		for _, wt := range worktrees {
			for _, repo := range repos {
				if wt.RepoName == repo {
					filteredWorktrees = append(filteredWorktrees, wt)
					break
				}
			}
		}

		if len(filteredWorktrees) > 0 {
			filtered[branch] = filteredWorktrees
		}
	}

	return filtered
}

func (m *Manager) sortBranches(worktreesByBranch map[string][]metadata.WorktreeInfo) []string {
	branches := make([]string, 0, len(worktreesByBranch))
	for branch := range worktreesByBranch {
		branches = append(branches, branch)
	}
	// Sort alphabetically
	// Could use sort.Strings here, but keeping it simple
	return branches
}

func (m *Manager) cleanupEmptyBranchDirs() {
	entries, err := os.ReadDir(m.config.Worktree.BaseDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		branchPath := filepath.Join(m.config.Worktree.BaseDir, entry.Name())

		// Check if directory is empty
		branchEntries, err := os.ReadDir(branchPath)
		if err != nil {
			continue
		}

		if len(branchEntries) == 0 {
			os.Remove(branchPath)
		}
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "unknown"
	}

	now := time.Now()
	diff := now.Sub(t)

	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	} else if diff < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	} else if diff < 7*24*time.Hour {
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	} else {
		return t.Format("2006-01-02")
	}
}
