package gogh

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/kyoh86/gogh/repo"
)

// Find a path of a local repository
func Find(name string) error {
	reposFound := []*repo.Local{}
	repo.Walk(func(repo *repo.Local) error {
		if repo.Matches(name) {
			reposFound = append(reposFound, repo)
		}
		return nil
	})

	switch len(reposFound) {
	case 0:
		spec, err := repo.NewSpec(name)
		if err != nil {
			return err
		}
		repo, err := repo.FromURL(spec.URL())
		if err != nil {
			return err
		}

		// if the directory exists
		if info, err := os.Stat(repo.FullPath); err == nil && info.IsDir() {
			fmt.Println(repo.FullPath)
			return nil
		}
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
