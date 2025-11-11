package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive TUI",
	Long: `Launch the Terminal User Interface for interactive repository management.

The TUI provides:
  - Repository browser
  - Worktree creation and deletion
  - Status overview
  - Sync operations`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("TUI not yet implemented")
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
