package lib

import (
	"fmt"
	"os"
	"strconv"

	"github.com/blang/semver/v4"
	"github.com/jbrunton/g3ops/services"
	"github.com/spf13/afero"
)

// CreateNewRelease - creates a new release
func CreateNewRelease(fs *afero.Afero, executor Executor, gitHubService services.GitHubService, clock Clock, g3ops *G3opsContext) {
	dir, newContext := CloneTempRepo(fs, executor, g3ops)
	defer os.RemoveAll(dir)

	manifest, err := newContext.GetReleaseManifest(fs)
	if err != nil {
		panic(err)
	}

	version, err := semver.Make(manifest.Version)
	if err != nil {
		panic(err)
	}
	fmt.Println("Current version:", version.String())

	version.IncrementPatch()
	fmt.Println("New version:", version.String())
	manifest.Version = version.String()
	err = newContext.SaveReleaseManifest(fs, manifest)
	if err != nil {
		panic(err)
	}

	branchName := fmt.Sprintf("release-%s-%s", version.String(), strconv.Itoa(int(clock.Now().UTC().Unix())))
	commitMessage := fmt.Sprintf("Update version to %s", version.String())
	CommitChanges(dir, []string{"manifest.yml"}, commitMessage, branchName, newContext, executor)

	repo, err := gitHubService.GetRepository(g3ops.RepoID)
	if err != nil {
		panic(err)
	}

	// TODO: only create PR if config.releases.createPullRequest is true
	newPr := &services.NewPullRequest{
		Title: commitMessage,
		Head:  branchName,
		Base:  *repo.DefaultBranch,
	}
	pr, err := gitHubService.CreatePullRequest(newPr, g3ops.RepoID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Created PR for release: %s\n", *pr.HTMLURL)
}
