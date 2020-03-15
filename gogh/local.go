package gogh

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/saracen/walker"
)

// Project repository specifier
type Project struct {
	FullPath  string
	RelPath   string
	PathParts []string
	Exists    bool
}

var (
	// ErrProjectNotFound is the error will be raised when a project is not found.
	ErrProjectNotFound = errors.New("project not found")
	// ErrProjectAlreadyExists is the error will be raised when a project already exists.
	ErrProjectAlreadyExists = errors.New("project already exists")
)

// FindProject will find a project (local repository) that matches exactly.
func FindProject(ev Env, repo *Repo) (*Project, error) {
	return findProject(ev, repo, Walk)
}

// FindProjectInPrimary will find a project (local repository) that matches exactly.
func FindProjectInPrimary(ev Env, repo *Repo) (*Project, error) {
	return findProject(ev, repo, WalkInPrimary)
}

func findProject(ev Env, repo *Repo, walker Walker) (*Project, error) {
	var project *Project

	// Find existing repository first
	if err := walker(ev, func(p *Project) error {
		if repo.Match(p) {
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

	return nil, ErrProjectNotFound
}

// FindOrNewProject will find a project (local repository) that matches exactly or create new one.
func FindOrNewProject(ev Env, repo *Repo) (*Project, error) {
	return findOrNewProject(ev, repo, Walk)
}

// FindOrNewProjectInPrimary will find a project (local repository) that matches exactly or create new one.
func FindOrNewProjectInPrimary(ev Env, repo *Repo) (*Project, error) {
	return findOrNewProject(ev, repo, WalkInPrimary)
}

func findOrNewProject(ev Env, repo *Repo, walker Walker) (*Project, error) {
	switch p, err := findProject(ev, repo, walker); err {
	case ErrProjectNotFound:
		// No repository found, returning new one
		return NewProject(ev, repo)
	case nil:
		return p, nil
	default:
		return nil, err
	}
}

// NewProject creates a project (local repository)
func NewProject(ev Env, repo *Repo) (*Project, error) {
	relPath := repo.RelPath()
	fullPath := filepath.Join(PrimaryRoot(ev), relPath)
	return &Project{
		FullPath:  fullPath,
		RelPath:   relPath,
		PathParts: []string{repo.host, repo.owner, repo.name},
		Exists:    isVcsDir(fullPath),
	}, nil
}

// FindProjectPath willl get a project (local repository) path from remote repository URL
func FindProjectPath(ev Env, repo *Repo) (string, error) {
	project, err := FindProject(ev, repo)
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
func (p *Project) IsInPrimaryRoot(ev Env) bool {
	return strings.HasPrefix(p.FullPath, PrimaryRoot(ev))
}

func isVcsDir(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}

// WalkFunc is the type of the function called for each repository visited by Walk / WalkInPrimary
type WalkFunc func(*Project) error

// Walker is the type of the function to visit each repository
type Walker func(Env, WalkFunc) error

// walkInPath thorugh projects (local repositories) in a path
func walkInPath(ev Env, root string, callback WalkFunc) error {
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

	return walker.Walk(root, func(path string, fi os.FileInfo) error {
		if !isVcsDir(path) {
			return nil
		}
		p, err := ParseProject(ev, root, path)
		if err != nil {
			return nil
		}
		if err := callback(p); err != nil {
			return err
		}
		return filepath.SkipDir
	})
}

func ParseProject(ev Env, root string, fullPath string) (*Project, error) {
	rel, err := filepath.Rel(root, fullPath)
	if err != nil {
		return nil, err
	}
	pathParts := strings.Split(rel, string(filepath.Separator))
	if len(pathParts) != 3 {
		return nil, errors.New("not supported project path")
	}
	if ev.GithubHost() != pathParts[0] {
		return nil, fmt.Errorf("not supported project host %q", pathParts[0])
	}
	if err := ValidateOwner(pathParts[1]); err != nil {
		return nil, err
	}
	if err := ValidateName(pathParts[2]); err != nil {
		return nil, err
	}
	return &Project{
		FullPath:  fullPath,
		RelPath:   filepath.ToSlash(rel),
		PathParts: pathParts,
		Exists:    true,
	}, nil
}

// WalkInPrimary thorugh projects (local repositories) in the first gogh.root directory
func WalkInPrimary(ev Env, callback WalkFunc) error {
	return walkInPath(ev, PrimaryRoot(ev), callback)
}

// Walk thorugh projects (local repositories) in gogh.root directories
func Walk(ev Env, callback WalkFunc) error {
	for _, root := range ev.Roots() {
		if err := walkInPath(ev, root, callback); err != nil {
			return err
		}
	}
	return nil
}

// Query searches projects (local repositories) with specified walker
func Query(ev Env, query string, walk Walker, callback WalkFunc) error {
	return walk(ev, func(p *Project) error {
		if query != "" && query != p.FullPath && !strings.Contains(p.RelPath, query) {
			return nil
		}

		return callback(p)
	})
}
