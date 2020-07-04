package lib

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// G3opsContext - current command context
type G3opsContext struct {
	Dir        string
	ConfigPath string
	Config     *G3opsConfig
	DryRun     bool
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

	if configPath == "" {
		configPath = os.Getenv("G3OPS_CONFIG")
	}
	if configPath == "" {
		configPath = ".g3ops/config.yml"
	}

	config, err := GetContextConfig(configPath)
	if err != nil {
		return nil, err
	}

	context = &G3opsContext{
		Config:     config,
		DryRun:     dryRun,
		ConfigPath: configPath,
		Dir:        filepath.Dir(configPath),
	}
	return context, nil
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

func init() {
	contextCache = make(map[*cobra.Command]*G3opsContext)
}
