package lib

import (
	"os"
	"testing"

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

	context, err := GetContext(fs, cmd)
	assert.NoError(t, err)

	assert.Equal(t, "test-context", context.Config.Name)
	assert.Equal(t, ".g3ops/config.yml", context.ConfigPath)
	assert.Equal(t, ".g3ops", context.Dir)
	assert.Equal(t, ".github", context.GitHubDir)
}

func TestGetContextConfigFlag(t *testing.T) {
	fs := CreateMemFs()
	fs.WriteFile("my-g3ops-context/my-config.yml", []byte("name: test-context"), 0644)
	cmd := testCommand()
	cmd.ParseFlags([]string{"--config", "my-g3ops-context/my-config.yml"})

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

	context, err := GetContext(fs, cmd)
	assert.NoError(t, err)

	assert.Equal(t, "test-context", context.Config.Name)
	assert.Equal(t, "my-g3ops-context/my-config.yml", context.ConfigPath)
	assert.Equal(t, "my-g3ops-context", context.Dir)
	os.Setenv("G3OPS_CONFIG", "")
}

func TestGetContextRelGitHubDir(t *testing.T) {
	fs := CreateMemFs()
	fs.WriteFile(".g3ops/config.yml", []byte("githubDir: ../path/to/.github"), 0644)
	cmd := testCommand()

	context, err := GetContext(fs, cmd)
	assert.NoError(t, err)

	assert.Equal(t, "../path/to/.github", context.GitHubDir)
}

func TestGetContextNestedRelGitHubDir(t *testing.T) {
	fs := CreateMemFs()
	fs.WriteFile("nested-app/.g3ops/config.yml", []byte("githubDir: ../.github"), 0644)
	cmd := testCommand()
	cmd.ParseFlags([]string{"--config", "nested-app/.g3ops/config.yml"})

	context, err := GetContext(fs, cmd)
	assert.NoError(t, err)

	assert.Equal(t, ".github", context.GitHubDir)
}

func TestGetContextAbsGitHubDir(t *testing.T) {
	fs := CreateMemFs()
	fs.WriteFile(".g3ops/config.yml", []byte("githubDir: /path/to/github-dir"), 0644)
	cmd := testCommand()

	context, err := GetContext(fs, cmd)
	assert.NoError(t, err)

	assert.Equal(t, "/path/to/github-dir", context.GitHubDir)
}
