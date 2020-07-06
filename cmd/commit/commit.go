package commit

import (
	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
)

// NewCommitCmd - new commit command
func NewCommitCmd(executor lib.Executor) *cobra.Command {
	cmd := &cobra.Command{
		Use: "commit",
	}
	cmd.AddCommand(newCommitBuildCmd(executor))
	return cmd
}
