package service

import (
	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build <service>",
	Short: "Build the given service",
	Args:  lib.ValidateArgs([]lib.ArgValidator{lib.ServiceValidator}),
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]
		context, err := lib.GetContext(cmd)
		if err != nil {
			panic(err)
		}
		serviceManifest, err := context.LoadServiceManifest(service)
		if err != nil {
			panic(err)
		}

		lib.Build(service, serviceManifest.Version, &context)
	},
}
