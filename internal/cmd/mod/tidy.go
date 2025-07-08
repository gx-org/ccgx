package mod

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdTidy = &cobra.Command{
	Use:   "tidy",
	Short: "update gx.mod",
	RunE:  cTidy,
}

func cTidy(cmd *cobra.Command, args []string) error {
	fmt.Println("Running tidy")
	return nil
}
