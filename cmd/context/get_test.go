package context

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

type testResult struct {
	err error
	out string
}

func execCommand(cmd *cobra.Command) testResult {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()

	return testResult{
		err: err,
		out: buf.String(),
	}
}

func TestGetCommandValidConfig(t *testing.T) {
	cmd := newContextGetCmd()
	cmd.Flags().String("config", "", "")
	cmd.Flags().Bool("dry-run", false, "")
	cmd.SetArgs([]string{"--config", "../../.g3ops/config.yml"})

	result := execCommand(cmd)

	actual := strings.TrimSpace(result.out)
	if actual != "sandbox" {
		t.Fatalf("expected \"%s\" got \"%s\"", "sandbox", actual)
	}
}

func TestGetCommandMissingConfig(t *testing.T) {
	cmd := newContextGetCmd()
	cmd.Flags().String("config", "", "")
	cmd.Flags().Bool("dry-run", false, "")

	result := execCommand(cmd)

	expected := "No current context found"
	actual := strings.TrimSpace(result.out)
	if actual != expected {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, actual)
	}
}
