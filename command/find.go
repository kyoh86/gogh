package command

import (
	"fmt"

	"github.com/kyoh86/gogh/gogh"
)

// Find a path of a local repository
func Find(ctx gogh.Context, remote *gogh.Remote) error {
	path, err := gogh.FindLocalPath(ctx, remote)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintln(ctx.Stdout(), path); err != nil {
		return err
	}

	return nil
}
