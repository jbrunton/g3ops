package context

import (
	"strings"
	"testing"

	"github.com/jbrunton/g3ops/test"
)

func TestGetCommandValidConfig(t *testing.T) {
	cmd := newContextGetCmd()
	cmd.Flags().String("config", "", "")
	cmd.Flags().Bool("dry-run", false, "")
	cmd.SetArgs([]string{"--config", "../../.g3ops/config.yml"})

	result := test.ExecCommand(cmd)

	actual := strings.TrimSpace(result.Out)
	if actual != "g3ops" {
		t.Fatalf("expected \"%s\" got \"%s\"", "g3ops", actual)
	}
}

func TestGetCommandMissingConfig(t *testing.T) {
	cmd := newContextGetCmd()
	cmd.Flags().String("config", "", "")
	cmd.Flags().Bool("dry-run", false, "")

	result := test.ExecCommand(cmd)

	expected := "No current context found"
	actual := strings.TrimSpace(result.Out)
	if actual != expected {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, actual)
	}
}
