package command

import (
	"log"

	"github.com/kyoh86/gogh/gogh"
	"github.com/pkg/errors"
)

// Where is a local project
func Where(ctx gogh.Context, primary bool, query string) error {
	var walk gogh.Walker = gogh.Walk
	if primary {
		walk = gogh.WalkInPrimary
	}

	formatter := gogh.FullPathFormatter()

	count := 0
	if err := gogh.Query(ctx, query, walk, func(p *gogh.Project) error {
		formatter.Add(p)
		count++
		return nil
	}); err != nil {
		return err
	}

	if count > 1 {
		log.Println("error: Multiple repositories are found")
	}

	if err := formatter.PrintAll(ctx.Stdout(), "\n"); err != nil {
		return err
	}
	if count > 1 {
		return errors.New("try more precise name")
	}
	return nil
}
