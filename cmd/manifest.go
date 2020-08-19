package cmd

import (
	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
)

func newManifestCheckCmd(container *lib.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check manifest status",
		RunE: func(cmd *cobra.Command, args []string) error {
			container.Logger.Printfln("testing 2")
			return nil
		},
	}
	return cmd
}

func newManifestCmd(container *lib.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "manifest",
		Short: "Manifest info",
	}
	cmd.AddCommand(newManifestCheckCmd(container))
	return cmd
}
