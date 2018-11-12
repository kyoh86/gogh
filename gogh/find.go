package gogh

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kyoh86/gogh/repo"
)

func Find(name string) error {
	reposFound := []*repo.Local{}
	repo.Walk(func(repo *repo.Local) error {
		if repo.Matches(name) {
			reposFound = append(reposFound, repo)
		}
		return nil
	})

	if len(reposFound) == 0 {
		spec, err := repo.NewSpec(name)

		if err == nil {
			repo, err := repo.FromURL(spec.URL())
			if err != nil {
				return err
			}

			// if the directory exists
			if info, err := os.Stat(repo.FullPath); err == nil && info.IsDir() {
				reposFound = append(reposFound, repo)
			}
		}
	}

	switch len(reposFound) {
	case 0:
		return errors.New("no repository found")

	case 1:
		fmt.Println(reposFound[0].FullPath)

	default:
		var lines []string
		lines = append(lines, "more than one repositories are found")
		lines = append(lines, "try more precise name")
		for _, repo := range reposFound {
			lines = append(lines, "- "+strings.Join(repo.PathParts, "/"))
		}
		return errors.New(strings.Join(lines, "\n"))
	}
	return nil
}
