package lib

// import (
// 	"errors"
// 	"fmt"
// 	"os"
// 	"strconv"

// 	"github.com/blang/semver/v4"
// 	"github.com/jbrunton/g3ops/services"
// 	"github.com/spf13/afero"
// )

// // ReleaseBuilder - struct for creating releases
// type ReleaseBuilder struct {
// 	fileSystem    *afero.Afero
// 	executor      Executor
// 	gitHubService services.GitHubService
// 	clock         Clock
// 	g3ops         *G3opsContext
// }

// // NewReleaseBuilder - constructor for ReleaseBuilder
// func NewReleaseBuilder(container *Container, g3ops *G3opsContext) *ReleaseBuilder {
// 	return &ReleaseBuilder{
// 		fileSystem:    container.FileSystem,
// 		executor:      container.Executor,
// 		gitHubService: container.GitHubService,
// 		clock:         container.Clock,
// 		g3ops:         g3ops,
// 	}
// }

// // CreateNewRelease - creates a new release
// func (builder *ReleaseBuilder) CreateNewRelease(name string, increment string) error {
// 	fmt.Println("name:", name, "increment:", increment)
// 	fs := builder.fileSystem
// 	g3ops := builder.g3ops
// 	gitHubService := builder.gitHubService

// 	dir, newContext := CloneTempRepo(fs, builder.executor, g3ops)
// 	defer os.RemoveAll(dir)

// 	manifest, err := newContext.GetReleaseManifest(fs)
// 	if err != nil {
// 		panic(err)
// 	}

// 	newVersion, err := getNewReleaseVersion(name, increment, manifest.Version)
// 	if err != nil {
// 		return err
// 	}

// 	manifest.Version = newVersion
// 	err = newContext.SaveReleaseManifest(fs, manifest)
// 	if err != nil {
// 		panic(err)
// 	}

// 	commitMessage := fmt.Sprintf("Update version to %s", newVersion)
// 	var branchName string
// 	if g3ops.Config.Releases.CreatePullRequest {
// 		branchName = fmt.Sprintf("release-%s-%s", newVersion, strconv.Itoa(int(builder.clock.Now().UTC().Unix())))
// 	} else {
// 		branchName = CurrentBranch(dir)
// 	}

// 	CommitChanges(dir, []string{"manifest.yml"}, commitMessage, branchName, newContext, builder.executor)

// 	if g3ops.Config.Releases.CreatePullRequest {
// 		repo, err := gitHubService.GetRepository(g3ops.RepoID)
// 		if err != nil {
// 			panic(err)
// 		}

// 		newPr := &services.NewPullRequest{
// 			Title: commitMessage,
// 			Head:  branchName,
// 			Base:  *repo.DefaultBranch,
// 		}
// 		if g3ops.DryRun {
// 			fmt.Printf("--dry-run passed, skipping pull request. Would have created PR:\n%#v\n", newPr)
// 		} else {
// 			pr, err := gitHubService.CreatePullRequest(newPr, g3ops.RepoID)
// 			if err != nil {
// 				fmt.Println(err)
// 				os.Exit(1)
// 			}
// 			fmt.Printf("Created PR for release: %s\n", *pr.HTMLURL)
// 		}
// 	}
// 	return nil
// }

// func getNewReleaseVersion(name string, increment string, currentVersion string) (string, error) {
// 	if name != "" {
// 		version, err := semver.Make(name)
// 		if err != nil {
// 			return "", fmt.Errorf("invalid version name: %q, should be in semver format", name)
// 		}
// 		fmt.Println("Updating to version:", name)
// 		return version.String(), nil
// 	}

// 	if currentVersion == "" {
// 		return "", errors.New("current version isn't set, specify the new version by name")
// 	}
// 	version, err := semver.Make(currentVersion)
// 	if err != nil {
// 		return "", fmt.Errorf("error parsing current version: %q, should be in semvar format", currentVersion)
// 	}
// 	fmt.Println("Current version:", version.String())

// 	switch increment {
// 	case "":
// 		version.IncrementPatch()
// 	case "patch":
// 		version.IncrementPatch()
// 	case "minor":
// 		version.IncrementMinor()
// 	case "major":
// 		version.IncrementMajor()
// 	default:
// 		panic(fmt.Errorf("Unexpected increment type: %q", increment))
// 	}
// 	fmt.Println("New version:", version.String())
// 	return version.String(), nil
// }
