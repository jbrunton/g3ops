package services

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
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

// HTTPGitHubService - concrete implementation of GitHubService
type HTTPGitHubService struct {
	client *github.Client
}

// GetRepository - returns the repository for the given context
func (service *HTTPGitHubService) GetRepository(repoID GitHubRepoID) (*github.Repository, error) {
	repo, _, err := service.client.Repositories.Get(context.Background(), repoID.Owner, repoID.Name)
	return repo, err
}

// CreatePullRequest - creates a pull request in the given repository
func (service *HTTPGitHubService) CreatePullRequest(newPr *NewPullRequest, repoID GitHubRepoID) (*github.PullRequest, error) {
	pr, _, err := service.client.PullRequests.Create(context.Background(), repoID.Owner, repoID.Name, toArg(newPr))
	return pr, err
}

// ListReleases - list the releases in the repo
func (service *HTTPGitHubService) ListReleases(repoID GitHubRepoID) ([]*github.RepositoryRelease, error) {
	releases, _, err := service.client.Repositories.ListReleases(context.Background(), repoID.Owner, repoID.Name, nil)
	return releases, err
}

// NewGitHubService - creates a new instance of an HTTPGitHubService
func NewGitHubService() *HTTPGitHubService {
	client := NewGitHubClient()
	return &HTTPGitHubService{client: client}
}

// NewGitHubClient - creates a new client using the GITHUB_TOKEN (if set)
func NewGitHubClient() *github.Client {
	// TODO: this token may not be needed any more
	token := os.Getenv("GITHUB_TOKEN")

	if token == "" {
		//fmt.Println("Warning: no GITHUB_TOKEN set. g3ops won't be able to authenticate, and some functionality won't be supported.")
		return github.NewClient(nil)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

func toArg(newPr *NewPullRequest) *github.NewPullRequest {
	return &github.NewPullRequest{
		Title: github.String(newPr.Title),
		Head:  github.String(newPr.Head),
		Base:  github.String(newPr.Base),
	}
}
