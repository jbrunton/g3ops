package lib

import (
	"fmt"
	"testing"

	"github.com/jbrunton/g3ops/cmd/styles"

	"github.com/jbrunton/g3ops/test"
	"github.com/spf13/cobra"
)

func validateFoo(c *cobra.Command, arg string) error {
	if arg != "foo" {
		return fmt.Errorf("Expected %s, got %s", "foo", arg)
	}
	return nil
}

func TestValidateMissingArgs(t *testing.T) {
	cmd := &cobra.Command{
		Use:  "a",
		Args: ValidateArgs([]ArgValidator{validateFoo}),
		Run:  test.EmptyRun,
	}

	result := test.ExecCommand(cmd)

	expectedErr := styles.StyleError("Missing arguments, expected 1 got 0")
	if result.Err.Error() != expectedErr {
		t.Fatalf("expected \"%s\" got \"%s\"", expectedErr, result.Err.Error())
	}
}

func TestValidateTooManyArgs(t *testing.T) {
	cmd := &cobra.Command{
		Use:  "a",
		Args: ValidateArgs([]ArgValidator{validateFoo}),
		Run:  test.EmptyRun,
	}

	cmd.SetArgs([]string{"foo", "bar"})
	result := test.ExecCommand(cmd)

	expectedErr := styles.StyleError("Too many arguments, expected 1 got 2")
	if result.Err.Error() != expectedErr {
		t.Fatalf("expected \"%s\" got \"%s\"", expectedErr, result.Err.Error())
	}
}

func TestValidateInvalidArg(t *testing.T) {
	cmd := &cobra.Command{
		Use:  "a",
		Args: ValidateArgs([]ArgValidator{validateFoo}),
		Run:  test.EmptyRun,
	}

	cmd.SetArgs([]string{"bar"})
	result := test.ExecCommand(cmd)

	expectedErr := `Expected foo, got bar`
	if result.Err.Error() != expectedErr {
		t.Fatalf("expected \"%s\" got \"%s\"", expectedErr, result.Err.Error())
	}
}

func TestValidateValidArg(t *testing.T) {
	cmd := &cobra.Command{
		Use:  "a",
		Args: ValidateArgs([]ArgValidator{validateFoo}),
		Run:  test.EmptyRun,
	}

	cmd.SetArgs([]string{"foo"})
	result := test.ExecCommand(cmd)

	if result.Err != nil {
		t.Fatalf("expected nil got \"%s\"", result.Err.Error())
	}
}
