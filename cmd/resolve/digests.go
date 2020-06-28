package resolve

import (
	"github.com/cppforlife/go-cli-ui/ui"
	kbldcmd "github.com/k14s/kbld/pkg/kbld/cmd"
	"github.com/spf13/cobra"
)

// func resolveDigests(input []string, env string, context *lib.G3opsContext) {
// 	o.FileFlags.Set(cmd)
// 	o.RegistryFlags.Set(cmd)
// }

func newResolveDigestsCmd() *cobra.Command {
	confUI := ui.NewConfUI(ui.NewNoopLogger())
	defer confUI.Flush()
	o := kbldcmd.NewResolveOptions(confUI)
	cmd := &cobra.Command{
		Use:   "digests",
		Short: "Build images and update references",
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.FileFlags.Set(cmd)
	o.RegistryFlags.Set(cmd)
	cmd.Flags().IntVar(&o.BuildConcurrency, "build-concurrency", 4, "Set maximum number of concurrent builds")
	cmd.Flags().BoolVar(&o.ImagesAnnotation, "images-annotation", true, "Annotate resources with images annotation")
	cmd.Flags().StringVar(&o.ImageMapFile, "image-map-file", "", "Set image map file (/cnab/app/relocation-mapping.json in CNAB)")
	cmd.Flags().StringVar(&o.LockOutput, "lock-output", "", "File path to emit configuration with resolved image references")
	return cmd
}
