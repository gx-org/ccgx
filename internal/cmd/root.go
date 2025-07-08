package cmd

import (
	"ccgx/internal/cmd/mod"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ccgx",
	Short: "Generate C++ bindings for GX.",
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(mod.Cmd)
}
