package context

import (
	"github.com/spf13/cobra"
)

// ContextCmd represents the context command
var ContextCmd = &cobra.Command{
	Use:   "context",
	Short: "A brief description of your command",
}

func init() {
	ContextCmd.AddCommand(newContextGetCmd())
	ContextCmd.AddCommand(describeCmd)
}
