package service

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/thoas/go-funk"

	"github.com/jbrunton/g3ops/lib"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

type command struct {
	cmd  string
	name string
	args []string
}

func parseCommand(cmd string) command {
	components := strings.Split(cmd, " ")
	return command{cmd, components[0], components[1:len(components)]}
	// var commands []command
	// commands := funk.Map(strings.Split(input, "\n"), func(cmd string) command {
	// 	components := strings.Split(cmd, " ")
	// 	return command{cmd: components[0], args: components[1:len(components)]}
	// })
	// return commands
}

func execCommand(command command) {
	fmt.Println("Running", aurora.Green(command.cmd).Bold(), "...")
	process := exec.Command(command.name, command.args...)

	stdout, err := process.StdoutPipe()
	if err != nil {
		panic(err)
	}
	process.Start()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println("  " + m)
	}
	process.Wait()
}

var buildCmd = &cobra.Command{
	Use:   "build <service>",
	Short: "Build the given service",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires the service name")
		}

		if len(args) > 1 {
			return errors.New("unexpected arguments, only service name expected")
		}

		ctx, err := lib.LoadContextManifest()
		if err != nil {
			panic(err)
		}

		for serviceName := range ctx.Services {
			if serviceName == args[0] {
				return nil
			}
		}

		return errors.New("unknown service: " + args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]

		ctx, err := lib.LoadContextManifest()
		if err != nil {
			panic(err)
		}

		os.Setenv("BUILD_SERVICE", serviceName)
		funk.ForEach(strings.Split(ctx.Ci.Defaults.Build.Command, "\n"), func(cmd string) {
			command := parseCommand(os.ExpandEnv(cmd))
			execCommand(command)
		})
	},
}
