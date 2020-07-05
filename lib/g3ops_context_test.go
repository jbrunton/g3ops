package lib

import (
	"os"
	"testing"

	"github.com/jbrunton/g3ops/test"
	"github.com/stretchr/testify/assert"

	"github.com/spf13/cobra"
)

func testCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "cmd",
	}
	cmd.Flags().Bool("dry-run", false, "Preview commands before executing")
	cmd.Flags().String("config", "", "Location of g3ops context config")
	return cmd
}

func TestGetContextDefaults(t *testing.T) {
	fs := CreateMemFs()
	fs.WriteFile(".g3ops/config.yml", []byte("name: test-context"), 0644)
	cmd := testCommand()
	//test.ExecCommand(cmd) // required in order to parse flags

	context, err := GetContext(fs, cmd)
	assert.NoError(t, err)

	assert.Equal(t, "test-context", context.Config.Name)
	assert.Equal(t, ".g3ops/config.yml", context.ConfigPath)
	assert.Equal(t, ".g3ops", context.Dir)
	assert.Equal(t, ".github", context.GithubDir)
}

func TestGetContextConfigFlag(t *testing.T) {
	fs := CreateMemFs()
	fs.WriteFile("my-g3ops-context/my-config.yml", []byte("name: test-context"), 0644)
	cmd := testCommand()
	cmd.SetArgs([]string{"--config", "my-g3ops-context/my-config.yml"})
	cmd.Execute() // required in order to parse flags

	context, err := GetContext(fs, cmd)
	assert.NoError(t, err)

	assert.Equal(t, "test-context", context.Config.Name)
	assert.Equal(t, "my-g3ops-context/my-config.yml", context.ConfigPath)
	assert.Equal(t, "my-g3ops-context", context.Dir)
}

func TestGetContextConfigEnvArg(t *testing.T) {
	fs := CreateMemFs()
	fs.WriteFile("my-g3ops-context/my-config.yml", []byte("name: test-context"), 0644)
	cmd := testCommand()
	os.Setenv("G3OPS_CONFIG", "my-g3ops-context/my-config.yml")
	test.ExecCommand(cmd) // required in order to parse flags

	context, err := GetContext(fs, cmd)
	assert.NoError(t, err)

	assert.Equal(t, "test-context", context.Config.Name)
	assert.Equal(t, "my-g3ops-context/my-config.yml", context.ConfigPath)
	assert.Equal(t, "my-g3ops-context", context.Dir)
	os.Setenv("G3OPS_CONFIG", "")
}

func TestGetContextRelGithubDir(t *testing.T) {
	fs := CreateMemFs()
	fs.WriteFile(".g3ops/config.yml", []byte("workflows:\n  githubDir: ../path/to/.github"), 0644)
	cmd := testCommand()

	context, err := GetContext(fs, cmd)
	assert.NoError(t, err)

	assert.Equal(t, "../path/to/.github", context.GithubDir)
}

func TestGetContextNestedRelGithubDir(t *testing.T) {
	fs := CreateMemFs()
	fs.WriteFile("nested-app/.g3ops/config.yml", []byte("workflows:\n  githubDir: ../.github"), 0644)
	cmd := testCommand()
	cmd.SetArgs([]string{"--config", "nested-app/.g3ops/config.yml"})
	cmd.Execute()
	//os.Setenv("G3OPS_CONFIG", "nested-app/.g3ops/config.yml")

	context, err := GetContext(fs, cmd)
	assert.NoError(t, err)

	assert.Equal(t, ".github", context.GithubDir)
}

func TestGetContextAbsGithubDir(t *testing.T) {
	fs := CreateMemFs()
	fs.WriteFile(".g3ops/config.yml", []byte("workflows:\n  githubDir: /path/to/github-dir"), 0644)
	cmd := testCommand()

	context, err := GetContext(fs, cmd)
	assert.NoError(t, err)

	assert.Equal(t, "/path/to/github-dir", context.GithubDir)
}
