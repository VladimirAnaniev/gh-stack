package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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
	// TODO: Implement stack branch creation
	fmt.Printf("Creating stacked branch: %s\n", branchName)
	return nil
}