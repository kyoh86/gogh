package command

import (
	"log"

	"github.com/kyoh86/gogh/gogh"
	"github.com/pkg/errors"
)

// Where is a local project
func Where(ctx gogh.Context, primary bool, exact bool, query string) error {
	log.Printf("info: Finding a repository by query %s", query)

	walk := gogh.Walk
	finder := gogh.FindProject
	if primary {
		walk = gogh.WalkInPrimary
		finder = gogh.FindProjectInPrimary
	}

	formatter := gogh.FullPathFormatter()

	if exact {
		repo, err := gogh.ParseRepo(query)
		if err != nil {
			return err
		}
		project, err := finder(ctx, repo)
		if err != nil {
			return err
		}
		formatter.Add(project)
	} else {
		if err := gogh.Query(ctx, query, walk, func(p *gogh.Project) error {
			formatter.Add(p)
			return nil
		}); err != nil {
			return err
		}
	}

	switch l := formatter.Len(); {
	case l == 1:
		if err := formatter.PrintAll(ctx.Stdout(), "\n"); err != nil {
			return err
		}
	case l < 1:
		log.Println("error: No repository is found")
		return gogh.ErrProjectNotFound
	default:
		log.Println("error: Multiple repositories are found")
		if err := formatter.PrintAll(ctx.Stderr(), "\n"); err != nil {
			return err
		}
		return errors.New("try more precise name")
	}
	return nil
}
