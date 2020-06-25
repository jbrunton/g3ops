package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/logrusorgru/aurora"
)

var flagsRegex *regexp.Regexp

// StyleUsage - styles the usage template to include color
func StyleUsage(cmd *cobra.Command) {
	usageTemplate := cmd.UsageTemplate()
	usageTemplate = strings.NewReplacer(
		`{{.UseLine}}`, `{{StyleCommand .UseLine}}`,
		`{{.CommandPath}}`, `{{StyleCommand .CommandPath}}`,
		`[command]`, `{{StyleCommand "[command]"}}`,
		`{{rpad .Name .NamePadding }}`, `{{rpad .Name .NamePadding | StyleCommand}}`,
		`.FlagUsages |`, `.FlagUsages | StyleFlags |`,
	).Replace(usageTemplate)
	headingRegex := regexp.MustCompile(`(?m)^\w.*:`)
	usageTemplate = headingRegex.ReplaceAllStringFunc(usageTemplate, func(heading string) string {
		return fmt.Sprintf(`{{Heading "%s"}}`, heading)
	})
	cmd.SetUsageTemplate(usageTemplate)
}

func styleHeading(s string) aurora.Value {
	return aurora.Bold(s)
}

func styleCommand(s string) aurora.Value {
	return aurora.Green(s).Bold()
}

func styleOptions(s string) aurora.Value {
	return aurora.Yellow(s).Bold()
}

func styleFlags(s string) string {
	var styledUsages []string
	for _, flagUsage := range strings.Split(s, "\n") {
		styledUsage := flagsRegex.ReplaceAllStringFunc(flagUsage, func(flag string) string {
			return styleOptions(flag).String()
		})
		styledUsages = append(styledUsages, styledUsage)
	}
	return strings.Join(styledUsages, "\n")
}

func init() {
	flagsRegex = regexp.MustCompile(`^\s+-\S,\s+--\S+|^\s+--\S+`)
	cobra.AddTemplateFunc("Heading", styleHeading)
	cobra.AddTemplateFunc("StyleCommand", styleCommand)
	cobra.AddTemplateFunc("StyleOptions", styleOptions)
	cobra.AddTemplateFunc("StyleFlags", styleFlags)
}
