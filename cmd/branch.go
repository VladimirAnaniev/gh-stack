package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vladimir-ananiev/gh-stacked/pkg/git"
)

var branchCmd = &cobra.Command{
	Use:   "branch <branch-name>",
	Short: "Create a new stacked branch",
	Long: `Create a new branch in the current stack.
	
If run from main/master, creates the first branch in a new stack.
If run from an existing stacked branch, creates a child branch.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branchName := args[0]
		return createStackedBranch(branchName)
	},
}

func createStackedBranch(branchName string) error {
	// Get current working directory (repository path)
	repoPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create branch service
	branchService := git.NewBranchService(repoPath)

	// Get current branch to use as parent
	currentBranch, err := branchService.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	fmt.Printf("Creating branch '%s' from '%s'...\n", branchName, currentBranch)

	// Create the new branch
	err = branchService.CreateBranch(branchName, currentBranch)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	fmt.Printf("✓ Created and switched to branch '%s'\n", branchName)
	return nil
}
