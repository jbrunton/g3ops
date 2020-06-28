package lib

import (
	"errors"
	"fmt"

	"github.com/jbrunton/g3ops/cmd/styles"
	"github.com/spf13/cobra"
)

// ArgValidator - function which returns an error if the argument is invalid
type ArgValidator func(cmd *cobra.Command, arg string) error

// ValidateArgs - Returns a function which validates all the given arguments
func ValidateArgs(argValidators []ArgValidator) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < len(argValidators) {
			return fmt.Errorf(styles.StyleError("Missing arguments, expected %d got %d"), len(argValidators), len(args))
		}
		if len(args) > len(argValidators) {
			return fmt.Errorf(styles.StyleError("Too many arguments, expected %d got %d"), len(argValidators), len(args))
		}
		for i, argValidator := range argValidators {
			err := argValidator(cmd, args[i])
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// ServiceValidator - validates the name of a service
func ServiceValidator(cmd *cobra.Command, arg string) error {
	context, err := GetContext(cmd)
	if err != nil {
		panic(err)
	}

	var serviceNames []string

	for serviceName := range context.Config.Services {
		if serviceName == arg {
			return nil
		}
		serviceNames = append(serviceNames, serviceName)
	}

	return errors.New(styles.StyleError(`Unknown service "` + arg + `". Valid options: ` + styles.StyleEnumOptions(serviceNames) + "."))
}
