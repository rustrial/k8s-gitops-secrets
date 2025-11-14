package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "seals",
	Short: "Seal secrets tool",
	Long:  `Seal (encrypt) Kubernetes secrets.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// Execute is the main entry point into seals CLI.
func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
