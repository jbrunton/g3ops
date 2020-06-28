package lib

// G3opsEnvironment - represents an environment manifest
type G3opsEnvironment struct {
	Deploy   g3opsEnvironmentDeployOpts
	Services []g3opsEnvironmentServiceOpts
}

type g3opsEnvironmentDeployOpts struct {
	Host string
}

type g3opsEnvironmentServiceOpts struct {
	Name    string
	Version string
}
