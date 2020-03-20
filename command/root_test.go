package command_test

import (
	"fmt"
	"strings"

	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/env"
)

func ExampleRoots() {
	yml := strings.NewReader("roots:\n  - /foo\n  - /bar")

	config, err := env.GetAccess(yml, env.EnvarPrefix)
	if err != nil {
		panic(err)
	}
	if err := command.Roots(&config, false); err != nil {
		panic(err)
	}
	fmt.Println()
	if err := command.Roots(&config, true); err != nil {
		panic(err)
	}

	// Output:
	// /foo
	//
	// /foo
	// /bar
}
