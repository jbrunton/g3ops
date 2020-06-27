package context

import (
	"fmt"

	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// GetCmd - g3ops context get
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Prints the current g3ops context",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := lib.LoadContextConfig()
		if err == nil {
			out, err := yaml.Marshal(&config)
			if err == nil {
				lib.PrintYaml(string(out))
			} else {
				panic(err)
			}
		} else {
			fmt.Println("No current context found")
		}
	},
}
