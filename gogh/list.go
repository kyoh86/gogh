package gogh

import (
	"fmt"
	"strings"

	"github.com/kyoh86/gogh/repo"
)

// List local repositories
func List(exact, fullpath, short, primary bool, query string) error {
	filter := filterFunc(exact, primary, query)

	repos := []*repo.Local{}

	repo.Walk(func(repo *repo.Local) error {
		if filter(repo) == false {
			return nil
		}

		repos = append(repos, repo)
		return nil
	})

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

func filterFunc(exact, primary bool, query string) func(*repo.Local) bool {
	switch {
	case query == "":
		if primary {
			return func(repo *repo.Local) bool {
				return repo.IsInPrimaryRoot()
			}
		} else {
			return func(_ *repo.Local) bool {
				return true
			}
		}
	case exact:
		if primary {
			return func(repo *repo.Local) bool {
				return repo.IsInPrimaryRoot() && repo.Matches(query)
			}
		} else {
			return func(repo *repo.Local) bool {
				return repo.Matches(query)
			}
		}
	default:
		if primary {
			return func(repo *repo.Local) bool {
				return repo.IsInPrimaryRoot() && strings.Contains(repo.NonHostPath(), query)
			}
		} else {
			return func(repo *repo.Local) bool {
				return strings.Contains(repo.NonHostPath(), query)
			}
		}
	}
}

func shortName(dups map[string]bool, repo *repo.Local) string {
	for _, p := range repo.Subpaths() {
		if dups[p] {
			continue
		}
		return p
	}
	return repo.FullPath
}
