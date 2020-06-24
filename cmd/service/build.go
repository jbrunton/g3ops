package service

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jbrunton/g3ops/cmd/context"
	"github.com/spf13/cobra"
)

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

		ctx, err := context.LoadContextManifest()
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
		// ctx, err := context.LoadContextManifest()
		// if err != nil {
		// 	panic(err)
		// }
		serviceName := args[0]

		ctx, err := context.LoadContextManifest()
		if err != nil {
			panic(err)
		}

		//command := ctx.Ci.Defaults.Build.Command
		//commandArgs := ctx.Ci.Defaults.Build.Args
		//process := exec.Command("docker-compose", "build", "$BUILD_SERVICE")
		//process := exec.Command("echo", "docker-compose build $BUILD_SERVICE")
		os.Setenv("BUILD_SERVICE", serviceName)
		//process := exec.Command("echo", os.ExpandEnv("$BUILD_SERVICE"))
		process := exec.Command(
			ctx.Ci.Defaults.Build.Command,
			strings.Split(os.ExpandEnv(ctx.Ci.Defaults.Build.Args), " ")...)

		stdout, err := process.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		process.Start()

		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
		}
		process.Wait()
	},
}
