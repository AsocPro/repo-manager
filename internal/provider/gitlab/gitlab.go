package gitlab

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/asocpro/repo-manager/internal/provider"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// Provider implements the provider.Provider interface for GitLab
type Provider struct {
	client *gitlab.Client
	config provider.AuthConfig
}

// New creates a new GitLab provider
func New(config provider.AuthConfig) (*Provider, error) {
	if config.URL == "" {
		config.URL = "https://gitlab.com"
	}

	var client *gitlab.Client
	var err error

	if config.Token != "" {
		client, err = gitlab.NewClient(config.Token, gitlab.WithBaseURL(config.URL))
	} else {
		client, err = gitlab.NewClient("", gitlab.WithBaseURL(config.URL))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}

	return &Provider{
		client: client,
		config: config,
	}, nil
}

// GetType returns the provider type
func (p *Provider) GetType() string {
	return "gitlab"
}

// Authenticate authenticates with GitLab
func (p *Provider) Authenticate(ctx context.Context) error {
	// Test authentication by getting current user
	_, _, err := p.client.Users.CurrentUser()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	return nil
}

// ListGroups lists all accessible groups
func (p *Provider) ListGroups(ctx context.Context) ([]provider.Group, error) {
	var allGroups []provider.Group

	opt := &gitlab.ListGroupsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	}

	for {
		groups, resp, err := p.client.Groups.ListGroups(opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list groups: %w", err)
		}

		for _, g := range groups {
			group := provider.Group{
				ID:          fmt.Sprintf("%d", g.ID),
				Name:        g.Name,
				Path:        g.Path,
				FullPath:    g.FullPath,
				Description: g.Description,
				WebURL:      g.WebURL,
			}

			if g.ParentID != 0 {
				group.ParentID = fmt.Sprintf("%d", g.ParentID)
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

// ListGroupRepos lists all repositories in a group
func (p *Provider) ListGroupRepos(ctx context.Context, groupPath string, opts provider.ListOptions) ([]provider.Repository, error) {
	var allRepos []provider.Repository

	listOpts := &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
		Archived:         &opts.Archived,
		IncludeSubGroups: gitlab.Ptr(true),
	}

	if opts.PerPage > 0 {
		listOpts.PerPage = opts.PerPage
	}
	if opts.Page > 0 {
		listOpts.Page = opts.Page
	}

	for {
		projects, resp, err := p.client.Groups.ListGroupProjects(groupPath, listOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to list group projects: %w", err)
		}

		for _, proj := range projects {
			// Apply filters
			if opts.NameFilter != "" {
				matched, _ := filepath.Match(opts.NameFilter, proj.Name)
				if !matched {
					continue
				}
			}

			if opts.MinStars > 0 && proj.StarCount < opts.MinStars {
				continue
			}

			if !opts.MinActivity.IsZero() && proj.LastActivityAt != nil {
				if proj.LastActivityAt.Before(opts.MinActivity) {
					continue
				}
			}

			repo := provider.Repository{
				ID:            fmt.Sprintf("%d", proj.ID),
				Name:          proj.Name,
				Path:          proj.Path,
				FullPath:      proj.PathWithNamespace,
				Description:   proj.Description,
				DefaultBranch: proj.DefaultBranch,
				SSHURL:        proj.SSHURLToRepo,
				HTTPURL:       proj.HTTPURLToRepo,
				WebURL:        proj.WebURL,
				Archived:      proj.Archived,
				Stars:         proj.StarCount,
				Forks:         proj.ForksCount,
			}

			if proj.LastActivityAt != nil {
				repo.LastActivity = *proj.LastActivityAt
			}

			// Extract group path from namespace
			if proj.Namespace != nil {
				repo.GroupPath = proj.Namespace.FullPath
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
	proj, _, err := p.client.Projects.GetProject(repoPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	repo := &provider.Repository{
		ID:            fmt.Sprintf("%d", proj.ID),
		Name:          proj.Name,
		Path:          proj.Path,
		FullPath:      proj.PathWithNamespace,
		Description:   proj.Description,
		DefaultBranch: proj.DefaultBranch,
		SSHURL:        proj.SSHURLToRepo,
		HTTPURL:       proj.HTTPURLToRepo,
		WebURL:        proj.WebURL,
		Archived:      proj.Archived,
		Stars:         proj.StarCount,
		Forks:         proj.ForksCount,
	}

	if proj.LastActivityAt != nil {
		repo.LastActivity = *proj.LastActivityAt
	}

	if proj.Namespace != nil {
		repo.GroupPath = proj.Namespace.FullPath
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

// Helper function to check if a string matches a glob pattern
func matchGlob(pattern, str string) bool {
	// Simple glob matching - can be enhanced
	if pattern == "" || pattern == "*" {
		return true
	}

	// Convert to lowercase for case-insensitive matching
	pattern = strings.ToLower(pattern)
	str = strings.ToLower(str)

	// Simple contains check if pattern has wildcards
	if strings.Contains(pattern, "*") {
		pattern = strings.Trim(pattern, "*")
		return strings.Contains(str, pattern)
	}

	return strings.Contains(str, pattern)
}
