package outputs

import (
	"github.com/spf13/cobra"
)

// OutputsCmd represents the context command
var OutputsCmd = &cobra.Command{
	Use: "set-outputs",
}

func init() {
	OutputsCmd.AddCommand(newBuildMatrixCmd())
}
