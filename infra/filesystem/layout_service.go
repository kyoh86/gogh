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
	root workspace.Root
}

// NewLayoutService creates a new instance of Layout
func NewLayoutService(root workspace.Root) *LayoutService {
	return &LayoutService{root: root}
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

	// Create a reference in the format host/owner/name
	return typ.Ptr(repository.NewReference(parts[0], parts[1], parts[2])), nil
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

	// Create a reference in the format host/owner/name
	return typ.Ptr(repository.NewReference(parts[0], parts[1], parts[2])), nil
}

func (l *LayoutService) PathFor(ref repository.Reference) string {
	return filepath.Join(l.root, ref.Host(), ref.Owner(), ref.Name())
}

func (l *LayoutService) CreateRepositoryFolder(ref repository.Reference) (string, error) {
	path := l.PathFor(ref)
	return path, os.MkdirAll(path, 0755)
}

func (l *LayoutService) DeleteRepository(ref repository.Reference) error {
	path := l.PathFor(ref)
	return os.RemoveAll(path)
}

// Ensure Layout implements workspace.Layout
var _ workspace.LayoutService = (*LayoutService)(nil)
