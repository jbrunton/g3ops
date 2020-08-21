package cmd

import (
	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
)

func newBuildCmd(container *lib.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Create build for manifest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fs := lib.CreateOsFs()
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
			if build != nil {
				container.Logger.Printfln("Build %s found for version %s, skipping", build.ID, manifest.Version)
			} else {
				lib.Build(manifest.Version, container.FileSystem, context, container.Executor)
			}
			return nil
		},
	}
	return cmd
}
