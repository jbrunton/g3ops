package environment

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jbrunton/g3ops/lib"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func newLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "Lists environments in the current context",
		Run: func(cmd *cobra.Command, args []string) {
			context, err := lib.GetContext(cmd)
			if err != nil {
				panic(err)
			}
			if err == nil {
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "Manifest"})
				table.SetColumnColor(
					tablewriter.Colors{tablewriter.FgGreenColor},
					tablewriter.Colors{tablewriter.FgYellowColor})
				for envName, env := range context.Config.Environments {
					//fmt.Println(serviceName)
					cwd, err := os.Getwd()
					if err != nil {
						panic(err)
					}
					relPath, err := filepath.Rel(cwd, env.Manifest)
					table.Append([]string{envName, relPath})
				}
				table.Render() // Send output
			} else {
				fmt.Println("No current context found")
			}
		},
	}
}
