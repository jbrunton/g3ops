package workflow

import (
	"fmt"
	"os"

	"github.com/jbrunton/g3ops/cmd/styles"

	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
)

func newGenerateWorkflowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generates workflow files",
		RunE: func(cmd *cobra.Command, args []string) error {
			context, err := lib.GetContext(cmd)
			if err != nil {
				fmt.Println(styles.StyleError(err.Error()))
				os.Exit(1)
			}
			valuesPath := context.Config.Ci.Workflows.Build.Values
			//targetPath := context.Config.Ci.Workflows.Build.Target
			_, err = os.Stat(valuesPath)
			if err != nil {
				fmt.Println(styles.StyleError(err.Error()))
				os.Exit(1)
			}
			lib.ExecCommand(fmt.Sprintf("ytt -f %s > %s",
				context.Config.Ci.Workflows.Build.Values,
				context.Config.Ci.Workflows.Build.Target), context)

			return nil
		},
	}
}

// WorkflowCmd represents the context command
var WorkflowCmd = &cobra.Command{
	Use: "workflow",
}

func init() {
	WorkflowCmd.AddCommand(newGenerateWorkflowCmd())
}
