package gogh

import (
	"fmt"
	"log"
)

// Root prints a gogh.root
func Root(ctx Context, all bool) error {
	if !all {
		_, err := fmt.Fprintln(ctx.Stdout(), ctx.PrimaryRoot())
		return err
	}
	log.Println("info: finding all roots...")
	for _, root := range ctx.Roots() {
		if _, err := fmt.Fprintln(ctx.Stdout(), root); err != nil {
			return err
		}
	}
	return nil
}
