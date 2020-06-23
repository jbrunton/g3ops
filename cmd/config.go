package cmd

import(
	"fmt"

	"github.com/spf13/viper"
)

// G3opsContext Represents a g3ops context
type G3opsContext struct {
	Name, Url string
}

func getG3opsContexts() []G3opsContext {
	var contexts []G3opsContext
	err := viper.UnmarshalKey("contexts", &contexts)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
		panic(err)
	}
	return contexts
}

func addG3opsContext(context G3opsContext) {
	contexts := getG3opsContexts()
	contexts = append(contexts, context)
	viper.Set("contexts", contexts)
}
