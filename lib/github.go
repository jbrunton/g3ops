package lib

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GitHubService - service for interacting with the GitHub API
type GitHubService interface {
	GetRepository(g3ops *G3opsContext) (*github.Repository, error)
	CreatePullRequest(newPr *NewPullRequest, g3ops *G3opsContext) (*github.PullRequest, error)
	ListReleases(g3ops *G3opsContext) ([]*github.RepositoryRelease, error)
}

// NewPullRequest - represents a new pull request
type NewPullRequest struct {
	Title string
	Head  string
	Base  string
}

func (newPr *NewPullRequest) toArg() *github.NewPullRequest {
	return &github.NewPullRequest{
		Title: github.String(newPr.Title),
		Head:  github.String(newPr.Head),
		Base:  github.String(newPr.Base),
	}
}

// HTTPGitHubService - concrete implementation of GitHubService
type HTTPGitHubService struct {
	client *github.Client
}

// GetRepository - returns the repository for the given context
func (service *HTTPGitHubService) GetRepository(g3ops *G3opsContext) (*github.Repository, error) {
	repo, _, err := service.client.Repositories.Get(context.Background(), g3ops.RepoOwnerName, g3ops.RepoName)
	return repo, err
}

// CreatePullRequest - creates a pull request in the given repository
func (service *HTTPGitHubService) CreatePullRequest(newPr *NewPullRequest, g3ops *G3opsContext) (*github.PullRequest, error) {
	pr, _, err := service.client.PullRequests.Create(context.Background(), g3ops.RepoOwnerName, g3ops.RepoName, newPr.toArg())
	return pr, err
}

// ListReleases - list the releases in the repo
func (service *HTTPGitHubService) ListReleases(g3ops *G3opsContext) ([]*github.RepositoryRelease, error) {
	releases, _, err := service.client.Repositories.ListReleases(context.Background(), g3ops.RepoOwnerName, g3ops.RepoName, nil)
	return releases, err
}

// NewGitHubService - creates a new instance of an HTTPGitHubService
func NewGitHubService() *HTTPGitHubService {
	client := NewGitHubClient()
	return &HTTPGitHubService{client: client}
}

// NewGitHubClient - creates a new client using the GITHUB_TOKEN (if set)
func NewGitHubClient() *github.Client {
	token := os.Getenv("GITHUB_TOKEN")

	if token == "" {
		fmt.Println("Warning: no GITHUB_TOKEN set. g3ops won't be able to authenticate, and some functionality won't be supported.")
		return github.NewClient(nil)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}
