package environment

import (
	"github.com/spf13/cobra"
)

// EnvironmentCmd represents the context command
var EnvironmentCmd = &cobra.Command{
	Use:   "environment",
	Short: "A brief description of your command",
}

func init() {
	EnvironmentCmd.AddCommand(newLsCmd())
	EnvironmentCmd.AddCommand(newDescribeCmd())
}
