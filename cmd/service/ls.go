package service

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/jbrunton/g3ops/lib"
	"github.com/olekukonko/tablewriter"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Lists services in the current context",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := lib.LoadContextConfig()
		if err == nil {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Manifest"})
			table.SetColumnColor(
				tablewriter.Colors{tablewriter.FgGreenColor},
				tablewriter.Colors{tablewriter.FgYellowColor})
			for serviceName, service := range config.Services {
				//fmt.Println(serviceName)
				cwd, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				relPath, err := filepath.Rel(cwd, service.Manifest)
				table.Append([]string{serviceName, relPath})
			}
			table.Render() // Send output
		} else {
			fmt.Println("No current context found")
		}
	},
}
