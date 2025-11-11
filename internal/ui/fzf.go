package ui

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/asocpro/repo-manager/internal/provider"
)

// FzfSelector provides fzf-based interactive selection
type FzfSelector struct{}

// NewFzfSelector creates a new fzf selector
func NewFzfSelector() *FzfSelector {
	return &FzfSelector{}
}

// CheckFzfInstalled checks if fzf is available
func (f *FzfSelector) CheckFzfInstalled() bool {
	cmd := exec.Command("fzf", "--version")
	return cmd.Run() == nil
}

// SelectRepositories allows multi-select of repositories using fzf
func (f *FzfSelector) SelectRepositories(repos []provider.Repository) ([]provider.Repository, error) {
	if !f.CheckFzfInstalled() {
		return nil, fmt.Errorf("fzf is not installed. Install it with: sudo dnf install fzf (or apt/brew/etc)")
	}

	if len(repos) == 0 {
		return nil, fmt.Errorf("no repositories to select from")
	}

	// Build input for fzf
	var input strings.Builder
	repoMap := make(map[string]provider.Repository)

	for _, repo := range repos {
		// Format: "fullpath | description | stars: X | last activity: YYYY-MM-DD"
		line := fmt.Sprintf("%s", repo.FullPath)

		if repo.Description != "" {
			line += fmt.Sprintf(" | %s", truncate(repo.Description, 60))
		}

		if repo.Stars > 0 {
			line += fmt.Sprintf(" | ⭐ %d", repo.Stars)
		}

		if !repo.LastActivity.IsZero() {
			line += fmt.Sprintf(" | 📅 %s", repo.LastActivity.Format("2006-01-02"))
		}

		if repo.Archived {
			line += " | 📦 archived"
		}

		input.WriteString(line + "\n")
		repoMap[repo.FullPath] = repo
	}

	// Run fzf with multi-select
	cmd := exec.Command("fzf",
		"--multi",
		"--height=80%",
		"--layout=reverse",
		"--border",
		"--prompt=Select repositories> ",
		"--preview=echo {}",
		"--preview-window=down:3:wrap",
		"--header=Press TAB to select multiple, ENTER to confirm",
	)

	cmd.Stdin = strings.NewReader(input.String())

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		// User cancelled
		return nil, nil
	}

	// Parse selected repositories
	selected := []provider.Repository{}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Extract the full path (first field before |)
		fullPath := strings.TrimSpace(strings.Split(line, "|")[0])

		if repo, ok := repoMap[fullPath]; ok {
			selected = append(selected, repo)
		}
	}

	return selected, nil
}

// SelectSingle allows single selection using fzf
func (f *FzfSelector) SelectSingle(items []string, prompt string) (string, error) {
	if !f.CheckFzfInstalled() {
		return "", fmt.Errorf("fzf is not installed")
	}

	if len(items) == 0 {
		return "", fmt.Errorf("no items to select from")
	}

	cmd := exec.Command("fzf",
		"--height=40%",
		"--layout=reverse",
		"--border",
		fmt.Sprintf("--prompt=%s> ", prompt),
	)

	cmd.Stdin = strings.NewReader(strings.Join(items, "\n"))

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		// User cancelled
		return "", nil
	}

	return strings.TrimSpace(out.String()), nil
}

// SelectProvider allows interactive provider selection
func (f *FzfSelector) SelectProvider(providers []string) (string, error) {
	return f.SelectSingle(providers, "Select provider")
}

// SelectGroup allows interactive group selection
func (f *FzfSelector) SelectGroup(groups []provider.Group) (string, error) {
	if !f.CheckFzfInstalled() {
		return "", fmt.Errorf("fzf is not installed")
	}

	if len(groups) == 0 {
		return "", fmt.Errorf("no groups to select from")
	}

	// Build input for fzf
	var input strings.Builder
	groupMap := make(map[string]string)

	for _, group := range groups {
		line := group.FullPath
		if group.Description != "" {
			line += fmt.Sprintf(" | %s", truncate(group.Description, 80))
		}

		input.WriteString(line + "\n")
		groupMap[strings.Split(line, "|")[0]] = group.FullPath
	}

	cmd := exec.Command("fzf",
		"--height=40%",
		"--layout=reverse",
		"--border",
		"--prompt=Select group/organization> ",
	)

	cmd.Stdin = strings.NewReader(input.String())

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		// User cancelled
		return "", nil
	}

	selected := strings.TrimSpace(out.String())
	fullPath := strings.TrimSpace(strings.Split(selected, "|")[0])

	if path, ok := groupMap[fullPath]; ok {
		return path, nil
	}

	return fullPath, nil
}

// SelectWorktree allows interactive worktree selection
func (f *FzfSelector) SelectWorktree(worktrees []string) (string, error) {
	return f.SelectSingle(worktrees, "Select worktree")
}

// truncate truncates a string to the specified length with ellipsis
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// FormatRepoList formats repositories for display
func FormatRepoList(repos []provider.Repository) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Found %d repositories:\n\n", len(repos)))
	sb.WriteString(strings.Repeat("=", 80) + "\n")

	for i, repo := range repos {
		sb.WriteString(fmt.Sprintf("%3d. %s\n", i+1, repo.FullPath))

		if repo.Description != "" {
			sb.WriteString(fmt.Sprintf("     %s\n", truncate(repo.Description, 70)))
		}

		details := []string{}
		if repo.Stars > 0 {
			details = append(details, fmt.Sprintf("⭐ %d", repo.Stars))
		}
		if !repo.LastActivity.IsZero() {
			details = append(details, fmt.Sprintf("📅 %s", repo.LastActivity.Format("2006-01-02")))
		}
		if repo.Archived {
			details = append(details, "📦 archived")
		}

		if len(details) > 0 {
			sb.WriteString(fmt.Sprintf("     %s\n", strings.Join(details, " | ")))
		}

		sb.WriteString("\n")
	}

	sb.WriteString(strings.Repeat("=", 80) + "\n")

	return sb.String()
}
