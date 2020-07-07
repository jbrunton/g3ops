package cmd

import (
	"strings"
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
	release := &github.RepositoryRelease{TagName: github.String("1.1.1")}
	gitHubService.On("ListReleases", repoID).Return([]*github.RepositoryRelease{release}, nil)

	// act
	result := test.ExecCommand(cmd)

	// assert
	expectedOutput := strings.Join([]string{
		"Release \"1.1.1\" already exists",
		"::set-output name=releaseRequired::0\n",
	}, "\n")
	assert.Equal(t, expectedOutput, result.Out)
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
	release := &github.RepositoryRelease{TagName: github.String("1.1.1")}
	gitHubService.On("ListReleases", repoID).Return([]*github.RepositoryRelease{release}, nil)

	// act
	result := test.ExecCommand(cmd)

	// assert
	expectedOutput := strings.Join([]string{
		"Release \"1.1.2\" not found, release required",
		"::set-output name=releaseRequired::1",
		"::set-output name=releaseName::1.1.2\n",
	}, "\n")
	assert.Equal(t, expectedOutput, result.Out)
}
