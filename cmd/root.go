package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vladimir-ananiev/gh-stack/pkg/github"
	"github.com/vladimir-ananiev/gh-stack/pkg/git"
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
	
	// Get the default branch
	defaultBranch, err := git.GetDefaultBranch(ctx)
	if err != nil {
		return fmt.Errorf("failed to get default branch: %w", err)
	}
	
	// Get all open PRs and build dependency tree
	prs, err := github.GetOpenPRs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get open PRs: %w", err)
	}
	
	tree := github.BuildDependencyTree(prs)
	
	// Checkout default branch and pull
	fmt.Printf("Checking out %s and pulling...\n", defaultBranch)
	if err := git.CheckoutAndPull(ctx, defaultBranch); err != nil {
		return fmt.Errorf("failed to update %s: %w", defaultBranch, err)
	}
	
	// Process tree in dependency order
	return github.ProcessCascadeRebase(ctx, tree)
}
