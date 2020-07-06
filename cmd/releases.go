package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbrunton/g3ops/lib"

	"github.com/olekukonko/tablewriter"

	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func newGithubClient() *github.Client {
	token := os.Getenv("GITHUB_TOKEN")

	if token == "" {
		fmt.Println("Warning: no GITHUB_TOKEN set. g3ops won't be able to authenticate, and some functionality won't be supported.")
		return github.NewClient(nil)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

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

			client := newGithubClient()
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

func newCreateReleaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Updates manifest to trigger to a new release",
		Run: func(cmd *cobra.Command, args []string) {
			fs := lib.CreateOsFs()
			g3ops, err := lib.GetContext(fs, cmd)
			if err != nil {
				panic(err)
			}

			dir, err := ioutil.TempDir("", strings.Join([]string{"g3ops", g3ops.RepoName, "*"}, "-"))
			if err != nil {
				log.Fatal(err)
			}

			lib.ExecCommand(fmt.Sprintf("git clone --depth 1 git@github.com:%s.git %s", g3ops.Config.Repo, dir), g3ops)

			// TODO: specify context, not config, and require .g3ops directory in context
			newContext, err := lib.NewContext(fs, filepath.Join(dir, ".g3ops", "config.yml"), g3ops.DryRun)
			if err != nil {
				panic(err)
			}
			manifest, err := newContext.GetReleaseManifest()
			if err != nil {
				panic(err)
			}

			fmt.Println("Version:", manifest.Version)

			defer os.RemoveAll(dir)
		},
	}
	return cmd
}

func newReleasesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "releases",
	}
	cmd.AddCommand(newListReleasesCmd())
	cmd.AddCommand(newCreateReleaseCmd())
	return cmd
}
