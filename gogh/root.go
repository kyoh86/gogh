package gogh

import (
	"fmt"
)

// Root prints a gogh.root
func Root(ctx Context, all bool) error {
	if !all {
		fmt.Println(ctx.PrimaryRoot())
		return nil
	}
	for _, root := range ctx.Roots() {
		fmt.Println(root)
	}
	return nil
}
