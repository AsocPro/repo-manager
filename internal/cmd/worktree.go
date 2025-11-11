package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var worktreeCmd = &cobra.Command{
	Use:   "worktree",
	Short: "Manage git worktrees",
	Long: `Create, list, and delete git worktrees organized by branch names.

Worktrees are organized in branch-based directories:
  ~/worktrees/feature-x/repo1/
  ~/worktrees/feature-x/repo2/`,
}

var createCmd = &cobra.Command{
	Use:   "create [branch]",
	Short: "Create a new worktree",
	Long: `Create a new worktree for one or more repositories.

If the branch doesn't exist, it will be created automatically.
Supports batch operations via multi-select or tag filtering.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("worktree create command not yet implemented")
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all worktrees",
	Long:  `List all worktrees with their status and branch information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("worktree list command not yet implemented")
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a worktree",
	Long:  `Delete one or more worktrees with safety checks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("worktree delete command not yet implemented")
	},
}

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch to a worktree",
	Long: `Interactively select and switch to a worktree.
Outputs a cd command for shell integration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("worktree switch command not yet implemented")
	},
}

var (
	createBranch  string
	createRepos   []string
	createTags    []string
	createAll     bool
)

func init() {
	rootCmd.AddCommand(worktreeCmd)

	worktreeCmd.AddCommand(createCmd)
	worktreeCmd.AddCommand(listCmd)
	worktreeCmd.AddCommand(deleteCmd)
	worktreeCmd.AddCommand(switchCmd)

	createCmd.Flags().StringVarP(&createBranch, "branch", "b", "", "Branch name for the worktree")
	createCmd.Flags().StringSliceVarP(&createRepos, "repos", "r", []string{}, "Specific repositories to create worktrees for")
	createCmd.Flags().StringSliceVarP(&createTags, "tags", "t", []string{}, "Create worktrees for repos with these tags")
	createCmd.Flags().BoolVarP(&createAll, "all", "a", false, "Create worktrees for all managed repos")
}
