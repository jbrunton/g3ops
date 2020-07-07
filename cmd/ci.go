package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func checkReleaseManifest(fs *afero.Afero, gitHubService lib.GitHubService, g3ops *lib.G3opsContext, cmd *cobra.Command) {
	releaseManifest, err := g3ops.GetReleaseManifest(fs)
	if err != nil {
		panic(err)
	}
	expectedVersion := releaseManifest.Version

	releases, err := gitHubService.ListReleases(g3ops)

	currentVersion := *releases[0].Name

	if currentVersion == expectedVersion {
		fmt.Fprintf(cmd.OutOrStdout(), "Release %q already exists\n", expectedVersion)
		fmt.Fprintf(cmd.OutOrStdout(), "::set-output name=buildRequired::0\n")
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Release %q not found, build required\n", expectedVersion)
		fmt.Fprintf(cmd.OutOrStdout(), "::set-output name=buildRequired::1\n")
	}
}

func checkServiceManifests(fs *afero.Afero, g3ops *lib.G3opsContext, cmd *cobra.Command) {
	buildTasks := []map[string]string{}
	for serviceName := range g3ops.Config.Services {
		serviceManifest, err := g3ops.LoadServiceManifest(serviceName)
		if err != nil {
			panic(err)
		}
		if lib.BuildExists(serviceName, serviceManifest.Version) {
			fmt.Fprintf(cmd.OutOrStdout(), "Build %q already exists for service %q\n", serviceManifest.Version, serviceName)
		} else {
			fmt.Fprintf(cmd.OutOrStdout(), "Build %q required for service %q\n", serviceManifest.Version, serviceName)
			buildTasks = append(buildTasks, map[string]string{
				"service": serviceName,
			})
		}
	}
	buildMatrix := map[string]interface{}{
		"include": buildTasks,
	}
	json, err := json.Marshal(&buildMatrix)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "::set-output name=buildMatrix::%s\n", json)
	if len(buildTasks) > 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "::set-output name=buildRequired::1\n")
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "::set-output name=buildRequired::0\n")
	}
}

func newBuildMatrixCmd(container *lib.Container) *cobra.Command {
	return &cobra.Command{
		Use:       "check",
		Short:     "Sets buildMatrix output describing any builds required",
		ValidArgs: []string{"release-manifest", "service-manifests"},
		Args:      cobra.OnlyValidArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fs := container.FileSystem
			g3ops, err := lib.GetContext(fs, cmd)
			if err != nil {
				panic(err)
			}

			switch args[0] {
			case "release-manifest":
				checkReleaseManifest(fs, container.GitHubService, g3ops, cmd)
			case "service-manifests":
				checkServiceManifests(fs, g3ops, cmd)
			default:
				panic(fmt.Errorf("Unexpected check: %q", args[0]))
			}
		},
	}
}

func newCiCmd(container *lib.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "ci",
		Short:  "Commands used by CI pipelines",
		Hidden: true,
	}
	cmd.AddCommand(newBuildMatrixCmd(container))
	return cmd
}
