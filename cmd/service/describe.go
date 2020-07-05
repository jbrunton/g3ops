package service

import (
	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// GetCmd - g3ops context get
var describeCmd = &cobra.Command{
	Use:   "describe <service>",
	Short: "Prints the current g3ops context",
	Args:  lib.ValidateArgs([]lib.ArgValidator{lib.ServiceValidator}),
	Run: func(cmd *cobra.Command, args []string) {
		fs := lib.CreateOsFs()
		context, err := lib.GetContext(fs, cmd)
		if err != nil {
			panic(err)
		}
		manifest, err := context.LoadServiceManifest(args[0])
		if err == nil {
			out, err := yaml.Marshal(&manifest)
			if err == nil {
				lib.PrintYaml(string(out))
			} else {
				panic(err)
			}
		}
	},
}
