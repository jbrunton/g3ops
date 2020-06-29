package workflow

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jbrunton/g3ops/cmd/styles"
	"github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v2"

	"github.com/jbrunton/g3ops/lib"
	"github.com/spf13/cobra"
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
#@     g3ops commit build ${{ matrix.service }}
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

type workflowLockFile struct {
	Build string
}

func newGenerateWorkflowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generates workflow files",
		Run: func(cmd *cobra.Command, args []string) {
			context, err := lib.GetContext(cmd)
			if err != nil {
				fmt.Println(styles.StyleError(err.Error()))
				os.Exit(1)
			}
			valuesPath := context.Config.Ci.Workflows.Build.Values
			//targetPath := context.Config.Ci.Workflows.Build.Target
			_, err = os.Stat(valuesPath)
			if err != nil {
				fmt.Println(styles.StyleError(err.Error()))
				os.Exit(1)
			}
			lib.ExecCommandI(fmt.Sprintf("ytt -f - -f %s > %s",
				context.Config.Ci.Workflows.Build.Values,
				context.Config.Ci.Workflows.Build.Target), context, template)

			buildWorkflow, err := ioutil.ReadFile(context.Config.Ci.Workflows.Build.Target)
			if err != nil {
				panic(err)
			}
			buildChecksum := fmt.Sprintf("%x", md5.Sum(buildWorkflow))

			lockInfo := workflowLockFile{
				Build: buildChecksum,
			}

			lockFilePath := ".g3ops/workflow-lock.yml"
			lockFileData, err := yaml.Marshal(&lockInfo)
			if err != nil {
				panic(err)
			}

			if context.DryRun {
				fmt.Println(aurora.Yellow(fmt.Sprintf("--dry-run passed, skipping update of %q", lockFilePath)))
			} else {
				err = ioutil.WriteFile(lockFilePath, lockFileData, 0644)
				if err != nil {
					panic(err)
				}
			}
		},
	}
}

// WorkflowCmd represents the context command
var WorkflowCmd = &cobra.Command{
	Use: "workflow",
}

func init() {
	WorkflowCmd.AddCommand(newGenerateWorkflowCmd())
}
