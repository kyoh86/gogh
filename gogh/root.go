package gogh

import (
	"fmt"
)

// Root prints a gogh.root
func Root(ctx Context, all bool) error {
	if !all {
		if _, err := fmt.Fprintln(ctx.Stdout(), ctx.PrimaryRoot()); err != nil {
			return err
		}
		return nil
	}
	for _, root := range ctx.Roots() {
		if _, err := fmt.Fprintln(ctx.Stdout(), root); err != nil {
			return err
		}
	}
	return nil
}
