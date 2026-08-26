package filesystem

import (
	"errors"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	corefs "github.com/kyoh86/gogh/v4/core/fs"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// WorkspaceService manages workspace roots stored in the filesystem
type WorkspaceService struct {
	roots           []workspace.Root
	primaryRoot     workspace.Root
	hostPathAliases workspace.HostPathAliases
	changed         bool
	mu              sync.RWMutex
	// You might need a config file path or other storage mechanism
}

// NewWorkspaceService creates a new instance of RootService
func NewWorkspaceService() workspace.WorkspaceService {
	return &WorkspaceService{
		roots: []workspace.Root{},
	}
}

func expandHome(path string) (string, error) {
	return corefs.ReplaceTildeWithHome(path)
}

func normalizedRoot(root workspace.Root) (workspace.Root, error) {
	expanded, err := expandHome(root)
	if err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(expanded)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

func storedRoot(root, normalized workspace.Root) workspace.Root {
	if root == "~" || strings.HasPrefix(root, "~/") {
		return root
	}
	return normalized
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
	s.mu.RLock()
	defer s.mu.RUnlock()
	normalized, err := normalizedRoot(root)
	if err != nil {
		normalized = root
	}
	return NewLayoutServiceWithHostPathAliases(normalized, s.hostPathAliases)
}

// GetPrimaryLayout returns a Layout for the primary root
func (s *WorkspaceService) GetPrimaryLayout() workspace.LayoutService {
	s.mu.RLock()
	defer s.mu.RUnlock()
	normalized, err := normalizedRoot(s.primaryRoot)
	if err != nil {
		normalized = s.primaryRoot
	}
	return NewLayoutServiceWithHostPathAliases(normalized, s.hostPathAliases)
}

// SetHostPathAliases sets host-to-path aliases used by layouts.
func (s *WorkspaceService) SetHostPathAliases(aliases workspace.HostPathAliases) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hostPathAliases = cloneHostPathAliases(aliases)
	s.changed = true
}

// GetHostPathAliases returns host-to-path aliases used by layouts.
func (s *WorkspaceService) GetHostPathAliases() workspace.HostPathAliases {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneHostPathAliases(s.hostPathAliases)
}

// SetPrimaryRoot sets the specified path as the primary workspace root
func (s *WorkspaceService) SetPrimaryRoot(path workspace.Root) error {
	normalizedPath, err := normalizedRoot(path)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, root := range s.roots {
		normalized, err := normalizedRoot(root)
		if err != nil {
			return err
		}
		if normalized == normalizedPath {
			s.primaryRoot = root
			s.changed = true
			return nil
		}
	}

	return errors.New("specified path is not registered as a root")
}

// AddRoot adds a new workspace root
func (s *WorkspaceService) AddRoot(root workspace.Root, asPrimary bool) error {
	normalizedPath, err := normalizedRoot(root)
	if err != nil {
		return err
	}
	storedPath := storedRoot(root, normalizedPath)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for duplicates
	for _, existing := range s.roots {
		normalized, err := normalizedRoot(existing)
		if err != nil {
			return err
		}
		if normalized == normalizedPath {
			return errors.New("root already exists")
		}
	}

	s.roots = append(s.roots, storedPath)

	// If this is the first root, make it the primary
	if len(s.roots) == 1 || asPrimary {
		s.primaryRoot = storedPath
	}

	s.changed = true

	return nil // You would typically persist this change
}

// RemoveRoot removes the specified workspace root
func (s *WorkspaceService) RemoveRoot(path workspace.Root) error {
	normalizedPath, err := normalizedRoot(path)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i, root := range s.roots {
		normalized, err := normalizedRoot(root)
		if err != nil {
			return err
		}
		if normalized == normalizedPath {
			// Remove from slice
			s.roots = slices.Delete(s.roots, i, i+1)

			// If we removed the primary root, update it
			if root == s.primaryRoot {
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
