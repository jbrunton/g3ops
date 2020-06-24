package context

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Prints the current g3ops context",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, err := LoadContextManifest()
		if err == nil {
			fmt.Println(ctx.Name)
		} else {
			fmt.Println("No current context found")
		}
	},
}
