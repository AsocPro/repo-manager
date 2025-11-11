package provider

import (
	"context"
	"time"
)

// Provider defines the interface for all Git providers
type Provider interface {
	// GetType returns the provider type (gitlab, github, gitea, generic)
	GetType() string

	// Authenticate authenticates with the provider using the given credentials
	Authenticate(ctx context.Context) error

	// ListGroups lists all accessible groups/organizations
	ListGroups(ctx context.Context) ([]Group, error)

	// ListGroupRepos lists all repositories in a group/organization
	ListGroupRepos(ctx context.Context, groupPath string, opts ListOptions) ([]Repository, error)

	// GetRepo gets details for a specific repository
	GetRepo(ctx context.Context, repoPath string) (*Repository, error)

	// GetCloneURL returns the clone URL for a repository (supports SSH/HTTPS)
	GetCloneURL(repo *Repository, useSSH bool) string
}

// Group represents a group/organization
type Group struct {
	ID          string
	Name        string
	Path        string
	FullPath    string
	Description string
	ParentID    string
	WebURL      string
}

// Repository represents a git repository
type Repository struct {
	ID            string
	Name          string
	Path          string
	FullPath      string
	Description   string
	DefaultBranch string
	SSHURL        string
	HTTPURL       string
	WebURL        string
	Archived      bool
	LastActivity  time.Time
	Stars         int
	Forks         int
	GroupPath     string
}

// ListOptions contains options for listing repositories
type ListOptions struct {
	// NameFilter filters repositories by name pattern (supports glob)
	NameFilter string

	// Archived includes archived repositories if true
	Archived bool

	// MinStars filters repositories with at least this many stars
	MinStars int

	// MinActivity filters repositories active since this time
	MinActivity time.Time

	// Page and PerPage for pagination
	Page    int
	PerPage int
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	Token    string
	URL      string
	Username string
	Password string
}

// ProviderFactory creates provider instances
type ProviderFactory interface {
	Create(authConfig AuthConfig) (Provider, error)
}
