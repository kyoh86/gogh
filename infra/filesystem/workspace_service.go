package filesystem

import (
	"errors"
	"path/filepath"
	"slices"
	"sync"

	"github.com/kyoh86/gogh/v3/core/workspace"
)

// WorkspaceService manages workspace roots stored in the filesystem
type WorkspaceService struct {
	roots       []workspace.Root
	defaultRoot workspace.Root
	mu          sync.RWMutex
	// You might need a config file path or other storage mechanism
}

// NewRootService creates a new instance of RootService
func NewRootService() *WorkspaceService {
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

// GetDefaultRoot returns the default workspace root
func (s *WorkspaceService) GetDefaultRoot() workspace.Root {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.defaultRoot
}

// SetDefaultRoot sets the specified path as the default workspace root
func (s *WorkspaceService) SetDefaultRoot(path workspace.Root) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	found := slices.Contains(s.roots, path)

	if !found {
		return errors.New("specified path is not registered as a root")
	}

	s.defaultRoot = path
	return nil // You would typically persist this change
}

// IsDefault checks if the specified path is the default workspace root
func (s *WorkspaceService) IsDefault(path workspace.Root) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return path == s.defaultRoot
}

// AddRoot adds a new workspace root
func (s *WorkspaceService) AddRoot(root workspace.Root, asDefault bool) error {
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

	// If this is the first root, make it default
	if len(s.roots) == 1 || asDefault {
		s.defaultRoot = absPath
	}

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

			// If we removed the default root, update default
			if path == s.defaultRoot {
				if len(s.roots) > 0 {
					s.defaultRoot = s.roots[0]
				} else {
					s.defaultRoot = ""
				}
			}

			return nil // You would typically persist this change
		}
	}

	return errors.New("root not found")
}

// Ensure RootService implements workspace.RootService
var _ workspace.WorkspaceService = (*WorkspaceService)(nil)
