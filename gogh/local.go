package gogh

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
)

// Project repository specifier
type Project struct {
	FullPath  string
	RelPath   string
	PathParts []string
	Exists    bool
}

var (
	// ProjectNotFound is the error will be raised when a project is not found.
	ProjectNotFound = errors.New("project not found")
	// ProjectNotFound is the error will be raised when a project already exists.
	ProjectAlreadyExists = errors.New("project already exists")
)

// FindProject will find a project (local repository) that matches exactly.
func FindProject(ctx Context, repo *Repo) (*Project, error) {
	if err := CheckRepoHost(ctx, repo); err != nil {
		return nil, err
	}
	relPath := repo.RelPath(ctx)
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

	return nil, ProjectNotFound
}

// FindOrNewProject will find a project (local repository) that matches exactly or create new one.
func FindOrNewProject(ctx Context, repo *Repo) (*Project, error) {
	switch p, err := FindProject(ctx, repo); err {
	case ProjectNotFound:
		// No repository found, returning new one
		return NewProject(ctx, repo)
	case nil:
		return p, nil
	default:
		return nil, err
	}
}

// NewProject creates a project (local repository)
func NewProject(ctx Context, repo *Repo) (*Project, error) {
	if err := CheckRepoHost(ctx, repo); err != nil {
		return nil, err
	}
	relPath := repo.RelPath(ctx)
	fullPath := filepath.Join(ctx.PrimaryRoot(), relPath)
	return &Project{
		FullPath:  fullPath,
		RelPath:   relPath,
		PathParts: []string{repo.Host(ctx), repo.Owner(ctx), repo.Name(ctx)},
		Exists:    isVcsDir(fullPath),
	}, nil
}

// FindProjectPath willl get a project (local repository) path from remote repository URL
func FindProjectPath(ctx Context, repo *Repo) (string, error) {
	project, err := FindProject(ctx, repo)
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
	stat, err := os.Stat(root)
	switch {
	case err == nil:
		// noop
	case os.IsNotExist(err):
		log.Printf("warn: root dir %s is not exist", root)
		return nil
	default:
		return err
	}
	if !stat.IsDir() {
		log.Printf("warn: root dir %s is not a directory", root)
		return nil
	}

	return godirwalk.Walk(root, &godirwalk.Options{
		Callback: func(path string, _ *godirwalk.Dirent) error {
			if !isVcsDir(path) {
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
		},
		FollowSymbolicLinks: true,
		Unsorted:            true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
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
			return nil
		}

		return callback(p)
	})
}
