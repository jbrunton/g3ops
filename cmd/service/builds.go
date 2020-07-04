package service

import (
	"os"

	"github.com/jbrunton/g3ops/lib"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var buildsCmd = &cobra.Command{
	Use: "builds",
}

var lsBuildsCmd = &cobra.Command{
	Use:   "ls <service>",
	Short: "Lists builds for the given service",
	Args:  lib.ValidateArgs([]lib.ArgValidator{lib.ServiceValidator}),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]
		catalog := lib.LoadBuildCatalog(service)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Version", "Build Time", "Build SHA", "Image", "ID"})
		table.SetColumnColor(
			tablewriter.Colors{tablewriter.FgYellowColor},
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{},
			tablewriter.Colors{})
		for _, build := range catalog.Builds {
			table.Append([]string{build.Version, build.FormatTimestamp(), build.BuildSha, build.ImageTag, build.ID})
		}
		table.Render()
	},
}

func init() {
	buildsCmd.AddCommand(lsBuildsCmd)
}
