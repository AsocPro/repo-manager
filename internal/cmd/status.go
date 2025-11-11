package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of repositories and worktrees",
	Long: `Display status information for all managed repositories and worktrees.

Shows:
  - Repository status (dirty, ahead/behind)
  - Worktree status
  - Orphaned worktree warnings
  - Activity summary`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("status command not yet implemented")
	},
}

var (
	statusVerbose bool
	statusRepos   []string
	statusTags    []string
)

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().BoolVarP(&statusVerbose, "verbose", "v", false, "Show detailed status")
	statusCmd.Flags().StringSliceVarP(&statusRepos, "repos", "r", []string{}, "Show status for specific repositories")
	statusCmd.Flags().StringSliceVarP(&statusTags, "tags", "t", []string{}, "Show status for repos with these tags")
}
