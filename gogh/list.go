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

	if err := walk(ctx, func(l *Local) error {
		if query != "" || !strings.Contains(l.RelPath, query) {
			log.Printf("debug: found one repository (%q) but it's not matched for query\n", l.FullPath)
			return nil
		}

		formatter.Add(l)
		return nil
	}); err != nil {
		return err
	}

	return formatter.PrintAll(ctx.Stdout(), "\n")
}
