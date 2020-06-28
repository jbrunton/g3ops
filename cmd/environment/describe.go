package environment

import (
	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func newDescribeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "describe",
		Short: "Prints the current g3ops context",
		Run: func(cmd *cobra.Command, args []string) {
			context, err := lib.GetContext(cmd)
			if err != nil {
				panic(err)
			}
			manifest, err := context.LoadEnvironmentManifest(args[0])
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
}
