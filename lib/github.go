package lib

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// NewGithubClient - creates a new client using the GITHUB_TOKEN (if set)
func NewGithubClient() *github.Client {
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
