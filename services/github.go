package services

import (
	"fmt"
	"regexp"

	"github.com/google/go-github/github"
)

// GitHubRepoID - unique identifier for a repo
type GitHubRepoID struct {
	Owner string
	Name  string
}

// GitHubService - service for interacting with the GitHub API
type GitHubService interface {
	GetRepository(repo GitHubRepoID) (*github.Repository, error)
	CreatePullRequest(newPr *NewPullRequest, repo GitHubRepoID) (*github.PullRequest, error)
	ListReleases(repo GitHubRepoID) ([]*github.RepositoryRelease, error)
}

// NewPullRequest - represents a new pull request
type NewPullRequest struct {
	Title string
	Head  string
	Base  string
}

// ParseRepoID - parses a string in the form owner/repo
func ParseRepoID(repo string) (GitHubRepoID, error) {
	regex := regexp.MustCompile(`^(\w+)/(\w+)$`)
	matches := regex.FindStringSubmatch(repo)
	if len(matches) == 0 {
		return GitHubRepoID{}, fmt.Errorf("invalid repo name: %q", repo)
	}

	return GitHubRepoID{
		Owner: matches[1],
		Name:  matches[2],
	}, nil
}
