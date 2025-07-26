package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	mergeAutoCascade bool
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge ready PRs in dependency order",
	Long: `Merge ready PRs in the correct dependency order.
	
Automatically determines merge order based on stack dependencies and
merges PRs that are ready (approved, CI passing, etc.). With --auto-cascade,
continues merging dependent PRs as they become ready.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return mergeStack(mergeAutoCascade)
	},
}

func init() {
	mergeCmd.Flags().BoolVar(&mergeAutoCascade, "auto-cascade", false, "Automatically merge dependent PRs as they become ready")
}

func mergeStack(autoCascade bool) error {
	// TODO: Implement stack merge logic
	fmt.Printf("Merging stack (auto-cascade: %v)\n", autoCascade)
	return nil
}
