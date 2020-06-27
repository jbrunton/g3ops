package service

import (
	"errors"
	"os"

	"github.com/jbrunton/g3ops/cmd/styles"
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
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(styles.StyleError("Argument <service> required"))
		}

		if len(args) > 1 {
			return errors.New(styles.StyleError("Unexpected arguments, only <service> expected"))
		}

		context, err := lib.GetContext(cmd)
		if err != nil {
			panic(err)
		}

		var serviceNames []string

		for serviceName := range context.Config.Services {
			if serviceName == args[0] {
				return nil
			}
			serviceNames = append(serviceNames, serviceName)
		}

		return errors.New(styles.StyleError(`Unknown service "` + args[0] + `". Valid options: ` + styles.StyleEnumOptions(serviceNames) + "."))
	},
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]
		catalog := lib.LoadBuildCatalog(service)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Version", "Build Time", "Build SHA", "Image", "ID"})
		table.SetColumnColor(
			tablewriter.Colors{tablewriter.FgGreenColor},
			tablewriter.Colors{tablewriter.FgYellowColor},
			tablewriter.Colors{tablewriter.FgYellowColor},
			tablewriter.Colors{tablewriter.FgYellowColor},
			tablewriter.Colors{tablewriter.FgYellowColor})
		for _, build := range catalog.Builds {
			table.Append([]string{build.Version, build.FormatTimestamp(), build.BuildSha, build.ImageTag, build.ID})
		}
		table.Render() // Send output
	},
}

func init() {
	buildsCmd.AddCommand(lsBuildsCmd)
}
