package service

import (
	"fmt"

	"github.com/jbrunton/g3ops/cmd/context"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists services in the current context",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, err := context.LoadContextManifest()
		if err == nil {
			for serviceName, _ := range ctx.Services {
				fmt.Println(serviceName)
			}
		} else {
			fmt.Println("No current context found")
		}
	},
}
