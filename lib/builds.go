package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

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

// G3opsBuildCatalog - represents a build catalog for a service
type G3opsBuildCatalog struct {
	Builds []G3opsBuild
}

func getCatalogFileName(service string) string {
	return fmt.Sprintf(".g3ops/builds/%s.yml", service)
}

const buildsDir = ".g3ops/builds"

// SaveBuild - saves a new build to the catalog for the given service
func SaveBuild(service string, build G3opsBuild) {
	catalog := LoadBuildCatalog(service)
	catalog.Builds = append(catalog.Builds, build)
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
