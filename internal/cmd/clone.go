package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/asocpro/repo-manager/internal/provider"
	"github.com/asocpro/repo-manager/internal/repository"
	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a group or individual repository",
	Long: `Clone repositories from GitLab, GitHub, Gitea, or generic Git providers.

You can clone an entire group/organization or individual repositories.
Supports filtering by name patterns, metadata, and interactive selection.

Examples:
  # Clone a single repository
  repo-manager clone --provider gitlab --repo mygroup/myrepo

  # Clone an entire group
  repo-manager clone --provider github --group myorg

  # Clone with interactive selection
  repo-manager clone --provider gitlab --group mygroup --interactive

  # Clone with filters
  repo-manager clone --provider github --group myorg --filter "backend-*" --min-stars 10`,
	RunE: runClone,
}

var (
	cloneProvider     string
	cloneGroup        string
	cloneRepo         string
	cloneFilter       string
	cloneInteractive  bool
	cloneSSH          bool
	cloneTags         []string
	cloneFlatten      bool
	cloneMinStars     int
	cloneMinActivity  string
	cloneArchived     bool
)

func init() {
	rootCmd.AddCommand(cloneCmd)

	cloneCmd.Flags().StringVarP(&cloneProvider, "provider", "p", "", "Git provider (gitlab, github, gitea, generic) [required]")
	cloneCmd.Flags().StringVarP(&cloneGroup, "group", "g", "", "Group/organization to clone")
	cloneCmd.Flags().StringVarP(&cloneRepo, "repo", "r", "", "Individual repository to clone")
	cloneCmd.Flags().StringVarP(&cloneFilter, "filter", "f", "", "Filter pattern for repo names (glob)")
	cloneCmd.Flags().BoolVarP(&cloneInteractive, "interactive", "i", false, "Interactively select repos to clone")
	cloneCmd.Flags().BoolVar(&cloneSSH, "ssh", false, "Use SSH instead of HTTPS for cloning")
	cloneCmd.Flags().StringSliceVarP(&cloneTags, "tags", "t", []string{}, "Tags to add to cloned repositories")
	cloneCmd.Flags().BoolVar(&cloneFlatten, "flatten", false, "Flatten subgroups (ignore hierarchy)")
	cloneCmd.Flags().IntVar(&cloneMinStars, "min-stars", 0, "Minimum number of stars")
	cloneCmd.Flags().StringVar(&cloneMinActivity, "min-activity", "", "Minimum activity date (YYYY-MM-DD)")
	cloneCmd.Flags().BoolVar(&cloneArchived, "archived", false, "Include archived repositories")

	cloneCmd.MarkFlagRequired("provider")
}

func runClone(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Validate flags
	if cloneGroup == "" && cloneRepo == "" {
		return fmt.Errorf("either --group or --repo must be specified")
	}

	if cloneGroup != "" && cloneRepo != "" {
		return fmt.Errorf("cannot specify both --group and --repo")
	}

	// Create repository manager
	repoManager := repository.NewManager(cfg)

	// Clone individual repository
	if cloneRepo != "" {
		opts := repository.CloneRepoOptions{
			Provider: cloneProvider,
			RepoPath: cloneRepo,
			UseSSH:   cloneSSH,
			Tags:     cloneTags,
		}

		return repoManager.CloneRepo(ctx, opts)
	}

	// Clone group
	if cloneGroup != "" {
		// Parse min activity date if provided
		var minActivity time.Time
		if cloneMinActivity != "" {
			var err error
			minActivity, err = time.Parse("2006-01-02", cloneMinActivity)
			if err != nil {
				return fmt.Errorf("invalid min-activity date format (expected YYYY-MM-DD): %w", err)
			}
		}

		// Use config default for flatten if not specified
		flattenSubgroups := cloneFlatten
		if !cmd.Flags().Changed("flatten") {
			flattenSubgroups = cfg.Defaults.FlattenSubgroups
		}

		opts := repository.CloneGroupOptions{
			Provider:         cloneProvider,
			GroupPath:        cloneGroup,
			UseSSH:           cloneSSH,
			Tags:             cloneTags,
			FlattenSubgroups: flattenSubgroups,
			Interactive:      cloneInteractive,
			FilterOptions: provider.ListOptions{
				NameFilter:  cloneFilter,
				Archived:    cloneArchived,
				MinStars:    cloneMinStars,
				MinActivity: minActivity,
			},
		}

		return repoManager.CloneGroup(ctx, opts)
	}

	return nil
}
