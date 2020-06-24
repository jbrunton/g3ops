package lib

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// PrintYaml - highlights the given code if pygmentize is available
// Inspired by https://github.com/pksunkara/pygments
func PrintYaml(code string) {

	if _, err := exec.LookPath("pygmentize"); err != nil {
		fmt.Println("could not find pygmentize")
		fmt.Println(code)
		return
	}

	cmd := exec.Command("pygmentize", "-fterminal256", "-lyaml", "-O style=monokai")
	cmd.Stdin = strings.NewReader(code)

	var out bytes.Buffer
	cmd.Stdout = &out

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("error running cmd")
		fmt.Println(stderr.String())
		fmt.Println(err)
		fmt.Println(code)
		return
	}

	fmt.Println(out.String())
}
