package lib

import (
	"fmt"
	"os"

	"github.com/thoas/go-funk"
)

func Deploy(context *G3opsContext, container *Container, version string, env string) error {
	build := FindBuild(version, context)
	if build == nil {
		return fmt.Errorf("Unable to find build for version %s", version)
	}

	envMap := map[string]string{
		"BUILD_VERSION": version,
		"ENVIRONMENT":   env,
		//"BUILD_SHA":       build.BuildSha,
		"BUILD_ID":        build.ID,
		"BUILD_TIMESTAMP": build.FormatTimestamp(),
		//"BUILD_TIMESTAMP_UNIX": string(build.Timestamp.Unix()),
	}

	fmt.Println("Configuring environment for deployment:")

	funk.ForEach(envMap, func(envvar string, envval string) {
		os.Setenv(envvar, envval)
	})
	// funk.ForEach(context.Config.Build.Env, func(envvar string, envtemplate string) {
	// 	envval := os.ExpandEnv(envtemplate)
	// 	envMap[envvar] = envval
	// 	os.Setenv(envvar, envval)
	// })
	funk.ForEach(envMap, func(envvar string, envval string) {
		fmt.Printf("  %s=%s\n", envvar, envval)
	})

	// tag := os.Getenv("TAG")
	// if tag == "" {
	// 	panic("TAG must be set")
	// }
	//build.ImageTag = tag

	opts := ExecOptions{DryRun: context.DryRun, Dir: context.ProjectDir}

	//fmt.Printf("Running command:\n%s\n", context.Config.Build.Command)
	executor := container.Executor
	executor.ExecCommand(context.Config.Deploy.Command, opts)

	return nil
}
