package context

import (
	"fmt"

	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// GetCmd - g3ops context get
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Prints the current g3ops context",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := ioutil.ReadFile(".g3ops.yml")
		if err == nil {
			//fmt.Println("Current context: " + string(context))
			ctx := G3opsContext{}
			yaml.Unmarshal(data, &ctx)
			fmt.Println(ctx.Name)
		} else {
			fmt.Println("No current context found")
		}
	},
}
