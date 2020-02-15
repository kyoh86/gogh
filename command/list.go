package command

import (
	"io"
	"os"

	"github.com/kyoh86/gogh/gogh"
)

// List local projects
func List(ctx gogh.Context, formatter gogh.ProjectListFormatter, primary bool, isPublic bool, query string) error {
	InitLog(ctx)

	var walk gogh.Walker = gogh.Walk
	if primary {
		walk = gogh.WalkInPrimary
	}

	if err := gogh.Query(ctx, query, walk, func(p *gogh.Project) error {
		if isPublic {
			repo, err := gogh.ParseProject(p)
			if err != nil {
				return err
			}
			pub, err := repo.IsPublic(ctx)
			if err != nil {
				return err
			}
			if !pub {
				return nil
			}
		}
		formatter.Add(p)
		return nil
	}); err != nil {
		return err
	}

	var stdout io.Writer = os.Stdout
	if ctx, ok := ctx.(gogh.IOContext); ok {
		stdout = ctx.Stdout()
	}
	return formatter.PrintAll(stdout, "\n")
}
