package cmd

import (
	"encoding/json"
	"os"

	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
)

func newManifestCheckCmd(container *lib.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check manifest status",
		RunE: func(cmd *cobra.Command, args []string) error {
			fs := container.FileSystem
			context, err := lib.GetContext(fs, cmd)
			if err != nil {
				return err
			}
			manifest, err := context.GetManifest(fs)
			if err != nil {
				return err
			}

			if os.Getenv("CI") == "1" {
				container.Logger.Printfln("Running in CI environment.\n")
			} else {
				container.Logger.Printfln("Not a CI environment. Set CI=1 for GitHub workflow output commands.\n")
			}

			container.Logger.Println("Checking build...")
			container.Logger.Printfln("  Manifest version: %s", manifest.Version)
			build := lib.FindBuild(manifest.Version, context)
			if build != nil {
				container.Logger.Printfln("  Existing build: %s, no build required", build.ID)
				if os.Getenv("CI") == "1" {
					container.Logger.Println("\n::set-output name=buildRequired::0")
				}
			} else {
				container.Logger.Printfln("  No build found for version %s, build required", manifest.Version)
				if os.Getenv("CI") == "1" {
					container.Logger.Println("\n::set-output name=buildRequired::1")
				}
			}

			type deploymentTask struct {
				Environment string `json:"environment"`
				Version     string `json:"version"`
			}
			deploymentTasks := []deploymentTask{}
			type deploymentMatrix struct {
				Include []deploymentTask `json:"include"`
			}

			for envName, envInfo := range manifest.Environments {
				container.Logger.Printfln("\nChecking %s...", envName)
				latestDeployment := lib.GetLatestDeployment(envName, context)
				container.Logger.Printfln("  Manifest version: %s", envInfo.Version)
				if latestDeployment != nil {
					container.Logger.Printfln("  Deployed version: %s", latestDeployment.Version)
				} else {
					container.Logger.Println("  Deployed version: <none>")
				}
				if latestDeployment == nil || envInfo.Version != latestDeployment.Version {
					container.Logger.Println("  Deployment required")
					deploymentTasks = append(deploymentTasks, deploymentTask{Environment: envName, Version: envInfo.Version})
				} else {
					container.Logger.Println("  No deployment required")
				}
			}
			container.Logger.Println()
			if os.Getenv("CI") == "1" {
				if len(deploymentTasks) > 0 {
					container.Logger.Println("::set-output name=deploymentsRequired::1")
					deploymentMatrix := deploymentMatrix{Include: deploymentTasks}
					deploymentMatrixBytes, err := json.Marshal(deploymentMatrix)
					if err != nil {
						return err
					}
					container.Logger.Printfln("::set-output name=deploymentMatrix::%s", string(deploymentMatrixBytes))
				} else {
					container.Logger.Println("::set-output name=deploymentsRequired::0")
				}
			}
			return nil
		},
	}
	return cmd
}

func newManifestCmd(container *lib.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manifest",
		Short: "Manifest info",
	}
	cmd.AddCommand(newManifestCheckCmd(container))
	return cmd
}
