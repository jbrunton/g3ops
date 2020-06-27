package commit

import (
	"github.com/spf13/cobra"
)

// CommitCmd represents the context command
var CommitCmd = &cobra.Command{
	Use: "commit",
}

func init() {
	CommitCmd.AddCommand(newCommitBuildCmd())
}
