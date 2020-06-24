package lib

import "os/exec"

// CurrentSha - returns the short form version of git rev-parse HEAD
func CurrentSha() string {
	out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		panic(err)
	}
	return string(out)
}
