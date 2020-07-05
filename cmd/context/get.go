package context

import (
	"fmt"

	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
)

func newContextGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Prints the current g3ops context",
		Run: func(cmd *cobra.Command, args []string) {
			fs := lib.CreateOsFs()
			context, err := lib.GetContext(fs, cmd)
			if err == nil {
				fmt.Fprintln(cmd.OutOrStdout(), context.Config.Name)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "No current context found")
			}
		},
	}
}
