package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/google/uuid"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
)

// G3opsBuild - represents information about a build
type G3opsBuild struct {
	ID        string    // e.g. 0c8bf7ef-2291-4dba-9e8e-f3d01093fd86
	Version   string    // e.g. 0.2.22
	BuildSha  string    `yaml:"buildSha"` // git build sha, e.g. cc87c1c
	ImageTag  string    `yaml:"imageTag"` // specified by user, but could be based on version + id, e.g. 0.2.22-0c8bf7ef-2291-4dba-9e8e-f3d01093fd86
	Timestamp time.Time // e.g. '2020-06-21T13:43:29.694Z'
}

// FormatTimestamp - human readable string
func (b G3opsBuild) FormatTimestamp() string {
	return b.Timestamp.Format(time.RFC822)
}

// G3opsBuildCatalog - represents a build catalog for a service
type G3opsBuildCatalog struct {
	Builds []G3opsBuild
}

func getCatalogFileName(service string) string {
	return fmt.Sprintf(".g3ops/builds/%s.yml", service)
}

const buildsDir = ".g3ops/builds"

// Build - creates a build for the service and updates the catalog
func Build(service string, version string, cmdCtx *G3opsCommandContext) {
	build, err := createBuild(service, version)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	envMap := map[string]string{
		"BUILD_SERVICE":        service,
		"BUILD_VERSION":        build.Version,
		"BUILD_SHA":            build.BuildSha,
		"BUILD_ID":             build.ID,
		"BUILD_TIMESTAMP":      build.FormatTimestamp(),
		"BUILD_TIMESTAMP_UNIX": string(build.Timestamp.Unix()),
	}

	fmt.Println("Configuring environment for build:")

	funk.ForEach(envMap, func(envvar string, envval string) {
		os.Setenv(envvar, envval)
	})
	funk.ForEach(cmdCtx.Config.Ci.Defaults.Build.Env, func(envvar string, envtemplate string) {
		envval := os.ExpandEnv(envtemplate)
		envMap[envvar] = envval
		os.Setenv(envvar, envval)
	})
	funk.ForEach(envMap, func(envvar string, envval string) {
		fmt.Printf("  %s=%s\n", envvar, envval)
	})

	tag := os.Getenv("TAG")
	if tag == "" {
		panic("TAG must be set")
	}
	build.ImageTag = tag

	funk.ForEach(strings.Split(cmdCtx.Config.Ci.Defaults.Build.Command, "\n"), func(cmd string) {
		command := parseCommand(os.ExpandEnv(cmd))
		if command.cmd != "" {
			execCommand(command, cmdCtx.DryRun)
		}
	})

	saveBuild(service, build)
}

func createBuild(service string, version string) (G3opsBuild, error) {
	err := validateVersion(service, version)
	if err != nil {
		return G3opsBuild{}, err
	}

	buildVersion := version
	buildID := uuid.New().String()
	buildSha := CurrentSha()
	buildTimestamp := time.Now().UTC()

	return G3opsBuild{
		ID:        buildID,
		Version:   buildVersion,
		BuildSha:  buildSha,
		Timestamp: buildTimestamp,
	}, nil
}

// SaveBuild - saves a new build to the catalog for the given service
func saveBuild(service string, build G3opsBuild) {
	catalog := LoadBuildCatalog(service)
	catalog.Builds = append([]G3opsBuild{build}, catalog.Builds...)
	fileName := getCatalogFileName(service)
	data, err := yaml.Marshal(&catalog)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(buildsDir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(buildsDir, os.ModePerm)
		} else {
			panic(err)
		}
	}
	// TODO: check builds dir exists
	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		panic(err)
	}
}

// LoadBuildCatalog - loads a build catalog for the given service
func LoadBuildCatalog(service string) G3opsBuildCatalog {
	fileName := getCatalogFileName(service)
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return G3opsBuildCatalog{}
		}
		panic(err)
	}

	catalog := G3opsBuildCatalog{}
	err = yaml.Unmarshal(data, &catalog)
	if err != nil {
		panic(err)
	}

	return catalog
}

func validateVersion(service string, version string) error {
	catalog := LoadBuildCatalog(service)

	_, err := semver.NewVersion(version)
	if err != nil {
		return fmt.Errorf("Invalid version name %q, must be a semantic version", version)
	}

	for _, build := range catalog.Builds {
		if build.Version == version {
			return fmt.Errorf("Build already exists for version %q", version)
		}
	}

	return nil
}
