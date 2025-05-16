package filesystem

import (
	"errors"
	"path/filepath"
	"slices"
	"sync"

	"github.com/kyoh86/gogh/v4/core/workspace"
)

// WorkspaceService manages workspace roots stored in the filesystem
type WorkspaceService struct {
	roots       []workspace.Root
	primaryRoot workspace.Root
	changed     bool
	mu          sync.RWMutex
	// You might need a config file path or other storage mechanism
}

// NewWorkspaceService creates a new instance of RootService
func NewWorkspaceService() workspace.WorkspaceService {
	return &WorkspaceService{
		roots: []workspace.Root{},
	}
}

// GetRoots returns all registered workspace roots
func (s *WorkspaceService) GetRoots() []workspace.Root {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]workspace.Root, len(s.roots))
	copy(result, s.roots)
	return result
}

// GetPrimaryRoot returns the primary workspace root
func (s *WorkspaceService) GetPrimaryRoot() workspace.Root {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.primaryRoot
}

// GetLayoutFor returns a Layout for the root
func (s *WorkspaceService) GetLayoutFor(root workspace.Root) workspace.LayoutService {
	return NewLayoutService(root)
}

// GetPrimaryLayout returns a Layout for the primary root
func (s *WorkspaceService) GetPrimaryLayout() workspace.LayoutService {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return NewLayoutService(s.primaryRoot)
}

// SetPrimaryRoot sets the specified path as the primary workspace root
func (s *WorkspaceService) SetPrimaryRoot(path workspace.Root) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	found := slices.Contains(s.roots, path)

	if !found {
		return errors.New("specified path is not registered as a root")
	}

	s.primaryRoot = path
	s.changed = true
	return nil // You would typically persist this change
}

// AddRoot adds a new workspace root
func (s *WorkspaceService) AddRoot(root workspace.Root, asPrimary bool) error {
	absPath, err := filepath.Abs(root)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for duplicates
	if slices.Contains(s.roots, absPath) {
		return errors.New("root already exists")
	}

	s.roots = append(s.roots, absPath)

	// If this is the first root, make it the primary
	if len(s.roots) == 1 || asPrimary {
		s.primaryRoot = absPath
	}

	s.changed = true

	return nil // You would typically persist this change
}

// RemoveRoot removes the specified workspace root
func (s *WorkspaceService) RemoveRoot(path workspace.Root) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, root := range s.roots {
		if root == path {
			// Remove from slice
			s.roots = slices.Delete(s.roots, i, i+1)

			// If we removed the primary root, update it
			if path == s.primaryRoot {
				if len(s.roots) > 0 {
					s.primaryRoot = s.roots[0]
				} else {
					s.primaryRoot = ""
				}
			}
			s.changed = true
			return nil // You would typically persist this change
		}
	}

	return errors.New("root not found")
}

// HasChanges implements workspace.WorkspaceService.
func (s *WorkspaceService) HasChanges() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.changed
}

// MarkSaved implements workspace.WorkspaceService.
func (s *WorkspaceService) MarkSaved() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.changed = false
}

// Ensure RootService implements workspace.RootService
var _ workspace.WorkspaceService = (*WorkspaceService)(nil)
