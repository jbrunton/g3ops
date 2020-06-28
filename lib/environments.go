package lib

import (
	"fmt"
)

// G3opsEnvironment - represents an environment manifest
type G3opsEnvironment struct {
	Deploy   g3opsEnvironmentDeployOpts
	Services []G3opsEnvironmentServiceOpts
}

type g3opsEnvironmentDeployOpts struct {
	Host string
}

// G3opsEnvironmentServiceOpts - service details for the manifest
type G3opsEnvironmentServiceOpts struct {
	Name    string
	Version string
}

// FindService - finds service details for the given service
func (env *G3opsEnvironment) FindService(name string) (G3opsEnvironmentServiceOpts, error) {
	for _, service := range env.Services {
		if service.Name == name {
			return service, nil
		}
	}
	return G3opsEnvironmentServiceOpts{}, fmt.Errorf("Could not find service %q", name)
}
