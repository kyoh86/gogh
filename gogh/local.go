package gogh

import (
	"errors"
	"fmt"
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
func FindProject(ev Env, spec *RepoSpec) (*Project, *Repo, error) {
	return findProject(ev, spec, Walk)
}

// FindProjectInPrimary will find a project (local repository) that matches exactly.
func FindProjectInPrimary(ev Env, spec *RepoSpec) (*Project, *Repo, error) {
	return findProject(ev, spec, WalkInPrimary)
}

func findProject(ev Env, spec *RepoSpec, walker Walker) (*Project, *Repo, error) {
	var project *Project

	repo, err := spec.Validate(ev)
	if err != nil {
		return nil, nil, err
	}

	// Find existing repository first
	if err := walker(ev, func(p *Project) error {
		if repo.Match(p) {
			project = p
			return filepath.SkipDir
		}
		return nil
	}); err != nil {
		return nil, nil, err
	}

	if project != nil {
		return project, repo, nil
	}

	return nil, repo, ErrProjectNotFound
}

// FindOrNewProject will find a project (local repository) that matches exactly or create new one.
func FindOrNewProject(ev Env, spec *RepoSpec) (*Project, *Repo, error) {
	return findOrNewProject(ev, spec, Walk)
}

// FindOrNewProjectInPrimary will find a project (local repository) that matches exactly or create new one.
func FindOrNewProjectInPrimary(ev Env, spec *RepoSpec) (*Project, *Repo, error) {
	return findOrNewProject(ev, spec, WalkInPrimary)
}

func findOrNewProject(ev Env, spec *RepoSpec, walker Walker) (*Project, *Repo, error) {
	switch proj, repo, err := findProject(ev, spec, walker); err {
	case ErrProjectNotFound:
		// No repository found, returning new one
		return newProject(ev, repo), repo, nil
	case nil:
		return proj, repo, nil
	default:
		return nil, nil, err
	}
}

// NewProject creates a project (local repository)
func NewProject(ev Env, spec *RepoSpec) (*Project, error) {
	repo, err := spec.Validate(ev)
	if err != nil {
		return nil, err
	}
	return newProject(ev, repo), nil
}

func newProject(ev Env, repo *Repo) *Project {
	relPath := repo.RelPath()
	fullPath := filepath.Join(PrimaryRoot(ev), relPath)
	return &Project{
		FullPath:  fullPath,
		RelPath:   relPath,
		PathParts: []string{repo.host, repo.owner, repo.name},
		Exists:    isVcsDir(fullPath),
	}
}

// FindProjectPath willl get a project (local repository) path from remote repository URL
func FindProjectPath(ev Env, spec *RepoSpec) (string, error) {
	project, _, err := FindProject(ev, spec)
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
		return nil
	default:
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("root %q is not a directory", root)
	}

	return walker.Walk(root, func(path string, fi os.FileInfo) error {
		p, err := ParseProject(ev, root, path)
		switch {
		case err == ErrTooShallowPath:
			return nil
		case err == nil:
			// continue
		case err == ErrTooDeepPath:
			return filepath.SkipDir
		default:
			switch err.(type) {
			case *ErrorUnsupportedHost, ErrorInvalidName, ErrorInvalidOwner:
				return filepath.SkipDir
			default:
				return nil
			}
		}

		if !isVcsDir(path) {
			return nil
		}

		if err := callback(p); err != nil {
			return err
		}
		return filepath.SkipDir
	})
}

var (
	ErrTooShallowPath = errors.New("too shallow path")
	ErrTooDeepPath    = errors.New("too deep path")
)

func ParseProject(ev Env, root string, fullPath string) (*Project, error) {
	rel, err := filepath.Rel(root, fullPath)
	if err != nil {
		return nil, err
	}
	pathParts := strings.Split(rel, string(filepath.Separator))
	switch len(pathParts) {
	case 3:
		// continue
	case 0, 1, 2:
		return nil, ErrTooShallowPath
	default:
		return nil, ErrTooDeepPath
	}
	if err := ValidateHost(ev, pathParts[0]); err != nil {
		return nil, err
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
