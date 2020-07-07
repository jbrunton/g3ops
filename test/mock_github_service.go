package test

import (
	"github.com/google/go-github/github"
	"github.com/jbrunton/g3ops/services"
	"github.com/stretchr/testify/mock"
)

// MockGitHubService - a mocked GitHubService for testing
type MockGitHubService struct {
	mock.Mock
}

// NewMockGitHubService - utility function to create a new test instance
func NewMockGitHubService() *MockGitHubService {
	return &MockGitHubService{}
}

// GetRepository - stubbed method
func (service *MockGitHubService) GetRepository(repo services.GitHubRepoID) (*github.Repository, error) {
	args := service.Called(repo)
	return args.Get(0).(*github.Repository), args.Error(1)
}

// CreatePullRequest - stubbed method
func (service *MockGitHubService) CreatePullRequest(newPr *services.NewPullRequest, repo services.GitHubRepoID) (*github.PullRequest, error) {
	args := service.Called(newPr, repo)
	return args.Get(0).(*github.PullRequest), args.Error(1)
}

// ListReleases - stubbed method
func (service *MockGitHubService) ListReleases(repo services.GitHubRepoID) ([]*github.RepositoryRelease, error) {
	args := service.Called(repo)
	return args.Get(0).([]*github.RepositoryRelease), args.Error(1)
}
