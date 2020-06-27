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

	config := G3opsConfig{}
	err = yaml.Unmarshal(data, &config)
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
