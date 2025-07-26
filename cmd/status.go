package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current stack status with tree visualization",
	Long: `Display the current stack and PR status with a beautiful tree visualization.
	
Shows branch relationships, PR statuses, CI states, and merge readiness
for the entire dependency tree.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showStackStatus()
	},
}

func showStackStatus() error {
	// TODO: Implement stack status visualization
	fmt.Println("Stack Status:")
	fmt.Println("main")
	fmt.Println("└── feature-1 ✓ #123 Ready to merge")
	fmt.Println("    └── feature-2 🔄 #124 CI running")
	fmt.Println("        └── feature-3 📝 #125 Draft")
	return nil
}