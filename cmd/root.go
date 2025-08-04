package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vladimir-ananiev/gh-stack/pkg/git"
	"github.com/vladimir-ananiev/gh-stack/pkg/github"
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
	prs, err := github.GetOpenPRs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get open PRs: %w", err)
	}

	tree := github.BuildDependencyTree(prs)
	fmt.Printf("Stack Status (current: %s)\n\n", currentBranch)
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
	prs, err := github.GetOpenPRs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get open PRs: %w", err)
	}

	tree := github.BuildDependencyTree(prs)

	// Find the tree containing the current branch
	currentTree := github.FindCurrentBranchTree(tree, currentBranch)
	if currentTree == nil {
		return fmt.Errorf("current branch %s not found in any PR tree", currentBranch)
	}

	// Get the base branch for this tree
	baseBranch := currentTree.PR.BaseRefName

	// Checkout base branch and pull
	fmt.Printf("Checking out %s and pulling...\n", baseBranch)
	if err := git.CheckoutAndPull(ctx, baseBranch); err != nil {
		return fmt.Errorf("failed to update %s: %w", baseBranch, err)
	}

	// Process only the current tree in dependency order
	if err := github.ProcessSingleTreeRebase(ctx, currentTree); err != nil {
		return err
	}

	// Restore original branch
	fmt.Printf("Returning to %s...\n", currentBranch)
	if err := git.CheckoutBranch(ctx, currentBranch); err != nil {
		return fmt.Errorf("failed to restore branch %s: %w", currentBranch, err)
	}

	return nil
}
