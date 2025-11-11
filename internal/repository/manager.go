package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/asocpro/repo-manager/internal/config"
	"github.com/asocpro/repo-manager/internal/git"
	"github.com/asocpro/repo-manager/internal/provider"
	"github.com/asocpro/repo-manager/internal/ui"
)

// Manager handles repository operations
type Manager struct {
	config          *config.Config
	providerManager *ProviderManager
	gitClient       *git.Client
	fzfSelector     *ui.FzfSelector
}

// NewManager creates a new repository manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config:          cfg,
		providerManager: NewProviderManager(cfg),
		gitClient:       git.New(),
		fzfSelector:     ui.NewFzfSelector(),
	}
}

// CloneRepoOptions contains options for cloning a repository
type CloneRepoOptions struct {
	Provider      string
	RepoPath      string
	UseSSH        bool
	Tags          []string
	DefaultBranch string
}

// CloneRepo clones an individual repository
func (m *Manager) CloneRepo(ctx context.Context, opts CloneRepoOptions) error {
	// Get provider
	prov, err := m.providerManager.GetProvider(ctx, opts.Provider)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	// Get repository details from provider
	repo, err := prov.GetRepo(ctx, opts.RepoPath)
	if err != nil {
		return fmt.Errorf("failed to get repository details: %w", err)
	}

	// Determine clone URL
	cloneURL := prov.GetCloneURL(repo, opts.UseSSH)

	// Get provider config to determine base directory
	providerCfg, ok := m.config.Providers[opts.Provider]
	if !ok {
		return fmt.Errorf("provider %s not found in configuration", opts.Provider)
	}

	// Determine destination path
	destPath := filepath.Join(providerCfg.BaseDir, repo.Name)

	// Check if already exists
	if _, err := os.Stat(destPath); err == nil {
		return fmt.Errorf("repository already exists at %s", destPath)
	}

	// Ensure base directory exists
	if err := os.MkdirAll(providerCfg.BaseDir, 0755); err != nil {
		return fmt.Errorf("failed to create base directory: %w", err)
	}

	// Clone the repository
	fmt.Printf("Cloning %s to %s...\n", repo.FullPath, destPath)
	if err := m.gitClient.Clone(cloneURL, destPath); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Determine default branch
	defaultBranch := opts.DefaultBranch
	if defaultBranch == "" {
		if repo.DefaultBranch != "" {
			defaultBranch = repo.DefaultBranch
		} else {
			// Try to detect from cloned repo
			detectedBranch, err := m.gitClient.GetDefaultBranch(destPath)
			if err == nil {
				defaultBranch = detectedBranch
			} else {
				defaultBranch = m.config.Defaults.DefaultBranch
			}
		}
	}

	// Add repository to configuration
	repoConfig := config.Repo{
		Name:          repo.Name,
		Provider:      opts.Provider,
		RemoteURL:     cloneURL,
		DefaultBranch: defaultBranch,
		BasePath:      destPath,
		Tags:          opts.Tags,
		Active:        true,
	}

	m.config.AddRepo(repo.Name, repoConfig)

	// Save configuration
	if err := m.config.Save(""); err != nil {
		return fmt.Errorf("warning: failed to save configuration: %w", err)
	}

	fmt.Printf("✓ Successfully cloned %s\n", repo.FullPath)
	return nil
}

// CloneGroupOptions contains options for cloning a group
type CloneGroupOptions struct {
	Provider         string
	GroupPath        string
	UseSSH           bool
	Tags             []string
	FlattenSubgroups bool
	FilterOptions    provider.ListOptions
	Interactive      bool
}

