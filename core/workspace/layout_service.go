package workspace

import (
	"errors"

	"github.com/kyoh86/gogh/v3/core/repository"
)

var ErrNotMatched = errors.New("repository not matched for a layout")

// LayoutService defines the layout structure of a repository under a root
type LayoutService interface {
	// GetRoot returns the root of the layout
	GetRoot() string

	// Match returns the repository reference corresponding the given path
	// If the path does not match the layout, it returns the error `repository.ErrNotMatched`
	// Example:
	// Match("github.com/owner/repo") returns "github.com/owner/repo"
	// Match("github.com/owner/repo/foo") returns "github.com/owner/repo"
	Match(path string) (*repository.Reference, error)

	// ExactMatch returns the repository reference corresponding exactly to the given path
	// If the path does not match the layout, it returns the error `repository.ErrNotMatched`
	// Example:
	// ExactMatch("github.com/owner/repo") returns "github.com/owner/repo"
	// ExactMatch("github.com/owner/repo/foo") returns `repository.ErrNotMatched`
	ExactMatch(path string) (*repository.Reference, error)

	// PathFor returns the path corresponding to the given reference
	PathFor(ref repository.Reference) string

	// CreateRepositoryFolder creates a new folder for the repository
	CreateRepositoryFolder(reference repository.Reference) (string, error)

	// DeleteRepository deletes the repository
	DeleteRepository(reference repository.Reference) error
}
