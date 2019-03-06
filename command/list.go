package command

import (
	"github.com/kyoh86/gogh/gogh"
)

// List local projects
func List(ctx gogh.Context, format gogh.ProjectListFormat, primary bool, query string) error {
	var walk gogh.Walker = gogh.Walk
	if primary {
		walk = gogh.WalkInPrimary
	}

	formatter, err := format.Formatter()
	if err != nil {
		return err
	}

	if err := gogh.Query(ctx, query, walk, func(p *gogh.Project) error {
		formatter.Add(p)
		return nil
	}); err != nil {
		return err
	}

	return formatter.PrintAll(ctx.Stdout(), "\n")
}
