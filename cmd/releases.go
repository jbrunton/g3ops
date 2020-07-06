package cmd

import (
	"context"
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

			client := lib.NewGithubClient()
			releases, _, err := client.Repositories.ListReleases(context.Background(), g3ops.RepoOwnerName, g3ops.RepoName, nil)
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

			lib.CreateNewRelease(fs, container.Executor, g3ops)
		},
	}
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
