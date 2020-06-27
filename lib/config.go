package lib

import (
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// G3opsConfig - type of current g3ops context
type G3opsConfig struct {
	Name         string
	Environments map[string]g3opsEnvironmentConfig
	Services     map[string]g3opsServiceConfig
	Ci           g3opsCiConfig
}

type g3opsEnvironmentConfig struct {
	Manifest string
}

type g3opsServiceConfig struct {
	Manifest string
}

type g3opsCiConfig struct {
	Defaults g3opsCiDefaultsConfig
}

type g3opsCiDefaultsConfig struct {
	Build g3opsBuildConfig
}

type g3opsBuildConfig struct {
	Env     map[string]string
	Command string
}

// G3opsCommandContext - current command context
type G3opsCommandContext struct {
	Config G3opsConfig
	DryRun bool
}

// GetCommandContext - returns the current command context
func GetCommandContext(cmd *cobra.Command) G3opsCommandContext {
	dryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		panic(err)
	}

	config, err := LoadContextConfig()
	if err != nil {
		panic(err)
	}

	return G3opsCommandContext{
		Config: config,
		DryRun: dryRun,
	}
}

// LoadContextConfig - finds and returns the G3opsConfig
func LoadContextConfig() (G3opsConfig, error) {
	data, err := ioutil.ReadFile(".g3ops/config.yml")

	if err != nil {
		return G3opsConfig{}, err
	}

	return parseConfig(data)
}

func parseConfig(input []byte) (G3opsConfig, error) {
	config := G3opsConfig{}
	err := yaml.Unmarshal(input, &config)
	if err != nil {
		panic(err)
	}

	for envName, env := range config.Environments {
		path, err := filepath.Abs(env.Manifest)
		if err != nil {
			panic(err)
		}
		env.Manifest = path
		config.Environments[envName] = env
	}

	for serviceName, service := range config.Services {
		path, err := filepath.Abs(service.Manifest)
		if err != nil {
			panic(err)
		}
		service.Manifest = path
		config.Services[serviceName] = service
	}

	return config, nil
}
