package gogh

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Project repository specifier
type Project struct {
	FullPath  string
	RelPath   string
	PathParts []string
	Exists    bool
}

// FindProject will get a project (local repository) from remote repository URL
func FindProject(ctx Context, remote *Remote) (*Project, error) {
	if err := CheckRemoteHost(ctx, remote); err != nil {
		return nil, err
	}
	relPath := remote.RelPath(ctx)
	var project *Project

	// Find existing repository first
	if err := Walk(ctx, func(p *Project) error {
		if p.RelPath == relPath {
			project = p
			return filepath.SkipDir
		}
		return nil
	}); err != nil {
		return nil, err
	}

	if project != nil {
		return project, nil
	}

	// No repository found, returning new one
	return NewProject(ctx, remote)
}

// NewProject creates a project (local repository)
func NewProject(ctx Context, remote *Remote) (*Project, error) {
	relPath := remote.RelPath(ctx)
	fullPath := filepath.Join(ctx.PrimaryRoot(), relPath)
	info, err := os.Stat(fullPath)
	exists, err := existsProject(fullPath, info, err)
	if err != nil {
		return nil, err
	}
	return &Project{
		FullPath:  fullPath,
		RelPath:   relPath,
		PathParts: []string{remote.Host(ctx), remote.Owner(ctx), remote.Name(ctx)},
		Exists:    exists,
	}, nil
}

func existsProject(path string, info os.FileInfo, err error) (bool, error) {
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

// FindProjectPath willl get a project (local repository) path from remote repository URL
func FindProjectPath(ctx Context, remote *Remote) (string, error) {
	project, err := FindProject(ctx, remote)
	if err != nil {
		return "", err
	}
	return project.FullPath, nil
}

// Subpaths returns lists of tail parts of relative path from the root directory (shortest first)
// for example, {"gogh", "kyoh86/gogh", "github.com/kyoh86/gogh"} for $root/github.com/kyoh86/gogh.
func (p *Project) Subpaths() []string {
	tails := make([]string, len(p.PathParts))

	for i := range p.PathParts {
		tails[i] = strings.Join(p.PathParts[len(p.PathParts)-(i+1):], "/")
	}

	return tails
}

// IsInPrimaryRoot check which the repository is in primary root directory for gogh
func (p *Project) IsInPrimaryRoot(ctx Context) bool {
	return strings.HasPrefix(p.FullPath, ctx.PrimaryRoot())
}

func isVcsDir(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}

// WalkFunc is the type of the function called for each repository visited by Walk / WalkInPrimary
type WalkFunc func(*Project) error

// Walker is the type of the function to visit each repository
type Walker func(Context, WalkFunc) error

// walkInPath thorugh projects (local repositories) in a path
func walkInPath(root string, callback WalkFunc) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		exists, err := existsProject(path, info, err)
		if err != nil {
			return err
		}
		if !exists {
			return nil
		}
		p, err := parseProject(root, path)
		if err != nil {
			return nil
		}
		if err := callback(p); err != nil {
			return err
		}
		return filepath.SkipDir
	})
}

func parseProject(root string, fullPath string) (*Project, error) {
	rel, err := filepath.Rel(root, fullPath)
	if err != nil {
		return nil, err
	}
	pathParts := strings.Split(rel, string(filepath.Separator))
	return &Project{
		FullPath:  fullPath,
		RelPath:   filepath.ToSlash(rel),
		PathParts: pathParts,
		Exists:    true,
	}, nil
}

// WalkInPrimary thorugh projects (local repositories) in the first gogh.root directory
func WalkInPrimary(ctx Context, callback WalkFunc) error {
	return walkInPath(ctx.PrimaryRoot(), callback)
}

// Walk thorugh projects (local repositories) in gogh.root directories
func Walk(ctx Context, callback WalkFunc) error {
	for _, root := range ctx.Roots() {
		if err := walkInPath(root, callback); err != nil {
			return err
		}
	}
	return nil
}

// Query searches projects (local repositories) with specified walker
func Query(ctx Context, query string, walk Walker, callback WalkFunc) error {
	return walk(ctx, func(p *Project) error {
		if query != "" && !strings.Contains(p.RelPath, query) {
			log.Printf("debug: found one repository (%q) but it's not matched for query\n", p.FullPath)
			return nil
		}

		return callback(p)
	})
}
