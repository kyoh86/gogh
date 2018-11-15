package repo

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kyoh86/gogh/internal/git"
	homedir "github.com/mitchellh/go-homedir"
)

// Local repository specifier
type Local struct {
	FullPath  string
	RelPath   string
	PathParts []string
}

// FromFullPath will get a local repository with a full path to the directory
func FromFullPath(fullPath string) (*Local, error) {
	rts, err := Roots()
	if err != nil {
		return nil, err
	}
	for _, root := range rts {
		if !strings.HasPrefix(fullPath, root) {
			continue
		}

		rel, err := filepath.Rel(root, fullPath)
		if err != nil {
			continue
		}
		pathParts := strings.Split(rel, string(filepath.Separator))
		return &Local{fullPath, filepath.ToSlash(rel), pathParts}, nil
	}

	return nil, fmt.Errorf("no repository found for: %s", fullPath)
}

// FromURL will get a local repository location from remote repository URL
func FromURL(remote *url.URL) (*Local, error) {
	pathParts := append(
		[]string{remote.Host}, strings.Split(remote.Path, "/")...,
	)
	relPath := strings.TrimSuffix(path.Join(pathParts...), ".git")

	var rep *Local

	// Find existing repository first
	if err := Walk(func(repo *Local) error {
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

	r, err := PrimaryRoot()
	if err != nil {
		return nil, err
	}

	// No repository found, returning new one
	return &Local{
		path.Join(r, relPath),
		relPath,
		pathParts,
	}, nil
}

// Subpaths returns lists of tail parts of relative path from the root directory (shortest first)
// for example, {"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"} for $root/github.com/kyoh86/gogh.
func (repo *Local) Subpaths() []string {
	tails := make([]string, len(repo.PathParts))

	for i := range repo.PathParts {
		tails[i] = strings.Join(repo.PathParts[len(repo.PathParts)-(i+1):], "/")
	}

	return tails
}

// NonHostPath will get a relative path from its hostname
func (repo *Local) NonHostPath() string {
	return strings.Join(repo.PathParts[1:], "/")
}

// IsInPrimaryRoot check which the repository is in primary root directory for gogh
func (repo *Local) IsInPrimaryRoot() bool {
	r, err := PrimaryRoot()
	if err != nil {
		return false
	}
	return strings.HasPrefix(repo.FullPath, r)
}

// Matches checks if any subpath of the repository equals the query.
func (repo *Local) Matches(pathQuery string) bool {
	for _, p := range repo.Subpaths() {
		if p == pathQuery {
			return true
		}
	}

	return false
}

func (repo *Local) hasDir(rel string) bool {
	fi, err := os.Stat(filepath.Join(repo.FullPath, rel))
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func isVcsDir(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}

// Walk thorugh local repositories in gogh.root directories
func Walk(callback func(*Local) error) error {
	rts, err := Roots()
	if err != nil {
		return err
	}
	for _, root := range rts {
		if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return nil
			}
			if !isVcsDir(path) {
				return nil
			}
			repo, err := FromFullPath(path)
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

var roots []string

// Roots returns cloned repositories' root directories.
// The root dirs are determined as following:
//
//   - If PM_ROOT environment variable is nonempty, use it as the only root dir.
//   - Otherwise, use the result of `git config --get-all gogh.root` as the dirs.
//   - Otherwise, fallback to the default root, `~/go/src`.
func Roots() ([]string, error) {
	if len(roots) == 0 {
		rts, err := getRoots()
		if err != nil {
			return nil, err
		}
		roots = rts
	}
	return roots, nil
}

func getRoots() ([]string, error) {
	envRoot := os.Getenv("PM_ROOT")
	if envRoot != "" {
		return filepath.SplitList(envRoot), nil
	}
	rts, err := git.GetAllConf("gogh.root")
	if err != nil {
		return nil, err
	}

	if len(rts) == 0 {
		home, err := homedir.Dir()
		if err != nil {
			return nil, err
		}
		rts = []string{filepath.Join(home, "go", "src")}
	}

	for i, v := range rts {
		path := filepath.Clean(v)
		_, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		rts[i], err = filepath.EvalSymlinks(path)
		if err != nil {
			return nil, err
		}
	}

	return rts, nil
}

// PrimaryRoot returns the first one of the root directories to clone repository.
func PrimaryRoot() (string, error) {
	rts, err := Roots()
	if err != nil {
		return "", err
	}
	return rts[0], nil
}
