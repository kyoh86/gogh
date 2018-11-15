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
	repos, err := findRepos(name)
	if err != nil {
		return err
	}

	switch len(repos) {
	case 0:
		path, err := getRepoFullPath(name)
		if err != nil {
			return err
		}
		fmt.Println(path)

	case 1:
		fmt.Println(repos[0].FullPath)

	default:
		var lines []string
		lines = append(lines, "more than one repositories are found")
		lines = append(lines, "try more precise name")
		for _, repo := range repos {
			lines = append(lines, "- "+strings.Join(repo.PathParts, "/"))
		}
		return errors.New(strings.Join(lines, "\n"))
	}
	return nil
}

func findRepos(name string) ([]*repo.Local, error) {
	var repos []*repo.Local
	return repos, repo.Walk(func(repo *repo.Local) error {
		if repo.Matches(name) {
			repos = append(repos, repo)
		}
		return nil
	})
}

func getRepoFullPath(name string) (string, error) {
	spec, err := repo.NewSpec(name)
	if err != nil {
		return "", err
	}
	repo, err := repo.FromURL(spec.URL())
	if err != nil {
		return "", err
	}

	// if the directory exists
	if info, err := os.Stat(repo.FullPath); err == nil && info.IsDir() {
		return repo.FullPath, nil
	}
	return "", errors.New("no repository found")
}
