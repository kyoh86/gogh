package gogh

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Find a path of a local repository
func Find(ctx Context, name string) error {
	repos, err := findRepos(ctx, name)
	if err != nil {
		return err
	}

	switch len(repos) {
	case 0:
		path, err := getRepoFullPath(ctx, name)
		if err != nil {
			return err
		}
		fmt.Fprintln(ctx.Stdout(), path)

	case 1:
		fmt.Fprintln(ctx.Stdout(), repos[0].FullPath)

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

func findRepos(ctx Context, name string) ([]*Repository, error) {
	var repos []*Repository
	return repos, Walk(ctx, func(repo *Repository) error {
		if repo.Matches(name) {
			repos = append(repos, repo)
		}
		return nil
	})
}

func getRepoFullPath(ctx Context, name string) (string, error) {
	spec, err := NewSpec(name)
	if err != nil {
		return "", err
	}
	repo, err := FromURL(ctx, spec.URL(ctx, false))
	if err != nil {
		return "", err
	}

	// if the directory exists
	if info, err := os.Stat(repo.FullPath); err == nil && info.IsDir() {
		return repo.FullPath, nil
	}
	return "", errors.New("no repository found")
}
