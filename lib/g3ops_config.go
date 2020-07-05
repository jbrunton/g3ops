package lib

import (
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// G3opsConfig - type of current g3ops context
type G3opsConfig struct {
	Name         string
	Environments map[string]g3opsEnvironmentConfig
	Services     map[string]g3opsServiceConfig
	Ci           g3opsCiConfig
	Workflows    g3opsWorkflowsConfig
}

type g3opsWorkflowsConfig struct {
	GithubDir string `yaml:"githubDir"`
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

var configCache map[*cobra.Command]*G3opsConfig

// GetContextConfig - finds and returns the G3opsConfig
func GetContextConfig(fs *afero.Afero, cmd *cobra.Command, path string) (*G3opsConfig, error) {
	config := configCache[cmd]
	if config != nil {
		return config, nil
	}
	config, err := loadContextConfig(fs, path)
	if err == nil {
		configCache[cmd] = config
	}
	return config, err
}

func loadContextConfig(fs *afero.Afero, path string) (*G3opsConfig, error) {
	data, err := fs.ReadFile(path)

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
	configCache = make(map[*cobra.Command]*G3opsConfig)
}
