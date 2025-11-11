package cmd

import (
	"fmt"
	"os"

	"github.com/asocpro/repo-manager/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	version = "0.1.0"
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "repo-manager",
	Short: "A tool for managing git repositories and worktrees",
	Long: `Repository Manager helps you manage multiple git repositories
and their worktrees organized by branch names.

It supports cloning groups from GitLab, GitHub, Gitea, and generic Git providers,
creating branch-based worktrees, and provides both CLI and TUI interfaces.`,
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/repo-manager/config.yaml)")
}

func initConfig() {
	var err error
	cfg, err = config.Load(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid configuration: %v\n", err)
		os.Exit(1)
	}
}
