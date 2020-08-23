package cmd

import (
	"errors"

	"github.com/jbrunton/g3ops/cmd/styles"
	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
)

func newDeployCmd(container *lib.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy given version to environment",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New(styles.StyleError("Arguments <version> <environment> expected"))
			}

			if len(args) > 2 {
				return errors.New(styles.StyleError("Unexpected arguments, only <version> <environment> expected"))
			}

			fs := container.FileSystem
			context, err := lib.GetContext(fs, cmd)
			if err != nil {
				return err
			}

			manifest, err := context.GetManifest(fs)
			if err != nil {
				return err
			}

			var environments []string

			for envName := range manifest.Environments {
				if envName == args[1] {
					return nil
				}
				environments = append(environments, envName)
			}

			return errors.New(styles.StyleError(`Unknown environment "` + args[1] + `". Valid options: ` + styles.StyleEnumOptions(environments) + "."))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			context, err := lib.GetContext(container.FileSystem, cmd)
			if err != nil {
				return err
			}
			return lib.Deploy(context, container, args[0], args[1])
			//return lib.Build(context, container)
		},
	}
	return cmd
}
