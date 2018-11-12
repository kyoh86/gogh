package gogh

import (
	"fmt"
	"strings"

	"github.com/kyoh86/gogh/repo"
)

func List(exact, fullpath, short, primary bool, query string) error {
	var filter func(*repo.Local) bool
	switch {
	case query == "":
		if primary {
			filter = func(repo *repo.Local) bool {
				return repo.IsInPrimaryRoot()
			}
		} else {
			filter = func(_ *repo.Local) bool {
				return true
			}
		}
	case exact:
		if primary {
			filter = func(repo *repo.Local) bool {
				return repo.IsInPrimaryRoot() && repo.Matches(query)
			}
		} else {
			filter = func(repo *repo.Local) bool {
				return repo.Matches(query)
			}
		}
	default:
		if primary {
			filter = func(repo *repo.Local) bool {
				return repo.IsInPrimaryRoot() && strings.Contains(repo.NonHostPath(), query)
			}
		} else {
			filter = func(repo *repo.Local) bool {
				return strings.Contains(repo.NonHostPath(), query)
			}
		}
	}

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

func shortName(dups map[string]bool, repo *repo.Local) string {
	for _, p := range repo.Subpaths() {
		if dups[p] {
			continue
		}
		return p
	}
	return repo.FullPath
}