// CloneGroup clones all repositories in a group/organization
func (m *Manager) CloneGroup(ctx context.Context, opts CloneGroupOptions) error {
	// Get provider
	prov, err := m.providerManager.GetProvider(ctx, opts.Provider)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	// List repositories in the group
	fmt.Printf("Fetching repositories from %s/%s...\n", opts.Provider, opts.GroupPath)
	repos, err := prov.ListGroupRepos(ctx, opts.GroupPath, opts.FilterOptions)
	if err != nil {
		return fmt.Errorf("failed to list group repositories: %w", err)
	}

	if len(repos) == 0 {
		fmt.Println("No repositories found")
		return nil
	}

	fmt.Printf("Found %d repositories\n", len(repos))

	// If interactive, let user select which repos to clone
	selectedRepos := repos
	if opts.Interactive {
		selectedRepos, err = m.selectRepositories(repos)
		if err != nil {
			return fmt.Errorf("failed to select repositories: %w", err)
		}
	}

	if len(selectedRepos) == 0 {
		fmt.Println("No repositories selected")
		return nil
	}

	// Clone each selected repository
	successCount := 0
	failCount := 0

	for i, repo := range selectedRepos {
		fmt.Printf("\n[%d/%d] Cloning %s...\n", i+1, len(selectedRepos), repo.FullPath)

		// Determine destination path
		destPath := m.getRepoDestPath(opts.Provider, &repo, opts.FlattenSubgroups)

		// Check if already exists
		if _, err := os.Stat(destPath); err == nil {
			fmt.Printf("⊗ Already exists: %s\n", destPath)
			continue
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			fmt.Printf("✗ Failed to create directory: %v\n", err)
			failCount++
			continue
		}

		// Clone the repository
		cloneURL := prov.GetCloneURL(&repo, opts.UseSSH)
		if err := m.gitClient.Clone(cloneURL, destPath); err != nil {
			fmt.Printf("✗ Failed to clone: %v\n", err)
			failCount++
			continue
		}

		// Determine default branch
		defaultBranch := repo.DefaultBranch
		if defaultBranch == "" {
			detectedBranch, err := m.gitClient.GetDefaultBranch(destPath)
			if err == nil {
				defaultBranch = detectedBranch
			} else {
				defaultBranch = m.config.Defaults.DefaultBranch
			}
		}

		// Add repository to configuration
		repoConfig := config.Repo{
			Name:          repo.Name,
			Provider:      opts.Provider,
			RemoteURL:     cloneURL,
			DefaultBranch: defaultBranch,
			BasePath:      destPath,
			Tags:          opts.Tags,
			Active:        true,
		}

		m.config.AddRepo(repo.Name, repoConfig)
		successCount++

		fmt.Printf("✓ Successfully cloned to %s\n", destPath)
	}

	// Save configuration
	if err := m.config.Save(""); err != nil {
		fmt.Printf("Warning: failed to save configuration: %v\n", err)
	}

	// Print summary
	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("Clone Summary:\n")
	fmt.Printf("  Success: %d\n", successCount)
	fmt.Printf("  Failed:  %d\n", failCount)
	fmt.Printf("  Total:   %d\n", len(selectedRepos))
	fmt.Printf(strings.Repeat("=", 60) + "\n")

	return nil
}

// getRepoDestPath determines the destination path for a repository
func (m *Manager) getRepoDestPath(providerName string, repo *provider.Repository, flattenSubgroups bool) string {
	providerCfg := m.config.Providers[providerName]
	baseDir := providerCfg.BaseDir

	if flattenSubgroups || repo.GroupPath == "" {
		// Flat structure: baseDir/repoName
		return filepath.Join(baseDir, repo.Name)
	}

	// Preserve hierarchy: baseDir/group/subgroup/repoName
	// Remove the repo name from FullPath to get the group path
	groupPath := strings.TrimSuffix(repo.FullPath, "/"+repo.Name)
	if groupPath == repo.Name {
		// No group path, just repo name
		return filepath.Join(baseDir, repo.Name)
	}

	return filepath.Join(baseDir, groupPath, repo.Name)
}

// selectRepositories allows interactive repository selection using fzf
func (m *Manager) selectRepositories(repos []provider.Repository) ([]provider.Repository, error) {
	// Check if fzf is available
	if !m.fzfSelector.CheckFzfInstalled() {
		fmt.Println("Warning: fzf not installed, cloning all repositories")
		fmt.Println("Install fzf for interactive selection: sudo dnf install fzf")
		return repos, nil
	}

	// Use fzf to select repositories
	selected, err := m.fzfSelector.SelectRepositories(repos)
	if err != nil {
		return nil, err
	}

	if selected == nil {
		// User cancelled
		return []provider.Repository{}, nil
	}

	return selected, nil
}
