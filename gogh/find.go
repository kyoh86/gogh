package gogh

import (
	"fmt"
)

// Find a path of a local repository
func Find(ctx Context, remote *Remote) error {
	path, err := FindLocalPath(ctx, remote)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintln(ctx.Stdout(), path); err != nil {
		return err
	}

	return nil
}
