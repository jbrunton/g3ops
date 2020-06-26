package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora"
)

type command struct {
	cmd  string
	name string
	args []string
}

func parseCommand(cmd string) command {
	components := strings.Split(cmd, " ")
	return command{strings.TrimSpace(cmd), components[0], components[1:len(components)]}
}

func execCommand(command command, dryRun bool) {
	fmt.Println("Running", aurora.Green(command.cmd).Bold(), "...")
	if dryRun {
		fmt.Println(aurora.Yellow("--dry-run passed, skipping command. Would have run:"))
		fmt.Println(aurora.Yellow("  " + command.cmd))
		return
	}

	process := exec.Command(command.name, command.args...)

	stdout, err := process.StdoutPipe()
	if err != nil {
		panic(err)
	}

	var stderr bytes.Buffer
	process.Stderr = &stderr

	process.Start()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println("  " + m)
	}
	err = process.Wait()
	if err != nil {
		fmt.Println(stderr.String())
		panic(err)
	}
}
