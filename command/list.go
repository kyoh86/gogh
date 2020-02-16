package command

import (
	"os"

	"github.com/kyoh86/gogh/gogh"
)

// List local projects
func List(ctx gogh.Context, formatter gogh.ProjectListFormatter, primary bool, query string) error {
	var walk gogh.Walker = gogh.Walk
	if primary {
		walk = gogh.WalkInPrimary
	}

	if err := gogh.Query(ctx, query, walk, func(p *gogh.Project) error {
		formatter.Add(p)
		return nil
	}); err != nil {
		return err
	}

	return formatter.PrintAll(os.Stdout, "\n")
}
