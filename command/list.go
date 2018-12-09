package command

import (
	"github.com/kyoh86/gogh/gogh"
)

// List local repositories
func List(ctx gogh.Context, format gogh.RepoListFormat, primary bool, query string) error {
	var walk gogh.Walker = gogh.Walk
	if primary {
		walk = gogh.WalkInPrimary
	}

	formatter, err := format.Formatter()
	if err != nil {
		return err
	}

	if err := gogh.Query(ctx, query, walk, func(l *gogh.Local) error {
		formatter.Add(l)
		return nil
	}); err != nil {
		return err
	}

	return formatter.PrintAll(ctx.Stdout(), "\n")
}
