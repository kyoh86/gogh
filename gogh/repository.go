package gogh

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Repository repository specifier
type Repository struct {
	FullPath  string
	RelPath   string
	PathParts []string
}

// FromFullPath will get a local repository with a full path to the directory
func FromFullPath(ctx Context, fullPath string) (*Repository, error) {
	rts := ctx.Roots()
	for _, root := range rts {
		if !strings.HasPrefix(fullPath, root) {
			continue
		}

		rel, err := filepath.Rel(root, fullPath)
		if err != nil {
			continue
		}
		pathParts := strings.Split(rel, string(filepath.Separator))
		return &Repository{fullPath, filepath.ToSlash(rel), pathParts}, nil
	}

	return nil, fmt.Errorf("no repository found for: %s", fullPath)
}

// CheckURL checks url is valid GitHub url
func CheckURL(ctx Context, url *url.URL) error {
	pathParts := strings.Split(strings.TrimRight(url.Path, "/"), "/")
	if len(pathParts) != 3 || len(pathParts[1]) == 0 || len(pathParts[2]) == 0 {
		return errors.New("URL should be formed 'schema://hostname/user/name'")
	}
	if url.Host == "github.com" {
		return nil
	}

	gheHosts := ctx.GHEHosts()

	for _, host := range gheHosts {
		if url.Host == host {
			return nil
		}
	}

	return fmt.Errorf("not supported host: %q", url.Host)
}

// FromURL will get a local repository location from remote repository URL
func FromURL(ctx Context, remote *url.URL) (*Repository, error) {
	if err := CheckURL(ctx, remote); err != nil {
		return nil, err
	}
	pathParts := append(
		[]string{remote.Host}, strings.Split(remote.Path, "/")...,
	)
	relPath := strings.TrimSuffix(path.Join(pathParts...), ".git")
	pathParts = strings.Split(relPath, "/")

	var rep *Repository

	// Find existing repository first
	if err := Walk(ctx, func(repo *Repository) error {
		if repo.RelPath == relPath {
			rep = repo
			return filepath.SkipDir
		}
		return nil
	}); err != nil {
		return nil, err
	}

	if rep != nil {
		return rep, nil
	}

	r := ctx.PrimaryRoot()

	// No repository found, returning new one
	return &Repository{
		path.Join(r, relPath),
		relPath,
		pathParts,
	}, nil
}

// Subpaths returns lists of tail parts of relative path from the root directory (shortest first)
// for example, {"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"} for $root/github.com/kyoh86/gogh.
func (repo *Repository) Subpaths() []string {
	tails := make([]string, len(repo.PathParts))

	for i := range repo.PathParts {
		tails[i] = strings.Join(repo.PathParts[len(repo.PathParts)-(i+1):], "/")
	}

	return tails
}

// NonHostPath will get a relative path from its hostname
func (repo *Repository) NonHostPath() string {
	return strings.Join(repo.PathParts[1:], "/")
}

// IsInPrimaryRoot check which the repository is in primary root directory for gogh
func (repo *Repository) IsInPrimaryRoot(ctx Context) bool {
	return strings.HasPrefix(repo.FullPath, ctx.PrimaryRoot())
}

// Matches checks if any subpath of the repository equals the query.
func (repo *Repository) Matches(pathQuery string) bool {
	for _, p := range repo.Subpaths() {
		if p == pathQuery {
			return true
		}
	}

	return false
}

func isVcsDir(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}

type WalkFunc func(*Repository) error
type Walker func(Context, WalkFunc) error

// walkInPath thorugh local repositories in a path
func walkInPath(ctx Context, path string, callback WalkFunc) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		switch {
		case err == nil:
			// noop
		case os.IsNotExist(err):
			return nil
		default:
			return err
		}
		if !info.IsDir() {
			return nil
		}
		if !isVcsDir(path) {
			return nil
		}
		repo, err := FromFullPath(ctx, path)
		if err != nil {
			return nil
		}
		if err := callback(repo); err != nil {
			return err
		}
		return filepath.SkipDir
	})
}

// WalkInPrimary thorugh local repositories in the first gogh.root directory
func WalkInPrimary(ctx Context, callback WalkFunc) error {
	return walkInPath(ctx, ctx.PrimaryRoot(), callback)
}

// Walk thorugh local repositories in gogh.root directories
func Walk(ctx Context, callback WalkFunc) error {
	for _, root := range ctx.Roots() {
		if err := walkInPath(ctx, root, callback); err != nil {
			return err
		}
	}
	return nil
}
