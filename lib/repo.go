package lib

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// CloneTempRepo - clones a shallow copy into a temp directory. Returns the directory path and the context for the new repo.
func CloneTempRepo(fs *afero.Afero, executor Executor, g3ops *G3opsContext) (string, *G3opsContext) {
	dir, err := fs.TempDir("", strings.Join([]string{"g3ops", g3ops.RepoName, "*"}, "-"))
	if err != nil {
		log.Fatal(err)
	}

	executor.ExecCommand(fmt.Sprintf("git clone --depth 1 git@github.com:%s.git %s", g3ops.Config.Repo, dir), ExecOptions{DryRun: g3ops.DryRun})

	// TODO: specify context, not config, and require .g3ops directory in context
	newContext, err := NewContext(fs, filepath.Join(dir, ".g3ops", "config.yml"), g3ops.DryRun)
	if err != nil {
		panic(err)
	}

	return dir, newContext
}
