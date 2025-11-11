package provider

import (
	"context"
	"fmt"

	"github.com/asocpro/repo-manager/internal/config"
	"github.com/asocpro/repo-manager/internal/provider/generic"
	"github.com/asocpro/repo-manager/internal/provider/gitea"
	"github.com/asocpro/repo-manager/internal/provider/github"
	"github.com/asocpro/repo-manager/internal/provider/gitlab"
)

// Manager manages provider instances
type Manager struct {
	config    *config.Config
	providers map[string]Provider
}

// NewManager creates a new provider manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config:    cfg,
		providers: make(map[string]Provider),
	}
}

// GetProvider returns a provider instance by name
func (m *Manager) GetProvider(ctx context.Context, name string) (Provider, error) {
	// Check if already initialized
	if p, ok := m.providers[name]; ok {
		return p, nil
	}

	// Get provider config
	providerCfg, ok := m.config.Providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found in configuration", name)
	}

	if !providerCfg.Enabled {
		return nil, fmt.Errorf("provider %s is not enabled", name)
	}

	// Create provider based on type
	authConfig := AuthConfig{
		Token: providerCfg.Token,
		URL:   providerCfg.URL,
	}

	var provider Provider
	var err error

	switch providerCfg.Type {
	case "gitlab":
		provider, err = gitlab.New(authConfig)
	case "github":
		provider, err = github.New(authConfig)
	case "gitea":
		provider, err = gitea.New(authConfig)
	case "generic":
		provider = generic.New()
	default:
		return nil, fmt.Errorf("unknown provider type: %s", providerCfg.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create provider %s: %w", name, err)
	}

	// Authenticate if needed
	if providerCfg.Type != "generic" && providerCfg.Token != "" {
		if err := provider.Authenticate(ctx); err != nil {
			return nil, fmt.Errorf("failed to authenticate with %s: %w", name, err)
		}
	}

	// Cache the provider
	m.providers[name] = provider

	return provider, nil
}

// ListProviders returns all configured provider names
func (m *Manager) ListProviders() []string {
	var names []string
	for name := range m.config.Providers {
		names = append(names, name)
	}
	return names
}

// ListEnabledProviders returns all enabled provider names
func (m *Manager) ListEnabledProviders() []string {
	var names []string
	for name, cfg := range m.config.Providers {
		if cfg.Enabled {
			names = append(names, name)
		}
	}
	return names
}
