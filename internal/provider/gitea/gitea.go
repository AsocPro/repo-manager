package gitea

import (
	"context"
	"fmt"
	"path/filepath"

	"code.gitea.io/sdk/gitea"
	"github.com/asocpro/repo-manager/internal/provider"
)

// Provider implements the provider.Provider interface for Gitea/Forgejo
type Provider struct {
	client *gitea.Client
	config provider.AuthConfig
}

// New creates a new Gitea/Forgejo provider
func New(config provider.AuthConfig) (*Provider, error) {
	if config.URL == "" {
		return nil, fmt.Errorf("Gitea/Forgejo requires a URL to be specified")
	}

	opts := []gitea.ClientOption{}
	if config.Token != "" {
		opts = append(opts, gitea.SetToken(config.Token))
	}

	client, err := gitea.NewClient(config.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gitea client: %w", err)
	}

	return &Provider{
		client: client,
		config: config,
	}, nil
}

// GetType returns the provider type
func (p *Provider) GetType() string {
	return "gitea"
}

// Authenticate authenticates with Gitea/Forgejo
func (p *Provider) Authenticate(ctx context.Context) error {
	// Test authentication by getting current user
	_, _, err := p.client.GetMyUserInfo()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	return nil
}

// ListGroups lists all accessible organizations
func (p *Provider) ListGroups(ctx context.Context) ([]provider.Group, error) {
	var allGroups []provider.Group

	// List organizations
	orgs, _, err := p.client.ListMyOrgs(gitea.ListOrgsOptions{
		ListOptions: gitea.ListOptions{
			Page:     1,
			PageSize: 100,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	for _, org := range orgs {
		group := provider.Group{
			ID:          fmt.Sprintf("%d", org.ID),
			Name:        org.UserName,
			Path:        org.UserName,
			FullPath:    org.UserName,
			Description: org.Description,
			WebURL:      org.Website,
		}

		allGroups = append(allGroups, group)
	}

	return allGroups, nil
}

// ListGroupRepos lists all repositories in an organization
func (p *Provider) ListGroupRepos(ctx context.Context, groupPath string, opts provider.ListOptions) ([]provider.Repository, error) {
	var allRepos []provider.Repository

	listOpts := gitea.ListOrgReposOptions{
		ListOptions: gitea.ListOptions{
			Page:     1,
			PageSize: 100,
		},
	}

	if opts.PerPage > 0 {
		listOpts.PageSize = opts.PerPage
	}
	if opts.Page > 0 {
		listOpts.Page = opts.Page
	}

	for {
		repos, _, err := p.client.ListOrgRepos(groupPath, listOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to list organization repositories: %w", err)
		}

		if len(repos) == 0 {
			break
		}

		for _, giteaRepo := range repos {
			// Apply filters
			if opts.NameFilter != "" {
				matched, _ := filepath.Match(opts.NameFilter, giteaRepo.Name)
				if !matched {
					continue
				}
			}

			if opts.MinStars > 0 && giteaRepo.Stars < opts.MinStars {
				continue
			}

			if !opts.MinActivity.IsZero() && giteaRepo.Updated != nil {
				if giteaRepo.Updated.Before(opts.MinActivity) {
					continue
				}
			}

			if !opts.Archived && giteaRepo.Archived {
				continue
			}

			repo := provider.Repository{
				ID:            fmt.Sprintf("%d", giteaRepo.ID),
				Name:          giteaRepo.Name,
				Path:          giteaRepo.Name,
				FullPath:      giteaRepo.FullName,
				Description:   giteaRepo.Description,
				DefaultBranch: giteaRepo.DefaultBranch,
				SSHURL:        giteaRepo.SSHURL,
				HTTPURL:       giteaRepo.CloneURL,
				WebURL:        giteaRepo.HTMLURL,
				Archived:      giteaRepo.Archived,
				Stars:         giteaRepo.Stars,
				Forks:         giteaRepo.Forks,
			}

			if giteaRepo.Updated != nil {
				repo.LastActivity = *giteaRepo.Updated
			}

			if giteaRepo.Owner != nil {
				repo.GroupPath = giteaRepo.Owner.UserName
			}

			allRepos = append(allRepos, repo)
		}

		// If pagination is specified, only get requested page
		if opts.Page > 0 {
			break
		}

		// Check if there are more pages
		if len(repos) < listOpts.PageSize {
			break
		}

		listOpts.Page++
	}

	return allRepos, nil
}

// GetRepo gets details for a specific repository
func (p *Provider) GetRepo(ctx context.Context, repoPath string) (*provider.Repository, error) {
	// Parse owner/repo from path
	owner, repoName, err := parseRepoPath(repoPath)
	if err != nil {
		return nil, err
	}

	giteaRepo, _, err := p.client.GetRepo(owner, repoName)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	repo := &provider.Repository{
		ID:            fmt.Sprintf("%d", giteaRepo.ID),
		Name:          giteaRepo.Name,
		Path:          giteaRepo.Name,
		FullPath:      giteaRepo.FullName,
		Description:   giteaRepo.Description,
		DefaultBranch: giteaRepo.DefaultBranch,
		SSHURL:        giteaRepo.SSHURL,
		HTTPURL:       giteaRepo.CloneURL,
		WebURL:        giteaRepo.HTMLURL,
		Archived:      giteaRepo.Archived,
		Stars:         giteaRepo.Stars,
		Forks:         giteaRepo.Forks,
	}

	if giteaRepo.Updated != nil {
		repo.LastActivity = *giteaRepo.Updated
	}

	if giteaRepo.Owner != nil {
		repo.GroupPath = giteaRepo.Owner.UserName
	}

	return repo, nil
}

// GetCloneURL returns the clone URL for a repository
func (p *Provider) GetCloneURL(repo *provider.Repository, useSSH bool) string {
	if useSSH {
		return repo.SSHURL
	}
	return repo.HTTPURL
}

// parseRepoPath parses owner/repo from a repository path
func parseRepoPath(repoPath string) (string, string, error) {
	parts := filepath.SplitList(repoPath)
	if len(parts) < 2 {
		// Try splitting by /
		parts = []string{}
		for _, p := range filepath.SplitList(repoPath) {
			parts = append(parts, p)
		}
	}

	// Simple split by /
	idx := 0
	for i, c := range repoPath {
		if c == '/' {
			idx = i
			break
		}
	}

	if idx == 0 {
		return "", "", fmt.Errorf("invalid repository path: %s (expected owner/repo)", repoPath)
	}

	owner := repoPath[:idx]
	repo := repoPath[idx+1:]

	return owner, repo, nil
}
