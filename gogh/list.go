package gogh

import (
	"os"
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
		if query != "" || !strings.Contains(repo.NonHostPath(), query) {
			return nil
		}

		formatter.Add(repo)
		return nil
	}); err != nil {
		return err
	}

	return formatter.PrintAll(os.Stdout, "\n")
}
