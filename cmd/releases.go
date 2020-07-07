package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jbrunton/g3ops/lib"

	"github.com/olekukonko/tablewriter"

	"github.com/spf13/cobra"
)

func newListReleasesCmd(container *lib.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "Lists services in the current context",
		Run: func(cmd *cobra.Command, args []string) {
			g3ops, err := lib.GetContext(container.FileSystem, cmd)
			if err != nil {
				panic(err)
			}

			releases, err := container.GitHubService.ListReleases(g3ops.RepoID)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Tag", "Assets", "Status"})
			for _, release := range releases {
				assets := []string{}
				for _, asset := range release.Assets {
					assets = append(assets, *asset.Name)
				}
				if len(assets) == 0 {
					assets = []string{"(none)"}
				}
				table.SetColumnColor(
					tablewriter.Colors{tablewriter.FgGreenColor},
					tablewriter.Colors{tablewriter.FgYellowColor},
					tablewriter.Colors{tablewriter.FgYellowColor},
					tablewriter.Colors{tablewriter.FgYellowColor},
				)
				var status string
				if *release.Draft {
					status = "DRAFT"
				} else {
					status = "PUBLISHED"
				}

				row := []string{*release.Name, *release.TagName, strings.Join(assets, ", "), status}
				table.Append(row)
			}
			table.Render()
		},
	}
	return cmd
}

func newCreateReleaseCmd(container *lib.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Updates manifest to trigger to a new release",
		Run: func(cmd *cobra.Command, args []string) {
			fs := lib.CreateOsFs()
			g3ops, err := lib.GetContext(fs, cmd)
			if err != nil {
				panic(err)
			}

			builder := lib.NewReleaseBuilder(container, g3ops)
			increment, error := cmd.Flags().GetString("increment")
			if error != nil {
				panic(error)
			}
			builder.CreateNewRelease(increment)
		},
	}
	cmd.Flags().String("increment", "", "The semver increment type: major, minor or patch")
	return cmd
}

func newReleasesCmd(container *lib.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use: "releases",
	}
	cmd.AddCommand(newListReleasesCmd(container))
	cmd.AddCommand(newCreateReleaseCmd(container))
	return cmd
}
