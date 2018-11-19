package gogh

import (
	"fmt"
	"strings"
)

// List local repositories
func List(ctx Context, exact, fullpath, short, primary bool, query string) error {
	filter := filterFunc(ctx, exact, primary, query)

	repos := []*LocalRepo{}

	if err := Walk(ctx, func(repo *LocalRepo) error {
		if !filter(repo) {
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

func filterFunc(ctx Context, exact, primary bool, query string) func(*LocalRepo) bool {
	switch {
	case query == "":
		if primary {
			return func(repo *LocalRepo) bool {
				return repo.IsInPrimaryRoot(ctx)
			}
		}
		return func(_ *LocalRepo) bool {
			return true
		}
	case exact:
		if primary {
			return func(repo *LocalRepo) bool {
				return repo.IsInPrimaryRoot(ctx) && repo.Matches(query)
			}
		}
		return func(repo *LocalRepo) bool {
			return repo.Matches(query)
		}
	default:
		if primary {
			return func(repo *LocalRepo) bool {
				return repo.IsInPrimaryRoot(ctx) && strings.Contains(repo.NonHostPath(), query)
			}
		}
		return func(repo *LocalRepo) bool {
			return strings.Contains(repo.NonHostPath(), query)
		}
	}
}

func shortName(dups map[string]bool, repo *LocalRepo) string {
	for _, p := range repo.Subpaths() {
		if dups[p] {
			continue
		}
		return p
	}
	return repo.FullPath
}
