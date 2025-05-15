package workspace

import (
	"errors"

	"github.com/kyoh86/gogh/v4/core/store"
)

var (
	// ErrRootNotFound is an error when the root is not found
	ErrRootNotFound = errors.New("repository root not found")

	// ErrRootAlreadyExists is an error when the root already exists
	ErrRootAlreadyExists = errors.New("repository root already exists")

	// ErrNoPrimaryRoot is an error when no primary root is configured
	ErrNoPrimaryRoot = errors.New("no primary repository root configured")
)

type Root = string

// WorkspaceService manages a collection of repository roots
type WorkspaceService interface {
	// GetRoots returns all registered roots
	GetRoots() []Root

	// GetPrimaryRoot returns the primary root
	GetPrimaryRoot() Root

	// GetLayoutFor returns a Layout for the root
	GetLayoutFor(root Root) LayoutService

	// GetPrimaryLayout returns a Layout for the primary root
	GetPrimaryLayout() LayoutService

	// SetPrimaryRoot sets the primary root
	SetPrimaryRoot(Root) error

	// AddRoot adds a new root
	AddRoot(root Root, asPrimary bool) error

	// RemoveRoot removes a root
	RemoveRoot(root Root) error

	store.Content
}
