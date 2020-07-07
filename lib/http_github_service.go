package lib

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/jbrunton/g3ops/services"
	"golang.org/x/oauth2"
)

// HTTPGitHubService - concrete implementation of GitHubService
type HTTPGitHubService struct {
	client *github.Client
}

// GetRepository - returns the repository for the given context
func (service *HTTPGitHubService) GetRepository(repoID services.GitHubRepoID) (*github.Repository, error) {
	repo, _, err := service.client.Repositories.Get(context.Background(), repoID.Owner, repoID.Name)
	return repo, err
}

// CreatePullRequest - creates a pull request in the given repository
func (service *HTTPGitHubService) CreatePullRequest(newPr *services.NewPullRequest, repoID services.GitHubRepoID) (*github.PullRequest, error) {
	pr, _, err := service.client.PullRequests.Create(context.Background(), repoID.Owner, repoID.Name, toArg(newPr))
	return pr, err
}

// ListReleases - list the releases in the repo
func (service *HTTPGitHubService) ListReleases(repoID services.GitHubRepoID) ([]*github.RepositoryRelease, error) {
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

func toArg(newPr *services.NewPullRequest) *github.NewPullRequest {
	return &github.NewPullRequest{
		Title: github.String(newPr.Title),
		Head:  github.String(newPr.Head),
		Base:  github.String(newPr.Base),
	}
}
