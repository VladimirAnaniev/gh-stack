package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/vladimir-ananiev/gh-stack/pkg/git"
	"github.com/vladimir-ananiev/gh-stack/pkg/github"
)

var (
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	hintStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
)

var rootCmd = &cobra.Command{
	Use:   "stack",
	Short: "Manage stacked pull requests",
	Long: `A simple CLI tool for managing stacked Pull Request workflows on GitHub.
	
Shows dependency tree of open PRs and handles cascading rebases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showStackStatus()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var cascadeCmd = &cobra.Command{
	Use:   "cascade",
	Short: "Cascade rebase all branches in dependency order",
	Long: `Checkout default branch (main/master), pull, then for each branch with PR targeting the default branch:
	1. Checkout branch, rebase on target, push
	2. For each dependent branch, checkout, rebase, push
	
Handles merged branches by dropping commits already in target.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cascadeRebase()
	},
}

func init() {
	rootCmd.AddCommand(cascadeCmd)
}

func showStackStatus() error {
	ctx := context.Background()

	// Get current branch
	currentBranch, err := git.GetCurrentBranch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// Get all open PRs and build dependency tree
	var prs []*github.PR
	err = spinner.New().
		Title("Fetching pull requests...").
		Action(func() {
			prs, err = github.GetOpenPRs(ctx)
		}).
		Run()
	if err != nil {
		return fmt.Errorf("failed to get open PRs: %w", err)
	}

	tree := github.BuildDependencyTree(prs)
	github.PrintTree(tree, currentBranch)

	return nil
}

func cascadeRebase() error {
	ctx := context.Background()

	// Get current branch to restore later
	currentBranch, err := git.GetCurrentBranch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	// Get all open PRs and build dependency tree
	var prs []*github.PR
	err = spinner.New().
		Title("Fetching pull requests...").
		Action(func() {
			prs, err = github.GetOpenPRs(ctx)
		}).
		Run()
	if err != nil {
		return fmt.Errorf("failed to get open PRs: %w", err)
	}

	tree := github.BuildDependencyTree(prs)

	// Find the tree containing the current branch
	currentTree := github.FindCurrentBranchTree(tree, currentBranch)
	if currentTree == nil {
		fmt.Printf("%s %s has no open PR or is not part of a PR stack\n\n", 
			errorStyle.Render("✗ Error:"), 
			warningStyle.Render(currentBranch))
		fmt.Printf("%s Switch to a branch that has an open PR to use cascade\n", 
			hintStyle.Render("Hint:"))
		return nil // Return nil to prevent cobra from showing the error again
	}

	// Get the base branch for this tree
	baseBranch := currentTree.PR.BaseRefName

	// Checkout base branch and pull
	err = spinner.New().
		Title(fmt.Sprintf("Updating %s...", baseBranch)).
		Action(func() {
			err = git.CheckoutAndPull(ctx, baseBranch)
		}).
		Run()
	if err != nil {
		return fmt.Errorf("failed to update %s: %w", baseBranch, err)
	}

	// Process only the current tree in dependency order
	err = spinner.New().
		Title(fmt.Sprintf("Rebasing %s → %s...", currentTree.PR.HeadRefName, currentTree.PR.BaseRefName)).
		Action(func() {
			err = github.ProcessSingleTreeRebase(ctx, currentTree)
		}).
		Run()
	if err != nil {
		return err
	}

	// Restore original branch
	err = spinner.New().
		Title(fmt.Sprintf("Returning to %s...", currentBranch)).
		Action(func() {
			err = git.CheckoutBranch(ctx, currentBranch)
		}).
		Run()
	if err != nil {
		return fmt.Errorf("failed to restore branch %s: %w", currentBranch, err)
	}

	return nil
}
