package lib

import (
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/mock"
)

func TestCreateRelease(t *testing.T) {
	_, g3ops := newTestContext()
	g3ops.Config = &G3opsConfig{
		Repo: "my/repo",
	}
	container := NewTestContainer(g3ops)
	container.Executor.(*TestExecutor).On("ExecCommand", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		command := args.Get(0).(string)
		cloneCommand := regexp.MustCompile(`^git clone --depth 1 git@github.com:my/repo.git (\S+)$`)
		matches := cloneCommand.FindStringSubmatch(command)
		if len(matches) > 0 {
			dir := matches[1]
			container.FileSystem.WriteFile(filepath.Join(dir, "manifest.yml"), []byte("version: 1.2.3"), 0644)
			container.FileSystem.WriteFile(filepath.Join(dir, ".g3ops/config.yml"), []byte("repo: my/repo"), 0644)
		}
	})
	container.GitHubService.(*MockGitHubService).On("GetRepository", g3ops).Return(&github.Repository{DefaultBranch: github.String("develop")}, nil)
	container.GitHubService.(*MockGitHubService).On("CreatePullRequest", mock.Anything, g3ops).
		Return(&github.PullRequest{
			HTMLURL: github.String("https://github.com/my/repo/pull/101"),
		}, nil).
		Run(func(args mock.Arguments) {
			newPr := args.Get(0).(*github.NewPullRequest)
			assert.Equal(t, "Update version to 1.2.4", *newPr.Title)
			assert.Contains(t, *newPr.Head, "release-1.2.4")
			assert.Equal(t, "develop", *newPr.Base)
		})

	CreateNewRelease(container.FileSystem, container.Executor, container.GitHubService, g3ops)

	container.GitHubService.(*MockGitHubService).AssertExpectations(t)
}
