package lib

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/logrusorgru/aurora"
)

// ExecCommands - execute a sequence of bash commands
// func ExecCommands(commands string, context *G3opsContext) {
// 	funk.ForEach(strings.Split(commands, "\n"), func(command string) {
// 		if strings.TrimSpace(command) != "" {
// 			ExecCommand(command, context)
// 		}
// 	})
// }

// Executor - adapter for executing commands
type Executor interface {
	ExecCommand(command string, opts ExecOptions)
}

// ExecOptions - options for executing a command
type ExecOptions struct {
	Input  string
	Dir    string
	DryRun bool
}

// CommandExecutor - executes commands
type CommandExecutor struct{}

// ExecCommand - executes a command
func (executor *CommandExecutor) ExecCommand(command string, opts ExecOptions) {
	if opts.DryRun {
		fmt.Println(aurora.Yellow("--dry-run passed, skipping command. Would have run:"))
		fmt.Println(aurora.Yellow("  " + command))
		return
	}

	fmt.Println(aurora.Bold("Running"), aurora.Green(command).Bold(), "...")

	process := exec.Command("bash", "-c", command)
	if opts.Input != "" {
		stdin, err := process.StdinPipe()
		if err != nil {
			panic(err)
		}
		stdin.Write([]byte(opts.Input))
		stdin.Close()
	}
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr
	process.Dir = opts.Dir

	err := process.Run()

	if err != nil {
		os.Exit(1)
	}
}

// NewCommandExecutor - returns a concrete CommandExecutor
func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{}
}
