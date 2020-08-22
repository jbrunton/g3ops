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
			return lib.Build(context, container)
		},
	}
	return cmd
}
