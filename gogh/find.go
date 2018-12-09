package gogh

import (
	"errors"
	"fmt"
	"log"
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
		if _, err := fmt.Fprintln(ctx.Stdout(), path); err != nil {
			return err
		}

	case 1:
		if _, err := fmt.Fprintln(ctx.Stdout(), repos[0].FullPath); err != nil {
			return err
		}

	default:
		log.Println("warn: more than one repositories are found")
		for _, repo := range repos {
			log.Println("warn: - " + strings.Join(repo.PathParts, "/"))
		}
		return errors.New("try more precise name")
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

func getRepoFullPath(ctx Context, nameQuery string) (string, error) {
	name, err := ParseRemoteName(nameQuery)
	if err != nil {
		return "", err
	}
	repo, err := FromURL(ctx, name.URL(ctx, false))
	if err != nil {
		return "", err
	}

	// if the directory exists
	if info, err := os.Stat(repo.FullPath); err == nil && info.IsDir() {
		return repo.FullPath, nil
	}
	return "", errors.New("no repository found")
}
