package lib

import (
	"github.com/jbrunton/g3ops/services"
	"github.com/spf13/afero"
)

// Container - DI container
type Container struct {
	FileSystem    *afero.Afero
	Executor      Executor
	GitHubService services.GitHubService
	Clock         Clock
}

// NewContainer - creates a new production container instance. Use NewTestContainer for testing.
func NewContainer() *Container {
	return &Container{
		FileSystem:    CreateOsFs(),
		Executor:      &CommandExecutor{},
		GitHubService: services.NewGitHubService(),
		Clock:         NewSystemClock(),
	}
}
