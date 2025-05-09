package workspace

import (
	"errors"

	"github.com/kyoh86/gogh/v3/core/repository"
)

var ErrNotMatched = errors.New("repository layout not matched")

// LayoutService defines the layout structure of a repository under a root
type LayoutService interface {
	// GetRoot returns the root of the layout
	GetRoot() string

	// Match returns the repository reference corresponding to the given path
	// If the path does not match the layout, it returns the error `repository.ErrNotMatched`
	Match(path string) (*repository.Reference, error)

	// PathFor returns the path corresponding to the given reference
	PathFor(ref repository.Reference) string

	// CreateRepositoryFolder creates a new folder for the repository
	CreateRepositoryFolder(reference repository.Reference) (string, error)

	// DeleteRepository deletes the repository
	DeleteRepository(reference repository.Reference) error
}
