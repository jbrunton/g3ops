package lib

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// G3opsConfig - type of current g3ops context
type G3opsConfig struct {
	Name         string
	Environments map[string]g3opsEnvironmentConfig
	Services     map[string]g3opsServiceConfig
	Ci           g3opsCiConfig
}

type g3opsWorkflowsConfig struct {
	Build g3opsWorkflowConfig
}

type g3opsWorkflowConfig struct {
	Values string
	Target string
}

type g3opsEnvironmentConfig struct {
	Manifest string
}

type g3opsServiceConfig struct {
	Manifest string
}

type g3opsCiConfig struct {
	Defaults  g3opsCiDefaultsConfig
	Workflows g3opsWorkflowsConfig
}

type g3opsCiDefaultsConfig struct {
	Build g3opsBuildConfig
}

type g3opsBuildConfig struct {
	Env     map[string]string
	Command string
}

var configCache map[string]*G3opsConfig

// GetContextConfig - finds and returns the G3opsConfig
func GetContextConfig(path string) (*G3opsConfig, error) {
	config := configCache[path]
	if config != nil {
		return config, nil
	}
	config, err := loadContextConfig(path)
	if err == nil {
		configCache[path] = config
	}
	return config, err
}

func loadContextConfig(path string) (*G3opsConfig, error) {
	if path == "" {
		path = ".g3ops/config.yml"
	}
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return parseConfig(data)
}

// G3opsService - type of current g3ops context
type G3opsService struct {
	Name    string
	Version string
}

func parseConfig(input []byte) (*G3opsConfig, error) {
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

	return &config, nil
}

func init() {
	configCache = make(map[string]*G3opsConfig)
}
