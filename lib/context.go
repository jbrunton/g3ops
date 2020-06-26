package lib

import (
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// G3opsContext - type of current g3ops context
type G3opsContext struct {
	Name         string
	Environments map[string]struct {
		Manifest string
	}
	Services map[string]struct {
		Manifest string
	}
	Ci struct {
		Defaults struct {
			Build struct {
				Env     map[string]string
				Command string
			}
		}
	}
}

// G3opsCommandContext - current command context
type G3opsCommandContext struct {
	Context G3opsContext
	DryRun  bool
}

// GetCommandContext - returns the current command context
func GetCommandContext(cmd *cobra.Command) G3opsCommandContext {
	dryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		panic(err)
	}

	ctx, err := LoadContextManifest()
	if err != nil {
		panic(err)
	}

	return G3opsCommandContext{
		Context: ctx,
		DryRun:  dryRun,
	}
}

// LoadContextManifest - finds and returns the G3opsContext
func LoadContextManifest() (G3opsContext, error) {
	data, err := ioutil.ReadFile(".g3ops/config.yml")

	if err != nil {
		return G3opsContext{}, err
	}

	ctx := G3opsContext{}
	err = yaml.Unmarshal(data, &ctx)
	if err != nil {
		panic(err)
	}

	for envName, env := range ctx.Environments {
		path, err := filepath.Abs(env.Manifest)
		if err != nil {
			panic(err)
		}
		env.Manifest = path
		ctx.Environments[envName] = env
	}

	for serviceName, service := range ctx.Services {
		path, err := filepath.Abs(service.Manifest)
		if err != nil {
			panic(err)
		}
		service.Manifest = path
		ctx.Services[serviceName] = service
	}

	return ctx, nil
}

// GetServiceNames - returns the list of services defined in the manifest
// func GetServiceNames() []string {
// 	ctx, err := LoadContextManifest()
// 	if err != nil {
// 		panic(err)
// 	}

// 	var serviceNames []string

// 	for serviceName := range ctx.Services {
// 		serviceNames = append(serviceNames, serviceName)
// 	}

// 	return serviceNames
// }
