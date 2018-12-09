package gogh

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Local repository specifier
type Local struct {
	FullPath  string
	RelPath   string
	PathParts []string
	Exists    bool
}

// FindLocal will get a local repository from remote repository URL
func FindLocal(ctx Context, remote *Remote) (*Local, error) {
	if err := CheckRemoteHost(ctx, remote); err != nil {
		return nil, err
	}
	relPath := remote.RelPath(ctx)
	var local *Local

	// Find existing repository first
	if err := Walk(ctx, func(l *Local) error {
		if l.RelPath == relPath {
			if local != nil {
				return errors.New("more than one repositories are found; try more precise name")
			}
			local = l
			return filepath.SkipDir
		}
		return nil
	}); err != nil {
		return nil, err
	}

	if local != nil {
		return local, nil
	}

	// No repository found, returning new one
	return NewLocal(ctx, remote)
}

// NewLocal creates a local repository
func NewLocal(ctx Context, remote *Remote) (*Local, error) {
	relPath := remote.RelPath(ctx)
	fullPath := filepath.Join(ctx.PrimaryRoot(), relPath)
	info, err := os.Stat(fullPath)
	exists, err := existsLocal(fullPath, info, err)
	if err != nil {
		return nil, err
	}
	return &Local{
		FullPath:  fullPath,
		RelPath:   relPath,
		PathParts: []string{remote.Host(ctx), remote.Owner(ctx), remote.Name(ctx)},
		Exists:    exists,
	}, nil
}

func existsLocal(path string, info os.FileInfo, err error) (bool, error) {
	switch {
	case err == nil:
		// noop
	case os.IsNotExist(err):
		return false, nil
	default:
		return false, err
	}
	return info.IsDir() && isVcsDir(path), nil
}

// FindLocalPath willl get a local repository path from remote repository URL
func FindLocalPath(ctx Context, remote *Remote) (string, error) {
	local, err := FindLocal(ctx, remote)
	if err != nil {
		return "", err
	}
	return local.FullPath, nil
}

// Subpaths returns lists of tail parts of relative path from the root directory (shortest first)
// for example, {"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"} for $root/github.com/kyoh86/gogh.
func (l *Local) Subpaths() []string {
	tails := make([]string, len(l.PathParts))

	for i := range l.PathParts {
		tails[i] = strings.Join(l.PathParts[len(l.PathParts)-(i+1):], "/")
	}

	return tails
}

// IsInPrimaryRoot check which the repository is in primary root directory for gogh
func (l *Local) IsInPrimaryRoot(ctx Context) bool {
	return strings.HasPrefix(l.FullPath, ctx.PrimaryRoot())
}

func isVcsDir(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}

// WalkFunc is the type of the function called for each repository visited by Walk / WalkInPrimary
type WalkFunc func(*Local) error

// Walker is the type of the function to visit each repository
type Walker func(Context, WalkFunc) error

// walkInPath thorugh local repositories in a path
func walkInPath(ctx Context, root string, callback WalkFunc) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		exists, err := existsLocal(path, info, err)
		if err != nil {
			return err
		}
		if !exists {
			return nil
		}
		l, err := parseLocal(ctx, root, path)
		if err != nil {
			return nil
		}
		if err := callback(l); err != nil {
			return err
		}
		return filepath.SkipDir
	})
}

func parseLocal(ctx Context, root string, fullPath string) (*Local, error) {
	rel, err := filepath.Rel(root, fullPath)
	if err != nil {
		return nil, err
	}
	pathParts := strings.Split(rel, string(filepath.Separator))
	return &Local{
		FullPath:  fullPath,
		RelPath:   filepath.ToSlash(rel),
		PathParts: pathParts,
		Exists:    true,
	}, nil
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
