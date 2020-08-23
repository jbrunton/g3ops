package cmd

import (
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
			container.Logger.Printfln("Manifest version: %s", manifest.Version)
			build := lib.FindBuild(manifest.Version, context)
			if os.Getenv("CI") == "1" {
				container.Logger.Printfln("Running in CI environment")
			} else {
				container.Logger.Printfln("Not a CI environment")
			}
			if build != nil {
				container.Logger.Printfln("Build %s found for version %s", build.ID, manifest.Version)
				if os.Getenv("CI") == "1" {
					container.Logger.Println("::set-output name=buildRequired::0")
				}
			} else {
				container.Logger.Printfln("No build found for version %s, build required", manifest.Version)
				if os.Getenv("CI") == "1" {
					container.Logger.Println("::set-output name=buildRequired::1")
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
