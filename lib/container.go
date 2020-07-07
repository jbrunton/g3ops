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

// Copy - creates a copy of the container
// func (container *Container) Copy() *Container {
// 	return &Container{
// 		FileSystem:    container.FileSystem,
// 		Executor:      container.Executor,
// 		GitHubService: container.GitHubService,
// 	}
// }

// NewContainer - creates a new production container instance. Use NewTestContainer for testing.
func NewContainer() *Container {
	return &Container{
		FileSystem:    CreateOsFs(),
		Executor:      &CommandExecutor{},
		GitHubService: NewGitHubService(),
	}
}
