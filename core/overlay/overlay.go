package overlay

import (
	"fmt"

	doublestar "github.com/bmatcuk/doublestar/v4"
	"github.com/kyoh86/gogh/v4/core/repository"
)

// Overlay represents the metadata for an overlay entry.
type Overlay struct {
	RepoPattern     string `json:"repoPattern"`     // Repository pattern (glob)
	ForInit         bool   `json:"forInit"`         // Whether the overlay is for initialization only
	RelativePath    string `json:"relativePath"`    // Relative path in the repository where the overlay file will be placed
	ContentLocation string `json:"contentLocation"` // Location of the content to be copied
}

func (ov Overlay) ID() string {
	return fmt.Sprintf("%q%v%q", ov.RepoPattern, ov.ForInit, ov.RelativePath)
}

func (ov Overlay) String() string {
	return ov.RelativePath + "@" + ov.RepoPattern
}

func (ov Overlay) Match(ref repository.Reference) (bool, error) {
	return doublestar.Match(ov.RepoPattern, ref.String())
}
