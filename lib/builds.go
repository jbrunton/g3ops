package lib

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// G3opsBuild - represents information about a build
type G3opsBuild struct {
	ID        string // e.g. 0c8bf7ef-2291-4dba-9e8e-f3d01093fd86
	Version   string // e.g. 0.2.22
	BuildSha  string `yaml:"buildSha"` // git build sha, e.g. cc87c1c
	ImageTag  string `yaml:"imageTag"` // specified by user, but could be based on version + id, e.g. 0.2.22-0c8bf7ef-2291-4dba-9e8e-f3d01093fd86
	Timestamp string // e.g. '2020-06-21T13:43:29.694Z'
}

// G3opsBuildCatalog - represents a build catalog for a service
type G3opsBuildCatalog struct {
	Builds []G3opsBuild
}

// SaveBuild - saves a new build to the catalog for the given service
func SaveBuild(service string, build G3opsBuild) {

}

// LoadBuildCatalog - loads a build catalog for the given service
func LoadBuildCatalog(service string) G3opsBuildCatalog {

	data, err := ioutil.ReadFile(fmt.Sprintf(".g3ops/builds/%s.yml", service))
	if err != nil {
		panic(err)
	}

	catalog := G3opsBuildCatalog{}
	err = yaml.Unmarshal(data, &catalog)
	if err != nil {
		panic(err)
	}

	return catalog
}
