package lib

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/logrusorgru/aurora"
	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/jbrunton/g3ops/cmd/styles"
)

const template = `
#@ load("@ytt:data", "data")

#@ def commit_step():
#@   user_name = data.values.git.user_name
#@   user_email = data.values.git.user_email
#@   main_branch = data.values.git.main_branch
#@   return """
#@     git config --global user.name '{0}'
#@     git config --global user.email '{1}'
#@     g3ops commit build ${{{{ matrix.service }}}}
#@     git push origin:{2}
#@   """.lstrip().format(user_name, user_email, main_branch)
#@ end

name: #@ data.values.workflow.name

on:
  pull_request:
    branches:
      - #@ data.values.git.main_branch
  push:
    branches:
      - #@ data.values.git.main_branch

env:
  MAIN_BRANCH: #@ data.values.git.main_branch

jobs:

  manifest_check:
    name: sandbox_manifest_check
    runs-on: ubuntu-latest
    if: github.event_name == 'push'
    outputs:
      buildMatrix: ${{ steps.check.outputs.buildMatrix }}
      buildRequired: ${{ steps.check.outputs.buildRequired }}
      
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '^1.14.4'

      - name: install g3ops
        run: go get github.com/jbrunton/g3ops
        
      - name: validate workflows
        run: g3ops workflows check

      - name: check manifest
        id: check
        run: g3ops set-outputs build-matrix

  build:
    name: sandbox_build
    runs-on: ubuntu-latest
    needs: manifest_check
    if: ${{ needs.manifest_check.outputs.buildRequired == true }}
    strategy:
      matrix: ${{ fromJson(needs.manifest_check.outputs.buildMatrix) }}
    env:
      G3OPS_DOCKER_ACCESS_TOKEN: ${{ secrets.G3OPS_DOCKER_ACCESS_TOKEN }}
      G3OPS_DOCKER_USERNAME: ${{ secrets.G3OPS_DOCKER_USERNAME }}

    steps:
      - uses: actions/checkout@v2
        with:
          token: ${{ secrets.G3OPS_ADMIN_ACCESS_TOKEN }}

      - uses: actions/setup-go@v2
        with:
          go-version: '^1.14.4'

      - name: install g3ops
        run: go get github.com/jbrunton/g3ops

      - name: build
        run: g3ops service build ${{ matrix.service }}

      - name: 'commit'
        run: #@ commit_step()
`

// ValidateWorkflows - returns an error if the workflows are out of date
func ValidateWorkflows(context *G3opsContext) error {
	expectedBuildWorkflow := GenerateWorkflow(context)
	buildWorkflowFile := context.Config.Ci.Workflows.Build.Target
	actualBuildWorkflow, err := ioutil.ReadFile(buildWorkflowFile)
	if err != nil {
		return err
	}
	if string(expectedBuildWorkflow) != string(actualBuildWorkflow) {
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(string(expectedBuildWorkflow), string(actualBuildWorkflow), false)
		fmt.Printf("Workflow %q is out of date. Diff:\n%s", buildWorkflowFile, dmp.DiffPrettyText(diffs))
		return errors.New("Workflows are out of date, please run g3ops workflows generate")
	}
	return nil
}

// GenerateWorkflowFile - generates workflow, saves to file and updates workflow-lock.json
func GenerateWorkflowFile(context *G3opsContext) {
	buildWorkflow := GenerateWorkflow(context)
	buildWorkflowPath := context.Config.Ci.Workflows.Build.Target

	if context.DryRun {
		fmt.Println(aurora.Yellow(fmt.Sprintf("--dry-run passed, would have updated file %q:", buildWorkflowPath)))
		fmt.Println(aurora.Yellow(string(buildWorkflow)))
	} else {
		err := ioutil.WriteFile(buildWorkflowPath, buildWorkflow, 0644)
		if err != nil {
			panic(err)
		}
	}
}

// GenerateWorkflow - generate workflow
func GenerateWorkflow(context *G3opsContext) []byte {
	valuesPath := context.Config.Ci.Workflows.Build.Values
	//targetPath := context.Config.Ci.Workflows.Build.Target
	_, err := os.Stat(valuesPath)
	if err != nil {
		fmt.Println(styles.StyleError(err.Error()))
		os.Exit(1)
	}

	process := exec.Command("bash", "-c", fmt.Sprintf("ytt -f - -f %s", valuesPath))
	stdin, err := process.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdin.Write([]byte(template))
	stdin.Close()
	var out bytes.Buffer
	process.Stdout = &out
	process.Stderr = os.Stderr

	err = process.Run()

	if err != nil {
		os.Exit(1)
	}

	return out.Bytes()
}
