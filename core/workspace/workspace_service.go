package workspace

import (
	"context"
	"errors"
	"iter"

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
	// TODO: Default -> Primary
	GetDefaultRoot() Root

	// GetLayoutFor returns a Layout for the root
	GetLayoutFor(root Root) LayoutService

	// GetDefaultLayout returns a Layout for the default root
	GetDefaultLayout() LayoutService

	// ListRepository retrieves a list of repositories under the all roots
	ListRepository(ctx context.Context, limit int) iter.Seq2[Repository, error]

	// SetDefaultRoot sets the default root
	SetDefaultRoot(Root) error

	// AddRoot adds a new root
	AddRoot(root Root, asDefault bool) error

	// RemoveRoot removes a root
	RemoveRoot(root Root) error
}

type Repository interface {
	// FullPath returns the full path of the repository
	FullPath() string
	// Path returns the path from root of the repository
	Path() string
	// Host returns the host of the repository
	Host() string
	// Owner returns the owner of the repository
	Owner() string
	// Name returns the name of the repository
	Name() string
}

// WorkspaceStore is a service for saving and loading workspaces
type WorkspaceStore store.Store[WorkspaceService]
