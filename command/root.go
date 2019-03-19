package command

import (
	"fmt"
	"log"

	"github.com/kyoh86/gogh/gogh"
)

// Root prints a gogh.root
func Root(ctx gogh.Context, all bool) error {
	if !all {
		_, err := fmt.Fprintln(ctx.Stdout(), gogh.PrimaryRoot(ctx))
		return err
	}
	log.Println("info: Finding all roots...")
	for _, root := range ctx.Roots() {
		if _, err := fmt.Fprintln(ctx.Stdout(), root); err != nil {
			return err
		}
	}
	return nil
}
