/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// err := viper.ReadInConfig() // Find and read the config file
		// if err != nil {             // Handle errors reading the config file
		// 	panic(fmt.Errorf("Fatal error config file: %s \n", err))
		// }
		viper.ReadInConfig()
		
		// var contexts []G3opsContext
		// err := viper.UnmarshalKey("contexts", &contexts)
		// if err != nil {
		// 	fmt.Printf("unable to decode into struct, %v", err)
		// 	err = nil
		// }
		contexts := getG3opsContexts()
		for _, ctx := range contexts {
			fmt.Println("  name:", ctx.Name)
			fmt.Println("  url:", ctx.Url)
			fmt.Println()
		}
	},
}

func init() {
	contextCmd.AddCommand(lsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
