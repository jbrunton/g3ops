package lib

func ExampleExecCommand() {
	executor := CommandExecutor{}

	executor.ExecCommand("echo Hello, World!", ExecOptions{})

	executor.ExecCommand("echo Hello, World!", ExecOptions{DryRun: true})

	// Output:
	// [1mRunning[0m [1;32mecho Hello, World![0m ...
	// Hello, World!
	// [33m--dry-run passed, skipping command. Would have run:[0m
	// [33m  echo Hello, World![0m
}
