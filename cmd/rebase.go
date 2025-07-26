package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	rebaseCascade bool
)

var rebaseCmd = &cobra.Command{
	Use:   "rebase",
	Short: "Rebase current branch with cascading updates",
	Long: `Rebase the current branch and optionally cascade changes through dependent branches.
	
When --cascade is used, automatically rebases all dependent branches in the stack
to maintain proper relationships and avoid conflicts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rebaseStack(rebaseCascade)
	},
}

func init() {
	rebaseCmd.Flags().BoolVar(&rebaseCascade, "cascade", false, "Cascade rebase through dependent branches")
}

func rebaseStack(cascade bool) error {
	// TODO: Implement stack rebase logic
	fmt.Printf("Rebasing stack (cascade: %v)\n", cascade)
	return nil
}
