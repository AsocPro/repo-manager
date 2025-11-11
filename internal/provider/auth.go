package provider

import (
	"fmt"
	"os/exec"
	"strings"
)

// AuthHelper provides utilities for the hybrid authentication system
type AuthHelper struct{}

// NewAuthHelper creates a new authentication helper
func NewAuthHelper() *AuthHelper {
	return &AuthHelper{}
}

// CheckGitCredentialHelper checks if git credential helper is configured
func (h *AuthHelper) CheckGitCredentialHelper() (bool, string, error) {
	cmd := exec.Command("git", "config", "--global", "credential.helper")
	output, err := cmd.Output()
	if err != nil {
		return false, "", nil // Not an error, just not configured
	}

	helper := strings.TrimSpace(string(output))
	return helper != "", helper, nil
}

// ValidateToken performs basic validation on API tokens
func (h *AuthHelper) ValidateToken(token string) error {
	if token == "" {
		return fmt.Errorf("token is empty")
	}

	// Basic length check (most tokens are at least 20 chars)
	if len(token) < 20 {
		return fmt.Errorf("token appears to be too short (expected at least 20 characters)")
	}

	return nil
}

// GetAuthStatus returns a summary of authentication status
func (h *AuthHelper) GetAuthStatus() AuthStatus {
	status := AuthStatus{
		GitCredentialHelper: false,
	}

	// Check git credential helper
	configured, helper, _ := h.CheckGitCredentialHelper()
	status.GitCredentialHelper = configured
	status.CredentialHelperType = helper

	return status
}

// AuthStatus represents the authentication status
type AuthStatus struct {
	GitCredentialHelper   bool
	CredentialHelperType  string
}

// String returns a human-readable representation of auth status
func (s AuthStatus) String() string {
	var sb strings.Builder

	sb.WriteString("Authentication Status:\n")
	sb.WriteString("=====================\n\n")

	// Git Credential Helper
	sb.WriteString("Git Credential Helper: ")
	if s.GitCredentialHelper {
		sb.WriteString(fmt.Sprintf("✓ Configured (%s)\n", s.CredentialHelperType))
	} else {
		sb.WriteString("✗ Not configured\n")
		sb.WriteString("  Run: git config --global credential.helper store\n")
		sb.WriteString("  Or:  git config --global credential.helper cache\n")
	}

	sb.WriteString("\nHybrid Authentication Model:\n")
	sb.WriteString("- API Tokens: Used for provider APIs (listing groups, repos, metadata)\n")
	sb.WriteString("- Git Credentials: Used for cloning repositories (HTTPS/SSH)\n")
	sb.WriteString("\nTo configure API tokens, edit ~/.config/repo-manager/config.yaml\n")

	return sb.String()
}

// AuthenticationGuide provides setup instructions
const AuthenticationGuide = `
Hybrid Authentication System
=============================

This tool uses a hybrid authentication approach:

1. API Tokens (for provider APIs)
   - Used to list groups, repositories, and fetch metadata
   - Configured in ~/.config/repo-manager/config.yaml
   - Example:
     providers:
       gitlab:
         token: "glpat-xxxxxxxxxxxx"
       github:
         token: "ghp_xxxxxxxxxxxx"

2. Git Credentials (for cloning)
   - Used when actually cloning repositories
   - Managed by git's credential helper system
   - Recommended configuration:

     For HTTPS (recommended):
       git config --global credential.helper store
       # Or for temporary caching:
       git config --global credential.helper "cache --timeout=3600"

     For SSH:
       - Set up SSH keys: ssh-keygen -t ed25519 -C "your_email@example.com"
       - Add to ssh-agent: ssh-add ~/.ssh/id_ed25519
       - Add public key to your provider (GitHub/GitLab/Gitea)

Getting API Tokens:
-------------------

GitLab:
  1. Go to https://gitlab.com/-/profile/personal_access_tokens
  2. Create token with 'read_api' and 'read_repository' scopes
  3. Add to config: token: "glpat-xxxxxxxxxxxx"

GitHub:
  1. Go to https://github.com/settings/tokens
  2. Generate new token (classic) with 'repo' scope
  3. Add to config: token: "ghp_xxxxxxxxxxxx"

Gitea/Forgejo:
  1. Go to your instance: https://your-gitea.com/user/settings/applications
  2. Generate new token
  3. Add to config: token: "xxxxxxxxxxxx"

Generic Git:
  - No API token needed
  - Only uses git credential helper for cloning
`
