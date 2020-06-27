package lib

import (
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// G3opsContext - current command context
type G3opsContext struct {
	Config G3opsConfig
	DryRun bool
}

// GetContext - returns the current command context
func GetContext(cmd *cobra.Command) (G3opsContext, error) {
	dryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		panic(err)
	}

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		panic(err)
	}

	config, err := LoadContextConfig(configPath)
	if err != nil {
		return G3opsContext{}, err
	}

	return G3opsContext{
		Config: config,
		DryRun: dryRun,
	}, nil
}

// LoadServiceManifest - finds and returns the G3opsService for the given service
func (context *G3opsContext) LoadServiceManifest(name string) (G3opsService, error) {
	serviceContext := context.Config.Services[name]

	data, err := ioutil.ReadFile(serviceContext.Manifest)

	if err != nil {
		return G3opsService{}, err
	}

	service := G3opsService{}
	err = yaml.Unmarshal(data, &service)
	if err != nil {
		panic(err)
	}

	return service, nil
}
