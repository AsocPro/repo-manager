package repository

import (
	"context"
	"fmt"

	"github.com/asocpro/repo-manager/internal/config"
	"github.com/asocpro/repo-manager/internal/provider"
	"github.com/asocpro/repo-manager/internal/provider/generic"
	"github.com/asocpro/repo-manager/internal/provider/gitea"
	"github.com/asocpro/repo-manager/internal/provider/github"
	"github.com/asocpro/repo-manager/internal/provider/gitlab"
)

// ProviderManager manages provider instances
type ProviderManager struct {
	config    *config.Config
	providers map[string]provider.Provider
}

// NewProviderManager creates a new provider manager
func NewProviderManager(cfg *config.Config) *ProviderManager {
	return &ProviderManager{
		config:    cfg,
		providers: make(map[string]provider.Provider),
	}
}

// GetProvider returns a provider instance by name
func (m *ProviderManager) GetProvider(ctx context.Context, name string) (provider.Provider, error) {
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
	authConfig := provider.AuthConfig{
		Token: providerCfg.Token,
		URL:   providerCfg.URL,
	}

	var prov provider.Provider
	var err error

	switch providerCfg.Type {
	case "gitlab":
		prov, err = gitlab.New(authConfig)
	case "github":
		prov, err = github.New(authConfig)
	case "gitea":
		prov, err = gitea.New(authConfig)
	case "generic":
		prov = generic.New()
	default:
		return nil, fmt.Errorf("unknown provider type: %s", providerCfg.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create provider %s: %w", name, err)
	}

	// Authenticate if needed
	if providerCfg.Type != "generic" && providerCfg.Token != "" {
		if err := prov.Authenticate(ctx); err != nil {
			return nil, fmt.Errorf("failed to authenticate with %s: %w", name, err)
		}
	}

	// Cache the provider
	m.providers[name] = prov

	return prov, nil
}

// ListProviders returns all configured provider names
func (m *ProviderManager) ListProviders() []string {
	var names []string
	for name := range m.config.Providers {
		names = append(names, name)
	}
	return names
}

// ListEnabledProviders returns all enabled provider names
func (m *ProviderManager) ListEnabledProviders() []string {
	var names []string
	for name, cfg := range m.config.Providers {
		if cfg.Enabled {
			names = append(names, name)
		}
	}
	return names
}
