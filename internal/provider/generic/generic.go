package generic

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/asocpro/repo-manager/internal/provider"
)

// Provider implements the provider.Provider interface for generic Git URLs
// This provider doesn't use any API, just plain Git URLs
type Provider struct{}

// New creates a new Generic Git provider
func New() *Provider {
	return &Provider{}
}

// GetType returns the provider type
func (p *Provider) GetType() string {
	return "generic"
}

// Authenticate is a no-op for generic provider (uses git credentials)
func (p *Provider) Authenticate(ctx context.Context) error {
	// Generic provider relies on git credential helper
	return nil
}

// ListGroups is not supported for generic provider
func (p *Provider) ListGroups(ctx context.Context) ([]provider.Group, error) {
	return nil, fmt.Errorf("ListGroups not supported for generic Git provider")
}

// ListGroupRepos is not supported for generic provider
func (p *Provider) ListGroupRepos(ctx context.Context, groupPath string, opts provider.ListOptions) ([]provider.Repository, error) {
	return nil, fmt.Errorf("ListGroupRepos not supported for generic Git provider")
}

// GetRepo creates a minimal repository object from a Git URL
// This is used when manually adding a repository by URL
func (p *Provider) GetRepo(ctx context.Context, repoURL string) (*provider.Repository, error) {
	// Parse repository name and details from URL
	name, fullPath := parseGitURL(repoURL)

	repo := &provider.Repository{
		Name:     name,
		Path:     name,
		FullPath: fullPath,
		HTTPURL:  repoURL,
		SSHURL:   repoURL,
	}

	// Try to determine if it's SSH or HTTPS
	if strings.HasPrefix(repoURL, "git@") || strings.HasPrefix(repoURL, "ssh://") {
		repo.SSHURL = repoURL
		repo.HTTPURL = convertToHTTPS(repoURL)
	} else if strings.HasPrefix(repoURL, "http://") || strings.HasPrefix(repoURL, "https://") {
		repo.HTTPURL = repoURL
		repo.SSHURL = convertToSSH(repoURL)
	}

	return repo, nil
}

// GetCloneURL returns the clone URL (just returns the URL as-is)
func (p *Provider) GetCloneURL(repo *provider.Repository, useSSH bool) string {
	if useSSH && repo.SSHURL != "" {
		return repo.SSHURL
	}
	return repo.HTTPURL
}

// parseGitURL extracts repository name and path from a Git URL
func parseGitURL(url string) (name, fullPath string) {
	// Remove .git suffix if present
	url = strings.TrimSuffix(url, ".git")

	// Handle different URL formats
	// SSH: git@github.com:owner/repo
	// HTTPS: https://github.com/owner/repo
	// Git: git://github.com/owner/repo

	var path string

	if strings.HasPrefix(url, "git@") {
		// SSH format: git@host:path
		parts := strings.SplitN(url, ":", 2)
		if len(parts) == 2 {
			path = parts[1]
		}
	} else {
		// URL format
		parts := strings.SplitN(url, "://", 2)
		if len(parts) == 2 {
			// Remove host
			hostAndPath := parts[1]
			pathParts := strings.SplitN(hostAndPath, "/", 2)
			if len(pathParts) == 2 {
				path = pathParts[1]
			}
		}
	}

	if path == "" {
		// Fallback: just use the URL
		path = url
	}

	// Extract name from path (last component)
	name = filepath.Base(path)
	fullPath = path

	return name, fullPath
}

// convertToHTTPS attempts to convert an SSH URL to HTTPS
func convertToHTTPS(sshURL string) string {
	if !strings.HasPrefix(sshURL, "git@") {
		return sshURL
	}

	// git@github.com:owner/repo -> https://github.com/owner/repo
	url := strings.TrimPrefix(sshURL, "git@")
	url = strings.Replace(url, ":", "/", 1)
	url = strings.TrimSuffix(url, ".git")

	return "https://" + url
}

// convertToSSH attempts to convert an HTTPS URL to SSH
func convertToSSH(httpsURL string) string {
	if !strings.HasPrefix(httpsURL, "https://") && !strings.HasPrefix(httpsURL, "http://") {
		return httpsURL
	}

	// https://github.com/owner/repo -> git@github.com:owner/repo
	url := strings.TrimPrefix(httpsURL, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimSuffix(url, ".git")

	parts := strings.SplitN(url, "/", 2)
	if len(parts) != 2 {
		return httpsURL
	}

	return fmt.Sprintf("git@%s:%s", parts[0], parts[1])
}
