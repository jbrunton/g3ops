package styles

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jbrunton/cobra"
)

// ConfigureUsageTemplate - styles the usage template to include color
func ConfigureUsageTemplate(cmd *cobra.Command) {
	usageTemplate := cmd.UsageTemplate()
	usageTemplate = strings.NewReplacer(
		`{{.UseLine}}`, `{{StyleCommandUsage .UseLine}}`,
		`{{.CommandPath}}`, `{{StyleCommandUsage .CommandPath}}`,
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

// ConfigureUnknownCommandErrorFunc - configures a new UnknownCommandErrorFunc with color styling
func ConfigureUnknownCommandErrorFunc(cmd *cobra.Command) {
	cmd.SuggestionsMinimumDistance = 2
	cmd.SetUnknownCommandErrorFunc(func(c *cobra.Command, arg string) error {
		errorMessage := StyleError(fmt.Sprintf("unknown command %q for %q", arg, c.CommandPath()))
		suggestionsString := ""
		if suggestions := c.SuggestionsFor(arg); len(suggestions) > 0 {
			suggestionsString += "\n\nDid you mean this?\n"
			for _, s := range suggestions {
				suggestionsString += fmt.Sprintf(StyleCommand("\t%v\n").String(), s)
			}
		}
		return errors.New(errorMessage + suggestionsString)
	})
}

func init() {
	cobra.AddTemplateFunc("Heading", StyleHeading)
	cobra.AddTemplateFunc("StyleCommand", StyleCommand)
	cobra.AddTemplateFunc("StyleCommandUsage", StyleCommandUsage)
	cobra.AddTemplateFunc("StyleOptions", StyleOptions)
	cobra.AddTemplateFunc("StyleFlags", StyleFlags)
}
