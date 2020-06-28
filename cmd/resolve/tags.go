package resolve

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/jbrunton/g3ops/lib"

	"github.com/spf13/cobra"
)

var repoRegex, imageRegex *regexp.Regexp

func resolveTags(input []string, context *lib.G3opsContext) {
	repoTags := make(map[string]string)
	for service := range context.Config.Services {
		serviceManifest, err := context.LoadServiceManifest(service)
		if err != nil {
			panic(err)
		}
		// Note: this is for testing. Should be based on environment files.
		build, err := lib.FindBuild(service, serviceManifest.Version)
		if err != nil {
			panic(err) // TODO: error nicely
		}
		repoTags[serviceManifest.Build.Repository] = build.ImageTag
	}
	fmt.Printf("%v\n", repoTags)
	for _, line := range input {
		if imageRegex.MatchString(line) {
			fmt.Println(repoRegex.ReplaceAllStringFunc(line, func(repo string) string {
				tag := repoTags[repo]
				if tag == "" {
					return repo
				}
				return fmt.Sprintf("%s:%s", repo, tag)
			}))
		} else {
			fmt.Println(line)
		}
	}
}

func newResolveTagsCmd() *cobra.Command {
	return &cobra.Command{
		Use: "tags",
		Run: func(cmd *cobra.Command, args []string) {
			info, err := os.Stdin.Stat()
			if err != nil {
				panic(err)
			}
			if info.Mode()&os.ModeNamedPipe == 0 {
				fmt.Println("No input provided. This command expects input from a pipe.")
				return
			}

			reader := bufio.NewReader(os.Stdin)
			var input []string

			for {
				line, _, err := reader.ReadLine()
				if err != nil && err == io.EOF {
					break
				}
				input = append(input, string(line))
			}

			context, err := lib.GetContext(cmd)
			if err != nil {
				panic(err)
			}
			resolveTags(input, context)
		},
	}
}

func init() {
	repoRegex = regexp.MustCompile(`[-\w]+\/[-\w]+`)
	imageRegex = regexp.MustCompile(`^\s*image:\s+"?` + repoRegex.String() + `"?\s*$`)
}
