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
	rts, err := ctx.Roots()
	if err != nil {
		return err
	}
	for _, root := range rts {
		fmt.Println(root)
	}
	return nil
}
