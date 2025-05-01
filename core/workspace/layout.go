package workspace

import (
	"errors"

	"github.com/kyoh86/gogh/v3/core/repository"
)

var ErrNotMatched = errors.New("repository layout not matched")

// Layout defines the layout structure of a repository
type Layout interface {
	// Match returns the repository reference corresponding to the given path
	// If the path does not match the layout, it returns the error `repository.ErrNotMatched`
	Match(root Root, path string) (*repository.Reference, error)

	// PathFor returns the path corresponding to the given reference
	PathFor(root Root, ref *repository.Reference) string
}
