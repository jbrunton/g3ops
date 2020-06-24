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
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jbrunton/g3ops/cmd/context"
	"github.com/jbrunton/g3ops/cmd/service"
	"github.com/logrusorgru/aurora"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
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
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().Bool("dry-run", false, "Preview commands before executing, also --dry-run")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(context.ContextCmd)
	rootCmd.AddCommand(service.ServiceCmd)
	cobra.AddTemplateFunc("Bold", aurora.Bold)
	cobra.AddTemplateFunc("StyleCommand", func(s string) string { return aurora.Green(s).Bold().String() })
	cobra.AddTemplateFunc("StyleOptions", func(s string) string { return aurora.Yellow(s).Bold().String() })
	flagsRegex := regexp.MustCompile(`^\s+-\S,\s+--\S+|^\s+--\S+`)
	cobra.AddTemplateFunc("StyleFlags", func(flagsUsage string) string {
		var styledUsages []string
		for _, flagUsage := range strings.Split(flagsUsage, "\n") {
			styledUsage := flagsRegex.ReplaceAllStringFunc(flagUsage, func(flag string) string {
				return aurora.Yellow(flag).Bold().String()
			})
			styledUsages = append(styledUsages, styledUsage)
		}
		return strings.Join(styledUsages, "\n")
	})
	rootCmd.SetUsageTemplate(`{{Bold "Usage:"}}{{if .Runnable}}
  {{StyleCommand .UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{StyleCommand .CommandPath}} {{StyleOptions "[command]"}}{{end}}{{if gt (len .Aliases) 0}}
{{Bold "Aliases:"}}
  {{.NameAndAliases}}{{end}}{{if .HasExample}}
{{Bold "Examples:"}}
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}
{{Bold "Available Commands:"}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding | StyleCommand}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}
{{Bold "Flags:"}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces | StyleFlags}}{{end}}{{if .HasAvailableInheritedFlags}}
{{Bold "Global Flags:"}}
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces | StyleFlags}}{{end}}{{if .HasHelpSubCommands}}
{{Bold "Additional help topics:"}}{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}`)
}
