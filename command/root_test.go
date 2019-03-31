package command_test

import (
	"fmt"

	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/config"
)

func ExampleRoot() {
	ctx := &config.Config{
		VRoot: config.PathListOption{"/foo", "/bar"},
	}

	if err := command.Root(ctx, false); err != nil {
		panic(err)
	}
	fmt.Println()
	if err := command.Root(ctx, true); err != nil {
		panic(err)
	}
	// Output:
	// /foo
	//
	// /foo
	// /bar
}
