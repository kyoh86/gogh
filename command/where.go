package command

import (
	"errors"
	"log"
	"os"

	"github.com/kyoh86/gogh/gogh"
)

// Where is a local project
func Where(ev gogh.Env, primary bool, query string) error {
	log.Printf("info: Finding a repository by query %s\n", query)

	walk := gogh.Walk
	if primary {
		walk = gogh.WalkInPrimary
	}

	formatter := gogh.FullPathFormatter()

	if err := gogh.Query(ev, query, walk, func(p *gogh.Project) error {
		formatter.Add(p)
		return nil
	}); err != nil {
		return err
	}

	switch l := formatter.Len(); {
	case l == 1:
		if err := formatter.PrintAll(os.Stdout, "\n"); err != nil {
			return err
		}
	case l < 1:
		log.Println("error: No repository is found")
		return gogh.ErrProjectNotFound
	default:
		log.Println("error: Multiple repositories are found")
		if err := formatter.PrintAll(os.Stderr, "\n"); err != nil {
			return err
		}
		return errors.New("try more precise name")
	}
	return nil
}
