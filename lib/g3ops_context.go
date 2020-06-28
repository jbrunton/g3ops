package lib

import (
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// G3opsContext - current command context
type G3opsContext struct {
	Config *G3opsConfig
	DryRun bool
}

var contextCache map[*cobra.Command]*G3opsContext

// GetContext - returns the current command context
func GetContext(cmd *cobra.Command) (*G3opsContext, error) {
	context := contextCache[cmd]
	if context != nil {
		return context, nil
	}

	dryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		panic(err)
	}

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		panic(err)
	}

	config, err := GetContextConfig(configPath)
	if err != nil {
		return nil, err
	}

	context = &G3opsContext{
		Config: config,
		DryRun: dryRun,
	}
	return context, nil
}

// LoadServiceManifest - finds and returns the G3opsService for the given service
func (context *G3opsContext) LoadServiceManifest(name string) (G3opsService, error) {
	serviceConfig := context.Config.Services[name]

	data, err := ioutil.ReadFile(serviceConfig.Manifest)

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

// LoadEnvironmentManifest - finds and returns the G3opsEnvironment for the given environment name
func (context *G3opsContext) LoadEnvironmentManifest(name string) (G3opsEnvironment, error) {
	environmentConfig := context.Config.Environments[name]

	data, err := ioutil.ReadFile(environmentConfig.Manifest)

	if err != nil {
		return G3opsEnvironment{}, err
	}

	envManifest := G3opsEnvironment{}
	err = yaml.Unmarshal(data, &envManifest)
	if err != nil {
		panic(err)
	}

	return envManifest, nil
}

func init() {
	contextCache = make(map[*cobra.Command]*G3opsContext)
}
