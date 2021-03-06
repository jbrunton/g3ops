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
	"os"

	"github.com/jbrunton/g3ops/lib"

	"github.com/spf13/cobra"

	"github.com/jbrunton/g3ops/cmd/context"
)

var cfgFile string

var rootCmd *cobra.Command

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Version - the build version
var Version = "development"

func newVersionCmd(container *lib.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		RunE: func(cmd *cobra.Command, args []string) error {
			container.Logger.Printfln("g3ops version %s", Version)
			return nil
		},
	}
	return cmd
}

// NewRootCommand creates a new root command
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "g3ops",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) {},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.PersistentFlags().Bool("dry-run", false, "Preview commands before executing")
	cmd.PersistentFlags().StringP("config", "c", "", "Location of g3ops context config")

	return cmd
}

func init() {
	//cobra.OnInitialize(initConfig)

	rootCmd = NewRootCommand()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	container := lib.NewContainer()
	rootCmd.AddCommand(context.ContextCmd)
	rootCmd.AddCommand(newManifestCmd(container))
	rootCmd.AddCommand(newBuildCmd(container))
	rootCmd.AddCommand(newDeployCmd(container))
	rootCmd.AddCommand(newVersionCmd(container))
}
