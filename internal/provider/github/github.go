package github

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/asocpro/repo-manager/internal/provider"
	"github.com/google/go-github/v57/github"
)

// Provider implements the provider.Provider interface for GitHub
type Provider struct {
	client *github.Client
	config provider.AuthConfig
}

// New creates a new GitHub provider
func New(config provider.AuthConfig) (*Provider, error) {
	var client *github.Client

	if config.Token != "" {
		client = github.NewClient(nil).WithAuthToken(config.Token)
	} else {
		client = github.NewClient(nil)
	}

	// Support for GitHub Enterprise
	if config.URL != "" && config.URL != "https://github.com" {
		var err error
		client, err = client.WithEnterpriseURLs(config.URL, config.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to set enterprise URL: %w", err)
		}
	}

	return &Provider{
		client: client,
		config: config,
	}, nil
}

// GetType returns the provider type
func (p *Provider) GetType() string {
	return "github"
}

// Authenticate authenticates with GitHub
func (p *Provider) Authenticate(ctx context.Context) error {
	// Test authentication by getting current user
	_, _, err := p.client.Users.Get(ctx, "")
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	return nil
}

// ListGroups lists all accessible organizations
func (p *Provider) ListGroups(ctx context.Context) ([]provider.Group, error) {
	var allGroups []provider.Group

	opt := &github.ListOptions{
		PerPage: 100,
		Page:    1,
	}

	for {
		orgs, resp, err := p.client.Organizations.List(ctx, "", opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list organizations: %w", err)
		}

		for _, org := range orgs {
			// Get full org details
			fullOrg, _, err := p.client.Organizations.Get(ctx, *org.Login)
			if err != nil {
				// If we can't get details, use basic info
				group := provider.Group{
					ID:       fmt.Sprintf("%d", *org.ID),
					Name:     *org.Login,
					Path:     *org.Login,
					FullPath: *org.Login,
					WebURL:   *org.URL,
				}
				allGroups = append(allGroups, group)
				continue
			}

			group := provider.Group{
				ID:       fmt.Sprintf("%d", *fullOrg.ID),
				Name:     *fullOrg.Login,
				Path:     *fullOrg.Login,
				FullPath: *fullOrg.Login,
				WebURL:   *fullOrg.HTMLURL,
			}

			if fullOrg.Description != nil {
				group.Description = *fullOrg.Description
			}

			allGroups = append(allGroups, group)
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allGroups, nil
}

// ListGroupRepos lists all repositories in an organization
func (p *Provider) ListGroupRepos(ctx context.Context, groupPath string, opts provider.ListOptions) ([]provider.Repository, error) {
	var allRepos []provider.Repository

	listOpts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
		Type: "all",
	}

	if opts.PerPage > 0 {
		listOpts.PerPage = opts.PerPage
	}
	if opts.Page > 0 {
		listOpts.Page = opts.Page
	}

	for {
		repos, resp, err := p.client.Repositories.ListByOrg(ctx, groupPath, listOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to list organization repositories: %w", err)
		}

		for _, ghRepo := range repos {
			// Apply filters
			if opts.NameFilter != "" {
				matched, _ := filepath.Match(opts.NameFilter, *ghRepo.Name)
				if !matched {
					continue
				}
			}

			if opts.MinStars > 0 && *ghRepo.StargazersCount < opts.MinStars {
				continue
			}

			if !opts.MinActivity.IsZero() && ghRepo.UpdatedAt != nil {
				if ghRepo.UpdatedAt.Before(opts.MinActivity) {
					continue
				}
			}

			if !opts.Archived && ghRepo.Archived != nil && *ghRepo.Archived {
				continue
			}

			repo := provider.Repository{
				ID:       fmt.Sprintf("%d", *ghRepo.ID),
				Name:     *ghRepo.Name,
				Path:     *ghRepo.Name,
				FullPath: *ghRepo.FullName,
				WebURL:   *ghRepo.HTMLURL,
				Archived: ghRepo.Archived != nil && *ghRepo.Archived,
			}

			if ghRepo.Description != nil {
				repo.Description = *ghRepo.Description
			}

			if ghRepo.DefaultBranch != nil {
				repo.DefaultBranch = *ghRepo.DefaultBranch
			}

			if ghRepo.SSHURL != nil {
				repo.SSHURL = *ghRepo.SSHURL
			}

			if ghRepo.CloneURL != nil {
				repo.HTTPURL = *ghRepo.CloneURL
			}

			if ghRepo.StargazersCount != nil {
				repo.Stars = *ghRepo.StargazersCount
			}

			if ghRepo.ForksCount != nil {
				repo.Forks = *ghRepo.ForksCount
			}

			if ghRepo.UpdatedAt != nil {
				repo.LastActivity = ghRepo.UpdatedAt.Time
			}

			if ghRepo.Owner != nil && ghRepo.Owner.Login != nil {
				repo.GroupPath = *ghRepo.Owner.Login
			}

			allRepos = append(allRepos, repo)
		}

		// If pagination is specified, only get requested page
		if opts.Page > 0 {
			break
		}

		if resp.NextPage == 0 {
			break
		}
		listOpts.Page = resp.NextPage
	}

	return allRepos, nil
}

// GetRepo gets details for a specific repository
func (p *Provider) GetRepo(ctx context.Context, repoPath string) (*provider.Repository, error) {
	// Parse owner/repo from path
	parts := strings.Split(repoPath, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository path: %s (expected owner/repo)", repoPath)
	}

	owner, repoName := parts[0], parts[1]

	ghRepo, _, err := p.client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	repo := &provider.Repository{
		ID:       fmt.Sprintf("%d", *ghRepo.ID),
		Name:     *ghRepo.Name,
		Path:     *ghRepo.Name,
		FullPath: *ghRepo.FullName,
		WebURL:   *ghRepo.HTMLURL,
		Archived: ghRepo.Archived != nil && *ghRepo.Archived,
	}

	if ghRepo.Description != nil {
		repo.Description = *ghRepo.Description
	}

	if ghRepo.DefaultBranch != nil {
		repo.DefaultBranch = *ghRepo.DefaultBranch
	}

	if ghRepo.SSHURL != nil {
		repo.SSHURL = *ghRepo.SSHURL
	}

	if ghRepo.CloneURL != nil {
		repo.HTTPURL = *ghRepo.CloneURL
	}

	if ghRepo.StargazersCount != nil {
		repo.Stars = *ghRepo.StargazersCount
	}

	if ghRepo.ForksCount != nil {
		repo.Forks = *ghRepo.ForksCount
	}

	if ghRepo.UpdatedAt != nil {
		repo.LastActivity = ghRepo.UpdatedAt.Time
	}

	if ghRepo.Owner != nil && ghRepo.Owner.Login != nil {
		repo.GroupPath = *ghRepo.Owner.Login
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
