package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync repositories and worktrees",
	Long: `Fetch updates from remote repositories.

Can sync base repositories, worktrees, or both.
Supports batch operations across multiple repos.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("sync command not yet implemented")
	},
}

var (
	syncBases     bool
	syncWorktrees bool
	syncRepos     []string
	syncTags      []string
	syncAll       bool
)

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().BoolVar(&syncBases, "bases", true, "Sync base repositories")
	syncCmd.Flags().BoolVar(&syncWorktrees, "worktrees", false, "Sync worktrees (pull)")
	syncCmd.Flags().StringSliceVarP(&syncRepos, "repos", "r", []string{}, "Specific repositories to sync")
	syncCmd.Flags().StringSliceVarP(&syncTags, "tags", "t", []string{}, "Sync repos with these tags")
	syncCmd.Flags().BoolVarP(&syncAll, "all", "a", false, "Sync all managed repos")
}
