package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd builds the root command for the node-creator CLI.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node-creator",
		Short: "Node.js monorepo project scaffolding tool",
	}

	cmd.AddCommand(newGenerateCmd())

	return cmd
}

// Execute runs the CLI application.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
