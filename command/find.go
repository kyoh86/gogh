package command

import (
	"errors"
	"log"
	"os"

	"github.com/kyoh86/gogh/gogh"
)

// Find a local project
func Find(ev gogh.Env, primary bool, spec *gogh.RepoSpec) error {
	log.Printf("info: Finding a repository %s\n", spec)

	finder := gogh.FindProject
	if primary {
		finder = gogh.FindProjectInPrimary
	}

	formatter := gogh.FullPathFormatter()

	project, _, err := finder(ev, spec)
	if err != nil {
		return err
	}
	formatter.Add(project)

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
