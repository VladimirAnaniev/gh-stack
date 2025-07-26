package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stacked",
	Short: "Manage stacked pull requests",
	Long: `A CLI tool for managing stacked Pull Request workflows on GitHub.
	
Simplifies creating dependent PRs, cascading rebases, and visualizing PR dependency trees.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(branchCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(cleanupCmd)
	rootCmd.AddCommand(rebaseCmd)
	rootCmd.AddCommand(mergeCmd)
}
