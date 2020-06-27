package context

import (
	"fmt"

	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Prints the current g3ops context",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := lib.LoadContextConfig()
		if err == nil {
			fmt.Println(config.Name)
		} else {
			fmt.Println("No current context found")
		}
	},
}
