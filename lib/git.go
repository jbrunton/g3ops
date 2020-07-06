package lib

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CurrentSha - returns the short form version of git rev-parse HEAD
func CurrentSha() string {
	out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(out))
}

// CommitChanges - commits and pushes changes to the filesystem
func CommitChanges(repoDir string, files []string, message string, branchName string, g3ops *G3opsContext) {
	os.Setenv("GIT_DIR", filepath.Join(repoDir, ".git"))
	os.Setenv("GIT_WORK_TREE", repoDir)
	defer os.Setenv("GIT_DIR", "")
	defer os.Setenv("GIT_WORK_TREE", "")

	ExecCommand(fmt.Sprintf("git add %s", strings.Join(files, " ")), g3ops)
	ExecCommand(fmt.Sprintf("git commit -m \"%s\"", message), g3ops)
	ExecCommand(fmt.Sprintf("git push origin HEAD:%s", branchName), g3ops)
	//ExecCommand(fmt.Sprintf("%s push origin $(%s rev-parse --abbrev-ref HEAD)", gitCommand, gitCommand), g3ops)
}
