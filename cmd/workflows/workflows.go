package workflows

import (
	"fmt"
	"os"

	"github.com/jbrunton/g3ops/cmd/styles"

	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
)

func newUpdateWorkflowsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Updates workflow files",
		Run: func(cmd *cobra.Command, args []string) {
			context, err := lib.GetContext(cmd)
			if err != nil {
				fmt.Println(styles.StyleError(err.Error()))
				os.Exit(1)
			}
			fs := lib.CreateOsFs()
			lib.UpdateWorkflows(fs, context)
		},
	}
}

func newCheckWorkflowsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Check workflow files are up to date",
		Run: func(cmd *cobra.Command, args []string) {
			context, err := lib.GetContext(cmd)
			if err != nil {
				fmt.Println(styles.StyleError(err.Error()))
				os.Exit(1)
			}
			fs := lib.CreateOsFs()
			err = lib.ValidateWorkflows(fs, context)
			if err != nil {
				fmt.Println(styles.StyleError(err.Error()))
				os.Exit(1)
			} else {
				fmt.Println(styles.StyleCommand("Workflows up to date"))
			}
		},
	}
}

func newInitWorkflowsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize g3ops workflow",
		Run: func(cmd *cobra.Command, args []string) {
			context, err := lib.GetContext(cmd)
			if err != nil {
				fmt.Println(styles.StyleError(err.Error()))
				os.Exit(1)
			}

			fs := lib.CreateOsFs()
			lib.InitWorkflows(fs, context)
		},
	}
}

// WorkflowsCmd represents the context command
var WorkflowsCmd = &cobra.Command{
	Use: "workflows",
}

func init() {
	WorkflowsCmd.AddCommand(newUpdateWorkflowsCmd())
	WorkflowsCmd.AddCommand(newCheckWorkflowsCmd())
	WorkflowsCmd.AddCommand(newInitWorkflowsCmd())
}
