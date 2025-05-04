package workspace

import (
	"errors"

	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/store"
)

var (
	// ErrRootNotFound is an error when the root is not found
	ErrRootNotFound = errors.New("repository root not found")

	// ErrRootAlreadyExists is an error when the root already exists
	ErrRootAlreadyExists = errors.New("repository root already exists")

	// ErrNoDefaultRoot is an error when no default root is configured
	ErrNoDefaultRoot = errors.New("no default repository root configured")
)

type Root = string

// WorkspaceService manages a collection of repository roots
type WorkspaceService interface {
	// GetRoots returns all registered roots
	GetRoots() []Root

	// GetDefaultRoot returns the default root
	GetDefaultRoot() Root

	// SetDefaultRoot sets the default root
	SetDefaultRoot(Root) error

	// AddRoot adds a new root
	AddRoot(root Root, asDefault bool) error

	// RemoveRoot removes a root
	RemoveRoot(root Root) error

	// CreateRepositoryDirectory creates a new directory for the repository
	CreateRepositoryDirectory(root Root, reference *repository.Reference) (string, error)

	// DeleteRepositoryDirectory deletes the directory for the repository
	DeleteRepositoryDirectory(root Root, reference *repository.Reference) error
}

// WorkspaceStore is a service for saving and loading workspaces
type WorkspaceStore store.Store[WorkspaceService]
