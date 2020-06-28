package test

import (
	"bytes"

	"github.com/spf13/cobra"
)

// TODO: move this to test package

// Result - returns the resulting error and the combined stdout and stderr output of the executed command
type Result struct {
	Err error
	Out string
}

// ExecCommand - runs a command under test and returns execution info in the TestResult return value
func ExecCommand(cmd *cobra.Command) Result {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()

	return Result{
		Err: err,
		Out: buf.String(),
	}
}

// EmptyRun - utility function for null command
func EmptyRun(*cobra.Command, []string) {}
