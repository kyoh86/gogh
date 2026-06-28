package filesystem

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
	"github.com/kyoh86/gogh/v4/typ"
)

// LayoutService is a filesystem-based standard repository layout implementation
type LayoutService struct {
	root            workspace.Root
	hostPathAliases workspace.HostPathAliases
	pathHostAliases map[string]string
}

// NewLayoutService creates a new instance of LayoutService
func NewLayoutService(root workspace.Root) *LayoutService {
	return NewLayoutServiceWithHostPathAliases(root, nil)
}

// NewLayoutServiceWithHostPathAliases creates a new LayoutService with host-to-path aliases.
func NewLayoutServiceWithHostPathAliases(root workspace.Root, aliases workspace.HostPathAliases) *LayoutService {
	hostPathAliases := cloneHostPathAliases(aliases)
	return &LayoutService{
		root:            root,
		hostPathAliases: hostPathAliases,
		pathHostAliases: reverseHostPathAliases(hostPathAliases),
	}
}

// GetRoot returns the root of the layout
func (l *LayoutService) GetRoot() string {
	return l.root
}

// Match returns the reference corresponding to the given path
func (l *LayoutService) Match(path string) (*repository.Reference, error) {
	// Get the relative path from the root
	relPath, err := filepath.Rel(l.root, path)
	if err != nil {
		return nil, workspace.ErrNotMatched
	}

	// Split the path components
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	if len(parts) < 3 {
		return nil, workspace.ErrNotMatched
	}
	if parts[0] == ".." || parts[0] == "." {
		return nil, workspace.ErrNotMatched
	}

	// Create a reference in the format host/owner/name.
	return typ.Ptr(repository.NewReference(l.hostFromPathName(parts[0]), parts[1], parts[2])), nil
}

// ExactMatch returns the reference corresponding exactly to the given path
func (l *LayoutService) ExactMatch(path string) (*repository.Reference, error) {
	// Get the relative path from the root
	relPath, err := filepath.Rel(l.root, path)
	if err != nil {
		return nil, workspace.ErrNotMatched
	}

	// Split the path components
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	if len(parts) != 3 {
		return nil, workspace.ErrNotMatched
	}

	// Create a reference in the format host/owner/name.
	return typ.Ptr(repository.NewReference(l.hostFromPathName(parts[0]), parts[1], parts[2])), nil
}

func (l *LayoutService) PathFor(ref repository.Reference) string {
	return filepath.Join(l.root, l.hostPathName(ref.Host()), ref.Owner(), ref.Name())
}

func (l *LayoutService) legacyPathFor(ref repository.Reference) string {
	return filepath.Join(l.root, ref.Host(), ref.Owner(), ref.Name())
}

func (l *LayoutService) PathsFor(ref repository.Reference) []string {
	path := l.PathFor(ref)
	legacyPath := l.legacyPathFor(ref)
	if legacyPath == path {
		return []string{path}
	}
	return []string{path, legacyPath}
}

func (l *LayoutService) hostPathName(host string) string {
	if alias, ok := l.hostPathAliases[host]; ok {
		return alias
	}
	return host
}

func (l *LayoutService) hostFromPathName(name string) string {
	if host, ok := l.pathHostAliases[name]; ok {
		return host
	}
	return name
}

func (l *LayoutService) CreateRepositoryFolder(ref repository.Reference) (string, error) {
	path := l.PathFor(ref)
	return path, os.MkdirAll(path, 0o755)
}

func (l *LayoutService) DeleteRepository(ref repository.Reference) error {
	for _, path := range l.PathsFor(ref) {
		if _, err := os.Stat(path); err == nil {
			return os.RemoveAll(path)
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	return os.RemoveAll(l.PathFor(ref))
}

// Ensure LayoutService implements workspace.LayoutService
var _ workspace.LayoutService = (*LayoutService)(nil)
