package command

import (
	"fmt"
	"log"

	"github.com/kyoh86/gogh/gogh"
)

// Root prints a gogh.root
func Root(ctx gogh.Context, all bool) error {
	if !all {
		fmt.Fprintln(ctx.Stdout(), ctx.PrimaryRoot())
		return nil
	}
	log.Println("info: Finding all roots...")
	for _, root := range ctx.Root() {
		fmt.Fprintln(ctx.Stdout(), root)
	}
	return nil
}
