package lib

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// G3opsService - type of current g3ops context
type G3opsService struct {
	Name    string
	Version string
}

// LoadServiceManifest - finds and returns the G3opsService for the given service
func LoadServiceManifest(name string) (G3opsService, error) {
	ctx, err := LoadContextManifest()
	if err != nil {
		panic(err)
	}

	serviceContext := ctx.Services[name]

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
