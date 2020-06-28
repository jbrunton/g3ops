package resolve

import (
	"github.com/spf13/cobra"
)

// ResolveCmd represents the resolve command
var ResolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "A brief description of your command",
}

func init() {
	ResolveCmd.AddCommand(newResolveTagsCmd())
}
