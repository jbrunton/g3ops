package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jbrunton/g3ops/lib"

	"github.com/olekukonko/tablewriter"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func newListReleasesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "Lists services in the current context",
		Run: func(cmd *cobra.Command, args []string) {
			fs := lib.CreateOsFs()
			g3ops, err := lib.GetContext(fs, cmd)
			if err != nil {
				panic(err)
			}
			token := os.Getenv("GITHUB_TOKEN")
			var client *github.Client
			if token != "" {
				ts := oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: token},
				)
				tc := oauth2.NewClient(context.Background(), ts)
				client = github.NewClient(tc)
			} else {
				fmt.Println("Warning: no GITHUB_TOKEN set. g3ops won't be able to authenticate, and some functionality won't be supported.")
				client = github.NewClient(nil)
			}

			releases, _, err := client.Repositories.ListReleases(context.Background(), g3ops.RepoOwnerName, g3ops.RepoName, nil)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Assets", "Status"})
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
				)
				var status string
				if *release.Draft {
					status = "DRAFT"
				} else {
					status = "PUBLISHED"
				}

				row := []string{*release.Name, strings.Join(assets, ", "), status}
				table.Append(row)
			}
			table.Render()
		},
	}
	return cmd
}

func newReleasesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "releases",
	}
	cmd.AddCommand(newListReleasesCmd())
	return cmd
}
