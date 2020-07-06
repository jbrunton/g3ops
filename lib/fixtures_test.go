package lib

import (
	"fmt"

	"github.com/google/go-github/github"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
)

func newTestContext() (*afero.Afero, *G3opsContext) {
	fs := CreateMemFs()
	context := &G3opsContext{
		Dir:       ".g3ops",
		GitHubDir: ".github/",
	}
	return fs, context
}

const invalidTemplate = `
local workflow = {
  on: {
    push: {
      branches: [ "develop" ],
    },
  }
};
std.manifestYamlDoc(workflow)
`

const invalidWorkflow = `# File generated by g3ops, do not modify
# Source: .g3ops/workflows/test.jsonnet
"on":
  "push":
    "branches":
    - "develop"
`

const exampleTemplate = `
local workflow = {
  on: {
    push: {
      branches: [ "develop" ],
    },
  },
	jobs: {
		test: {
			"runs-on": "ubuntu-latest",
			steps: [
			  { run: "echo Hello, World!" }
      ]
    }
  }
};
std.manifestYamlDoc(workflow)
`

const exampleWorkflow = `# File generated by g3ops, do not modify
# Source: .g3ops/workflows/test.jsonnet
"jobs":
  "test":
    "runs-on": "ubuntu-latest"
    "steps":
    - "run": "echo Hello, World!"
"on":
  "push":
    "branches":
    - "develop"
`

func newTestWorkflowDefinition(name string, content string) *WorkflowDefinition {
	return &WorkflowDefinition{
		Name:        name,
		Source:      fmt.Sprintf(".g3ops/workflows/%s.jsonnet", name),
		Destination: fmt.Sprintf(".github/workflows/%s.yml", name),
		Content:     content,
	}
}

type TestExecutor struct {
	mock.Mock
}

func (executor *TestExecutor) ExecCommand(command string, opts ExecOptions) {
	executor.Called(command, opts)
	fmt.Println("Running " + command)
}

func NewTestContainer(g3ops *G3opsContext) Container {
	return Container{
		FileSystem:    CreateMemFs(),
		Executor:      &TestExecutor{},
		GitHubService: &MockGitHubService{},
	}
}

type MockGitHubService struct {
	mock.Mock
}

func (service *MockGitHubService) GetRepository(g3ops *G3opsContext) (*github.Repository, error) {
	args := service.Called(g3ops)
	return args.Get(0).(*github.Repository), args.Error(1)
}

func (service *MockGitHubService) CreatePullRequest(newPr *github.NewPullRequest, g3ops *G3opsContext) (*github.PullRequest, error) {
	args := service.Called(newPr, g3ops)
	return args.Get(0).(*github.PullRequest), args.Error(1)
}
