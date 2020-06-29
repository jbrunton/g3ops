package lib

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/thoas/go-funk"
)

// ExecCommands - execute a sequence of bash commands
func ExecCommands(commands string, context *G3opsContext) {
	funk.ForEach(strings.Split(commands, "\n"), func(command string) {
		if strings.TrimSpace(command) != "" {
			ExecCommand(command, context)
		}
	})
}

// ExecCommand - executes a bash command
func ExecCommand(command string, context *G3opsContext) {
	ExecCommandI(command, context, "")
}

// ExecCommandI - execute a bash command
func ExecCommandI(command string, context *G3opsContext, input string) {
	if context.DryRun {
		fmt.Println(aurora.Yellow("--dry-run passed, skipping command. Would have run:"))
		fmt.Println(aurora.Yellow("  " + command))
		return
	}

	fmt.Println("Running", aurora.Green(command).Bold(), "...")

	process := exec.Command("bash", "-c", command)
	if input != "" {
		stdin, err := process.StdinPipe()
		if err != nil {
			panic(err)
		}
		stdin.Write([]byte(input))
		stdin.Close()
	}
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr

	err := process.Run()

	if err != nil {
		os.Exit(1)
	}
}
