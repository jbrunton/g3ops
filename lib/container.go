package lib

import (
	"os"

	"github.com/jbrunton/g3ops/services"
	"github.com/jbrunton/gflows/io"
	"github.com/spf13/afero"
)

// Container - DI container
type Container struct {
	FileSystem    *afero.Afero
	Executor      Executor
	GitHubService services.GitHubService
	Clock         Clock
	Logger        *io.Logger
}

// NewContainer - creates a new production container instance. Use NewTestContainer for testing.
func NewContainer() *Container {
	logger := io.NewLogger(os.Stdout, true)
	return &Container{
		FileSystem:    CreateOsFs(),
		Executor:      &CommandExecutor{},
		GitHubService: services.NewGitHubService(),
		Clock:         NewSystemClock(),
		Logger:        logger,
	}
}
