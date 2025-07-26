package cmd

import (
	"github.com/spf13/cobra"
)

var (
	pushDraft bool
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push current branch and create/update PR",
	Long: `Smart push: opens PR first time, updates existing PR after.
	
Amends commits with stack metadata, pushes to origin, and creates or updates
the GitHub Pull Request with proper stack dependency information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return pushStackedBranch(pushDraft)
	},
}

func init() {
	pushCmd.Flags().BoolVar(&pushDraft, "draft", false, "Create PR as draft")
}

func pushStackedBranch(draft bool) error {

	return nil
}
