package gogh

import (
	"fmt"
	"strings"
)

// List local repositories
func List(ctx Context, fullpath, short, primary bool, query string) error {
	var walk Walker = Walk
	if primary {
		walk = WalkInPrimary
	}

	repos := []*Repository{}

	if err := walk(ctx, func(repo *Repository) error {
		if query == "" || strings.Contains(repo.NonHostPath(), query) {
			return nil
		}

		repos = append(repos, repo)
		return nil
	}); err != nil {
		return err
	}

	if short {
		// mark duplicated subpath
		dups := map[string]bool{}
		for _, repo := range repos {
			for _, p := range repo.Subpaths() {
				// (false, not ok) -> (false, ok) -> (true, ok) -> (true, ok) and so on
				_, dups[p] = dups[p]
			}
		}
		for _, repo := range repos {
			fmt.Println(shortName(dups, repo))
		}
	} else {
		for _, repo := range repos {
			if fullpath {
				fmt.Println(repo.FullPath)
			} else {
				fmt.Println(repo.RelPath)
			}
		}
	}
	return nil
}

func shortName(dups map[string]bool, repo *Repository) string {
	for _, p := range repo.Subpaths() {
		if dups[p] {
			continue
		}
		return p
	}
	return repo.FullPath
}
