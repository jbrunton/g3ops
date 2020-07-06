package lib

import (
	"github.com/spf13/afero"
)

// Container - DI container
type Container struct {
	FileSystem *afero.Afero
	Executor   Executor
}

// Copy - creates a copy of the container
func (container *Container) Copy() *Container {
	return &Container{
		FileSystem: container.FileSystem,
		Executor:   container.Executor,
	}
}

// NewContainer - creates a new production container instance. Use NewTestContainer for testing.
func NewContainer() *Container {
	return &Container{
		FileSystem: CreateOsFs(),
		Executor:   &CommandExecutor{},
	}
}
