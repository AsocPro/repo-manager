package cmd

import (
	"fmt"

	"github.com/asocpro/repo-manager/internal/worktree"
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
	Use:   "create",
	Short: "Create a new worktree",
	Long: `Create a new worktree for one or more repositories.

If the branch doesn't exist, it will be created automatically.
Supports batch operations via multi-select or tag filtering.

Examples:
  # Create worktree for all repos on branch feature-x
  repo-manager worktree create --branch feature-x --all

  # Create worktree for specific repos
  repo-manager worktree create --branch bugfix-123 --repos repo1,repo2

  # Create with interactive selection
  repo-manager worktree create --branch feature-y --interactive

  # Create for repos with specific tags
  repo-manager worktree create --branch dev --tags frontend,active`,
	RunE: runWorktreeCreate,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all worktrees",
	Long: `List all worktrees with their status and branch information.

Examples:
  # List all worktrees
  repo-manager worktree list

  # List worktrees for a specific branch
  repo-manager worktree list --branch feature-x

  # List with verbose details
  repo-manager worktree list --verbose`,
	RunE: runWorktreeList,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a worktree",
	Long: `Delete one or more worktrees with safety checks.

Examples:
  # Delete all worktrees for a branch
  repo-manager worktree delete --branch feature-x

  # Delete with interactive selection
  repo-manager worktree delete --interactive

  # Force delete (skip safety checks)
  repo-manager worktree delete --branch old-feature --force`,
	RunE: runWorktreeDelete,
}

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch to a worktree",
	Long: `Interactively select and switch to a worktree.
Outputs a cd command for shell integration.

To use this command with your shell, add this function to your ~/.bashrc or ~/.zshrc:

  wt() {
    eval $(repo-manager worktree switch "$@")
  }

Then use: wt
Or:       wt --branch feature-x`,
	RunE: runWorktreeSwitch,
}

var (
	createBranch      string
	createRepos       []string
	createTags        []string
	createAll         bool
	createInteractive bool

	listBranch  string
	listRepos   []string
	listVerbose bool

	deleteBranch      string
	deleteRepos       []string
	deleteForce       bool
	deleteInteractive bool

	switchBranch string
)

func init() {
	rootCmd.AddCommand(worktreeCmd)

	worktreeCmd.AddCommand(createCmd)
	worktreeCmd.AddCommand(listCmd)
	worktreeCmd.AddCommand(deleteCmd)
	worktreeCmd.AddCommand(switchCmd)

	// Create flags
	createCmd.Flags().StringVarP(&createBranch, "branch", "b", "", "Branch name for the worktree [required]")
	createCmd.Flags().StringSliceVarP(&createRepos, "repos", "r", []string{}, "Specific repositories to create worktrees for")
	createCmd.Flags().StringSliceVarP(&createTags, "tags", "t", []string{}, "Create worktrees for repos with these tags")
	createCmd.Flags().BoolVarP(&createAll, "all", "a", false, "Create worktrees for all managed repos")
	createCmd.Flags().BoolVarP(&createInteractive, "interactive", "i", false, "Interactively select repository")
	createCmd.MarkFlagRequired("branch")

	// List flags
	listCmd.Flags().StringVarP(&listBranch, "branch", "b", "", "Filter by branch name")
	listCmd.Flags().StringSliceVarP(&listRepos, "repos", "r", []string{}, "Filter by repository names")
	listCmd.Flags().BoolVarP(&listVerbose, "verbose", "v", false, "Show detailed information")

	// Delete flags
	deleteCmd.Flags().StringVarP(&deleteBranch, "branch", "b", "", "Delete worktrees for this branch")
	deleteCmd.Flags().StringSliceVarP(&deleteRepos, "repos", "r", []string{}, "Delete worktrees for these repositories")
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "Force delete (skip safety checks)")
	deleteCmd.Flags().BoolVarP(&deleteInteractive, "interactive", "i", false, "Interactively select worktree to delete")

	// Switch flags
	switchCmd.Flags().StringVarP(&switchBranch, "branch", "b", "", "Filter by branch name")
}

func runWorktreeCreate(cmd *cobra.Command, args []string) error {
	// Validate that at least one selection method is provided
	if !createAll && len(createRepos) == 0 && len(createTags) == 0 && !createInteractive {
		return fmt.Errorf("must specify one of: --all, --repos, --tags, or --interactive")
	}

	wtManager := worktree.NewManager(cfg)

	opts := worktree.CreateOptions{
		Branch:      createBranch,
		Repos:       createRepos,
		Tags:        createTags,
		All:         createAll,
		Interactive: createInteractive,
	}

	return wtManager.CreateWorktrees(opts)
}

func runWorktreeList(cmd *cobra.Command, args []string) error {
	wtManager := worktree.NewManager(cfg)

	opts := worktree.ListOptions{
		Branch:  listBranch,
		Repos:   listRepos,
		Verbose: listVerbose,
	}

	return wtManager.ListWorktrees(opts)
}

func runWorktreeDelete(cmd *cobra.Command, args []string) error {
	// Validate that at least one selection method is provided
	if deleteBranch == "" && len(deleteRepos) == 0 && !deleteInteractive {
		return fmt.Errorf("must specify one of: --branch, --repos, or --interactive")
	}

	wtManager := worktree.NewManager(cfg)

	opts := worktree.DeleteOptions{
		Branch:      deleteBranch,
		Repos:       deleteRepos,
		Force:       deleteForce,
		Interactive: deleteInteractive,
	}

	return wtManager.DeleteWorktrees(opts)
}

func runWorktreeSwitch(cmd *cobra.Command, args []string) error {
	wtManager := worktree.NewManager(cfg)
	return wtManager.SwitchWorktree(switchBranch)
}
