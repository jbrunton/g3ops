package lib

import (
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/google/go-github/github"
	"github.com/jbrunton/g3ops/services"
	"github.com/jbrunton/g3ops/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateReleasePR(t *testing.T) {
	// arrange
	_, g3ops := newTestContext()
	repoID := services.GitHubRepoID{
		Owner: "my",
		Name:  "repo",
	}
	g3ops.Config = &G3opsConfig{
		Repo:     "my/repo",
		Releases: g3opsReleasesConfig{CreatePullRequest: true},
	}
	g3ops.RepoID = repoID
	container := NewTestContainer(g3ops)
	container.Clock = NewTestClock(time.Unix(123456789, 0))
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
	container.GitHubService.(*test.MockGitHubService).On("GetRepository", repoID).Return(&github.Repository{DefaultBranch: github.String("develop")}, nil)
	expectedPullRequest := services.NewPullRequest{
		Title: "Update version to 1.2.4",
		Head:  "release-1.2.4-123456789",
		Base:  "develop",
	}
	container.GitHubService.(*test.MockGitHubService).On("CreatePullRequest", &expectedPullRequest, repoID).
		Return(&github.PullRequest{
			HTMLURL: github.String("https://github.com/my/repo/pull/101"),
		}, nil)

	// act
	builder := NewReleaseBuilder(container, g3ops)
	builder.CreateNewRelease("", "")

	// assert
	container.GitHubService.(*test.MockGitHubService).AssertExpectations(t)
}

func TestGetNewReleaseVersionIncrements(t *testing.T) {
	var incrementTests = []struct {
		increment string
		expected  string
	}{
		{"", "1.2.4"},
		{"patch", "1.2.4"},
		{"minor", "1.3.0"},
		{"major", "2.0.0"},
	}
	for _, test := range incrementTests {
		version, err := getNewReleaseVersion("", test.increment, "1.2.3")
		assert.NoError(t, err)
		assert.Equal(t, test.expected, version)
	}
}

func TestGetNewReleaseVersionIncrementErrors(t *testing.T) {
	_, err := getNewReleaseVersion("", "minor", "invalid version")
	assert.EqualError(t, err, "error parsing current version: \"invalid version\", should be in semvar format")

	_, err = getNewReleaseVersion("", "minor", "")
	assert.EqualError(t, err, "current version isn't set, specify the new version by name")
}

func TestGetNewReleaseVersionName(t *testing.T) {
	version, err := getNewReleaseVersion("2.0.0", "", "1.2.3")
	assert.NoError(t, err)
	assert.Equal(t, "2.0.0", version)
}

func TestGetNewReleaseVersionNameError(t *testing.T) {
	_, err := getNewReleaseVersion("invalid version", "", "1.2.3")
	assert.EqualError(t, err, "invalid version name: \"invalid version\", should be in semver format")
}
