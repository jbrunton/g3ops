package context

import (
	"fmt"

	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// G3opsContext - type of current g3ops context
type G3opsContext struct {
	Name string
}

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

func init() {
	//ContextCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
