// Package mod implements the mod command and subcommands.
package mod

import "github.com/spf13/cobra"

// Cmd is the implementation of the mod command.
var Cmd = &cobra.Command{
	Use:   "mod",
	Short: "GX module commands",
}

func init() {
	Cmd.AddCommand(cmdInit)
	Cmd.AddCommand(cmdTidy)
}
