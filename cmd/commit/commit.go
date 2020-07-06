package commit

import (
	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
)

// NewCommitCmd - new commit command
func NewCommitCmd(container *lib.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use: "commit",
	}
	cmd.AddCommand(newCommitBuildCmd(container))
	return cmd
}
