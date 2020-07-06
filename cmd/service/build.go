package service

import (
	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
)

func newBuildServiceCmd(executor lib.Executor) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build <service>",
		Short: "Build the given service",
		Args:  lib.ValidateArgs([]lib.ArgValidator{lib.ServiceValidator}),
		Run: func(cmd *cobra.Command, args []string) {
			service := args[0]
			fs := lib.CreateOsFs()
			context, err := lib.GetContext(fs, cmd)
			if err != nil {
				panic(err)
			}
			serviceManifest, err := context.LoadServiceManifest(service)
			if err != nil {
				panic(err)
			}

			lib.Build(service, serviceManifest.Version, context, executor)
		},
	}
	return cmd
}
