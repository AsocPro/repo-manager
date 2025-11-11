package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	Version   string              `yaml:"version"`
	Providers map[string]Provider `yaml:"providers"`
	Worktree  WorktreeConfig      `yaml:"worktree"`
	Defaults  Defaults            `yaml:"defaults"`
	Repos     map[string]Repo     `yaml:"repos,omitempty"`
}

// Provider represents a Git provider configuration
type Provider struct {
	Type     string `yaml:"type"` // gitlab, github, gitea, generic
	URL      string `yaml:"url,omitempty"`
	Token    string `yaml:"token,omitempty"`
	BaseDir  string `yaml:"base_dir"`
	Enabled  bool   `yaml:"enabled"`
}

// WorktreeConfig represents worktree-specific configuration
type WorktreeConfig struct {
	BaseDir string `yaml:"base_dir"`
}

// Defaults represents default behavior settings
type Defaults struct {
	FlattenSubgroups bool   `yaml:"flatten_subgroups"`
	DefaultBranch    string `yaml:"default_branch"`
}

// Repo represents metadata for a tracked repository
type Repo struct {
	Name          string   `yaml:"name"`
	Provider      string   `yaml:"provider"`
	RemoteURL     string   `yaml:"remote_url"`
	DefaultBranch string   `yaml:"default_branch"`
	BasePath      string   `yaml:"base_path"`
	Tags          []string `yaml:"tags,omitempty"`
	Active        bool     `yaml:"active"`
	LastAccessed  string   `yaml:"last_accessed,omitempty"`
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "repo-manager")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(configDir, "config.yaml"), nil
}

// Load loads the configuration from the config file
func Load(cfgFile string) (*Config, error) {
	if cfgFile == "" {
		var err error
		cfgFile, err = GetConfigPath()
		if err != nil {
			return nil, err
		}
	}

	// If config doesn't exist, create a default one
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		config := Default()
		if err := config.Save(cfgFile); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return config, nil
	}

	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Save saves the configuration to the config file
func (c *Config) Save(cfgFile string) error {
	if cfgFile == "" {
		var err error
		cfgFile, err = GetConfigPath()
		if err != nil {
			return err
		}
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cfgFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Default returns a default configuration
func Default() *Config {
	homeDir, _ := os.UserHomeDir()

	return &Config{
		Version: "1",
		Providers: map[string]Provider{
			"gitlab": {
				Type:    "gitlab",
				URL:     "https://gitlab.com",
				BaseDir: filepath.Join(homeDir, "git", "bases", "gitlab"),
				Enabled: false,
			},
			"github": {
				Type:    "github",
				URL:     "https://github.com",
				BaseDir: filepath.Join(homeDir, "git", "bases", "github"),
				Enabled: false,
			},
			"gitea": {
				Type:    "gitea",
				BaseDir: filepath.Join(homeDir, "git", "bases", "gitea"),
				Enabled: false,
			},
		},
		Worktree: WorktreeConfig{
			BaseDir: filepath.Join(homeDir, "worktrees"),
		},
		Defaults: Defaults{
			FlattenSubgroups: false,
			DefaultBranch:    "main",
		},
		Repos: make(map[string]Repo),
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Worktree.BaseDir == "" {
		return fmt.Errorf("worktree base_dir cannot be empty")
	}

	for name, provider := range c.Providers {
		if provider.BaseDir == "" {
			return fmt.Errorf("provider %s: base_dir cannot be empty", name)
		}
		if provider.Type == "" {
			return fmt.Errorf("provider %s: type cannot be empty", name)
		}
	}

	return nil
}

// AddRepo adds a repository to the configuration
func (c *Config) AddRepo(name string, repo Repo) {
	if c.Repos == nil {
		c.Repos = make(map[string]Repo)
	}
	c.Repos[name] = repo
}

// GetRepo retrieves a repository from the configuration
func (c *Config) GetRepo(name string) (Repo, bool) {
	repo, ok := c.Repos[name]
	return repo, ok
}

// RemoveRepo removes a repository from the configuration
func (c *Config) RemoveRepo(name string) {
	delete(c.Repos, name)
}

// ListRepos returns all repositories, optionally filtered by tags
func (c *Config) ListRepos(tags ...string) []Repo {
	var repos []Repo

	for _, repo := range c.Repos {
		if len(tags) == 0 {
			repos = append(repos, repo)
			continue
		}

		// Check if repo has all required tags
		hasAllTags := true
		for _, requiredTag := range tags {
			found := false
			for _, repoTag := range repo.Tags {
				if repoTag == requiredTag {
					found = true
					break
				}
			}
			if !found {
				hasAllTags = false
				break
			}
		}

		if hasAllTags {
			repos = append(repos, repo)
		}
	}

	return repos
}
