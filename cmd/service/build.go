package service

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/thoas/go-funk"

	"github.com/jbrunton/cobra"
	"github.com/jbrunton/g3ops/cmd/styles"
	"github.com/jbrunton/g3ops/lib"
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

var buildCmd = &cobra.Command{
	Use:   "build <service>",
	Short: "Build the given service",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New(styles.StyleError("Argument <service> required"))
		}

		if len(args) > 1 {
			return errors.New(styles.StyleError("Unexpected arguments, only <service> expected"))
		}

		ctx, err := lib.LoadContextManifest()
		if err != nil {
			panic(err)
		}

		var serviceNames []string

		for serviceName := range ctx.Services {
			if serviceName == args[0] {
				return nil
			}
			serviceNames = append(serviceNames, serviceName)
		}

		return errors.New(styles.StyleError(`Unknown service "` + args[0] + `". Valid options: ` + styles.StyleEnumOptions(serviceNames) + "."))
	},
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := args[0]
		var dryRun bool
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			panic(err)
		}

		ctx, err := lib.LoadContextManifest()
		if err != nil {
			panic(err)
		}

		serviceManifest, err := lib.LoadServiceManifest(serviceName)
		if err != nil {
			panic(err)
		}

		build, err := lib.CreateBuild(serviceName, serviceManifest.Version)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		envMap := map[string]string{
			"BUILD_SERVICE":        serviceName,
			"BUILD_VERSION":        build.Version,
			"BUILD_SHA":            build.BuildSha,
			"BUILD_ID":             build.ID,
			"BUILD_TIMESTAMP":      build.FormatTimestamp(),
			"BUILD_TIMESTAMP_UNIX": string(build.Timestamp.Unix()),
		}

		fmt.Println("Configuring environment for build:")

		funk.ForEach(envMap, func(envvar string, envval string) {
			os.Setenv(envvar, envval)
		})
		funk.ForEach(ctx.Ci.Defaults.Build.Env, func(envvar string, envtemplate string) {
			envval := os.ExpandEnv(envtemplate)
			envMap[envvar] = envval
			os.Setenv(envvar, envval)
		})
		funk.ForEach(envMap, func(envvar string, envval string) {
			fmt.Printf("  %s=%s\n", envvar, envval)
		})

		tag := os.Getenv("TAG")
		if tag == "" {
			panic("TAG must be set")
		}
		build.ImageTag = tag

		funk.ForEach(strings.Split(ctx.Ci.Defaults.Build.Command, "\n"), func(cmd string) {
			command := parseCommand(os.ExpandEnv(cmd))
			if command.cmd != "" {
				execCommand(command, dryRun)
			}
		})

		lib.SaveBuild(serviceName, build)
	},
}
