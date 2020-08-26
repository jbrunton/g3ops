package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/logrusorgru/aurora"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v2"
)

type G3opsDeployment struct {
	ID        string
	Version   string
	Timestamp time.Time
}

type G3opsDeploymentCatalog struct {
	Deployments []G3opsDeployment
}

func getDeploymentsCatalogFileName(context *G3opsContext, environment string) string {
	return path.Join(context.Dir, fmt.Sprintf("deployments/%s.yml", environment))
}

func createDeployment(version string, context *G3opsContext) (G3opsDeployment, error) {
	deploymentID := uuid.New().String()
	deploymentTimestamp := time.Now().UTC()

	return G3opsDeployment{
		ID:      deploymentID,
		Version: version,
		//BuildSha:  buildSha,
		Timestamp: deploymentTimestamp,
	}, nil
}

// LoadDeploymentsCatalog - loads deployment catalog for the given environment
func LoadDeploymentsCatalog(context *G3opsContext, environment string) G3opsDeploymentCatalog {
	fileName := getDeploymentsCatalogFileName(context, environment)
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return G3opsDeploymentCatalog{}
		}
		panic(err)
	}

	catalog := G3opsDeploymentCatalog{}
	err = yaml.Unmarshal(data, &catalog)
	if err != nil {
		panic(err)
	}

	return catalog
}

func saveDeployment(deployment G3opsDeployment, environment string, context *G3opsContext) {
	catalog := LoadDeploymentsCatalog(context, environment)
	catalog.Deployments = append([]G3opsDeployment{deployment}, catalog.Deployments...)
	fileName := getDeploymentsCatalogFileName(context, environment)
	data, err := yaml.Marshal(&catalog)
	if err != nil {
		panic(err)
	}

	if context.DryRun {
		fmt.Println(aurora.Yellow(fmt.Sprintf("--dry-run passed, skipping update of %q", fileName)))
	} else {
		if _, err := os.Stat(buildsDir); err != nil {
			if os.IsNotExist(err) {
				os.MkdirAll(buildsDir, os.ModePerm)
			} else {
				panic(err)
			}
		}

		err = ioutil.WriteFile(fileName, data, 0644)
		if err != nil {
			panic(err)
		}
	}
}

// GetLatestDeployment - returns the latest deployment (if there is one)
func GetLatestDeployment(environment string, context *G3opsContext) *G3opsDeployment {
	catalog := LoadDeploymentsCatalog(context, environment)
	if len(catalog.Deployments) > 0 {
		return &catalog.Deployments[0]
	}
	return nil
}

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

	deployment, err := createDeployment(version, context)
	if err != nil {
		return err
	}

	saveDeployment(deployment, env, context)

	gitCommand := strings.Join([]string{
		fmt.Sprintf("git add .g3ops/deployments/%s.yml", env),
		fmt.Sprintf(`git commit --allow-empty -m "Deploy v%s to %s"`, version, env),
	}, "\n")
	executor.ExecCommand(gitCommand, opts)

	return nil
}
