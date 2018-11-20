package gogh

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// LocalRepo repository specifier
type LocalRepo struct {
	FullPath  string
	RelPath   string
	PathParts []string
}

// FromFullPath will get a local repository with a full path to the directory
func FromFullPath(ctx Context, fullPath string) (*LocalRepo, error) {
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
		return &LocalRepo{fullPath, filepath.ToSlash(rel), pathParts}, nil
	}

	return nil, fmt.Errorf("no repository found for: %s", fullPath)
}

// FromURL will get a local repository location from remote repository URL
func FromURL(ctx Context, remote *url.URL) (*LocalRepo, error) {
	pathParts := append(
		[]string{remote.Host}, strings.Split(remote.Path, "/")...,
	)
	relPath := strings.TrimSuffix(path.Join(pathParts...), ".git")

	var rep *LocalRepo

	// Find existing repository first
	if err := Walk(ctx, func(repo *LocalRepo) error {
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
	return &LocalRepo{
		path.Join(r, relPath),
		relPath,
		pathParts,
	}, nil
}

// Subpaths returns lists of tail parts of relative path from the root directory (shortest first)
// for example, {"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"} for $root/github.com/kyoh86/gogh.
func (repo *LocalRepo) Subpaths() []string {
	tails := make([]string, len(repo.PathParts))

	for i := range repo.PathParts {
		tails[i] = strings.Join(repo.PathParts[len(repo.PathParts)-(i+1):], "/")
	}

	return tails
}

// NonHostPath will get a relative path from its hostname
func (repo *LocalRepo) NonHostPath() string {
	return strings.Join(repo.PathParts[1:], "/")
}

// IsInPrimaryRoot check which the repository is in primary root directory for gogh
func (repo *LocalRepo) IsInPrimaryRoot(ctx Context) bool {
	return strings.HasPrefix(repo.FullPath, ctx.PrimaryRoot())
}

// Matches checks if any subpath of the repository equals the query.
func (repo *LocalRepo) Matches(pathQuery string) bool {
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

// Walk thorugh local repositories in gogh.root directories
func Walk(ctx Context, callback func(*LocalRepo) error) error {
	for _, root := range ctx.Roots() {
		if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
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
		}); err != nil {
			return err
		}
	}
	return nil
}
