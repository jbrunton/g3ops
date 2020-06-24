package service

import (
	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// GetCmd - g3ops context get
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Prints the current g3ops context",
	Run: func(cmd *cobra.Command, args []string) {
		manifest, err := lib.LoadServiceManifest(args[0])
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
