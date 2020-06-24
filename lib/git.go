package lib

import (
	"os/exec"
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
