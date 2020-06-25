package styles

import (
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

func init() {
	cobra.AddTemplateFunc("Heading", StyleHeading)
	cobra.AddTemplateFunc("StyleCommand", StyleCommand)
	cobra.AddTemplateFunc("StyleCommandUsage", StyleCommandUsage)
	cobra.AddTemplateFunc("StyleOptions", StyleOptions)
	cobra.AddTemplateFunc("StyleFlags", StyleFlags)
}
