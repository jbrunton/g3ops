package outputs

import (
	"encoding/json"
	"fmt"

	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
)

func newBuildMatrixCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "build-matrix",
		Short: "Sets buildMatrix output describing any builds required",
		Run: func(cmd *cobra.Command, args []string) {
			context, err := lib.GetContext(cmd)
			if err != nil {
				panic(err)
			}

			buildTasks := []map[string]string{}
			for serviceName := range context.Config.Services {
				serviceManifest, err := context.LoadServiceManifest(serviceName)
				if err != nil {
					panic(err)
				}
				if lib.BuildExists(serviceName, serviceManifest.Version) {
					fmt.Fprintf(cmd.OutOrStdout(), "Build %q already exists for service %q\n", serviceManifest.Version, serviceName)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "Build %q required for service %q\n", serviceManifest.Version, serviceName)
					buildTasks = append(buildTasks, map[string]string{
						"service": serviceName,
					})
				}
			}
			buildMatrix := map[string]interface{}{
				"include": buildTasks,
			}
			json, err := json.Marshal(&buildMatrix)
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "::set-output name=buildMatrix::%s\n", json)
			if len(buildTasks) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "::set-output name=buildRequired::1\n")
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "::set-output name=buildRequired::0\n")
			}
		},
	}
}
