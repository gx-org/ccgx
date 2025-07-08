package mod

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdInit = &cobra.Command{
	Use:   "init",
	Short: "create gx.mod",
	RunE:  cInit,
}

func cInit(cmd *cobra.Command, args []string) error {
	fmt.Println("Running init")
	return nil
}
