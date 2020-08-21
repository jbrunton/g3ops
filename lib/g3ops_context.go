package lib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jbrunton/g3ops/services"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// G3opsContext - current command context
type G3opsContext struct {
	Dir        string
	ProjectDir string
	ConfigPath string
	GitHubDir  string
	Config     *G3opsConfig
	DryRun     bool
	RepoID     services.GitHubRepoID
	Manifest   *G3opsManifest
}

var contextCache map[*cobra.Command]*G3opsContext

// NewContext - returns a context for the given config
func NewContext(fs *afero.Afero, configPath string, dryRun bool) (*G3opsContext, error) {
	contextDir := filepath.Dir(configPath)
	projectDir := filepath.Dir(contextDir)

	config, err := GetContextConfig(fs, configPath)
	if err != nil {
		return nil, err
	}

	githubDir := config.GitHubDir
	if githubDir == "" {
		githubDir = ".github/"
	}
	if !filepath.IsAbs(githubDir) {
		githubDir = filepath.Join(filepath.Dir(contextDir), githubDir)
	}

	context := &G3opsContext{
		Config:     config,
		DryRun:     dryRun,
		ConfigPath: configPath,
		GitHubDir:  githubDir,
		Dir:        contextDir,
		ProjectDir: projectDir,
	}

	if config.Repo != "" {
		repoID, err := services.ParseRepoID(config.Repo)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		context.RepoID = repoID
	}
	return context, nil
}

// GetContext - returns the current command context
func GetContext(fs *afero.Afero, cmd *cobra.Command) (*G3opsContext, error) {
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

	return NewContext(fs, configPath, dryRun)
}

// LoadServiceManifest - finds and returns the G3opsService for the given service
// func (context *G3opsContext) LoadServiceManifest(name string) (G3opsService, error) {
// 	serviceContext := context.Config.Services[name]

// 	data, err := ioutil.ReadFile(serviceContext.Manifest)

// 	if err != nil {
// 		return G3opsService{}, err
// 	}

// 	service := G3opsService{}
// 	err = yaml.Unmarshal(data, &service)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return service, nil
// }

func init() {
	contextCache = make(map[*cobra.Command]*G3opsContext)
}

// GetManifest - returns the manifest (if it exists)
func (context *G3opsContext) GetManifest(fs *afero.Afero) (G3opsManifest, error) {
	data, err := fs.ReadFile(filepath.Join(filepath.Dir(context.Dir), "manifest.yml")) // TODO: read config

	if err != nil {
		return G3opsManifest{}, err
	}

	manifest := G3opsManifest{}
	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		panic(err)
	}

	return manifest, nil
}

// SaveReleaseManifest - updates the release manifest
// func (context *G3opsContext) SaveReleaseManifest(fs *afero.Afero, manifest G3opsReleaseManifest) error {
// 	out, err := yaml.Marshal(&manifest)
// 	if err != nil {
// 		return err
// 	}
// 	err = fs.WriteFile(filepath.Join(filepath.Dir(context.Dir), "manifest.yml"), out, 0644) // TODO: read config
// 	return err
// }
