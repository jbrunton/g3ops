package cmd

import (
	"testing"

	"github.com/jbrunton/g3ops/services"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"

	"github.com/jbrunton/g3ops/lib"
	"github.com/jbrunton/g3ops/test"
)

func TestCheckReleaseManifestExistingVersion(t *testing.T) {
	// arrange
	fs := lib.CreateMemFs()
	gitHubService := test.NewMockGitHubService()
	container := &lib.Container{
		FileSystem:    fs,
		GitHubService: gitHubService,
	}
	repoID := services.GitHubRepoID{Owner: "my", Name: "repo"}
	cmd := newCiCheckCmd(container)
	cmd.Flags().Bool("dry-run", false, "")
	cmd.Flags().String("config", "", "")
	fs.WriteFile(".g3ops/config.yml", []byte("repo: my/repo"), 0644) // TODO: make config file optional?
	cmd.SetArgs([]string{"release-manifest"})

	fs.WriteFile("manifest.yml", []byte("version: 1.1.1"), 0644) // TODO: make config file optional?
	release := &github.RepositoryRelease{Name: github.String("1.1.1")}
	gitHubService.On("ListReleases", repoID).Return([]*github.RepositoryRelease{release}, nil)

	// act
	result := test.ExecCommand(cmd)

	// assert
	assert.Equal(t, "Release \"1.1.1\" already exists\n::set-output name=buildRequired::0\n", result.Out)
}

func TestCheckReleaseManifestNewVersion(t *testing.T) {
	// arrange
	fs := lib.CreateMemFs()
	gitHubService := test.NewMockGitHubService()
	container := &lib.Container{
		FileSystem:    fs,
		GitHubService: gitHubService,
	}
	repoID := services.GitHubRepoID{Owner: "my", Name: "repo"}
	cmd := newCiCheckCmd(container)
	cmd.Flags().Bool("dry-run", false, "")
	cmd.Flags().String("config", "", "")
	fs.WriteFile(".g3ops/config.yml", []byte("repo: my/repo"), 0644) // TODO: make config file optional?
	cmd.SetArgs([]string{"release-manifest"})

	fs.WriteFile("manifest.yml", []byte("version: 1.1.2"), 0644) // TODO: make config file optional?
	release := &github.RepositoryRelease{Name: github.String("1.1.1")}
	gitHubService.On("ListReleases", repoID).Return([]*github.RepositoryRelease{release}, nil)

	// act
	result := test.ExecCommand(cmd)

	// assert
	assert.Equal(t, "Release \"1.1.2\" not found, build required\n::set-output name=buildRequired::1\n", result.Out)
}
