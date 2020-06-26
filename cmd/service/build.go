package service

import (
	"errors"

	"github.com/jbrunton/cobra"
	"github.com/jbrunton/g3ops/cmd/styles"
	"github.com/jbrunton/g3ops/lib"
)

var buildCmd = &cobra.Command{
	Use:   "build <service>",
	Short: "Build the given service",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(styles.StyleError("Argument <service> required"))
		}

		if len(args) > 1 {
			return errors.New(styles.StyleError("Unexpected arguments, only <service> expected"))
		}

		ctx, err := lib.LoadContextManifest()
		if err != nil {
			panic(err)
		}

		var serviceNames []string

		for serviceName := range ctx.Services {
			if serviceName == args[0] {
				return nil
			}
			serviceNames = append(serviceNames, serviceName)
		}

		return errors.New(styles.StyleError(`Unknown service "` + args[0] + `". Valid options: ` + styles.StyleEnumOptions(serviceNames) + "."))
	},
	Run: func(cmd *cobra.Command, args []string) {
		service := args[0]

		serviceManifest, err := lib.LoadServiceManifest(service)
		if err != nil {
			panic(err)
		}

		cmdCtx := lib.GetCommandContext(cmd)
		lib.Build(service, serviceManifest.Version, &cmdCtx)
	},
}
