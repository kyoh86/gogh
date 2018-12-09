package gogh

import (
	"log"
	"strings"
)

// List local repositories
func List(ctx Context, format RepoListFormat, primary bool, query string) error {
	var walk Walker = Walk
	if primary {
		walk = WalkInPrimary
	}

	formatter, err := format.Formatter()
	if err != nil {
		return err
	}

	if err := walk(ctx, func(repo *Repository) error {
		if query != "" || !strings.Contains(repo.Name().String(), query) {
			log.Printf("debug: found one repository (%q) but it's not matched for query\n", repo.FullPath)
			return nil
		}

		formatter.Add(repo)
		return nil
	}); err != nil {
		return err
	}

	return formatter.PrintAll(ctx.Stdout(), "\n")
}
