package overlay

import (
	doublestar "github.com/bmatcuk/doublestar/v4"
	"github.com/kyoh86/gogh/v4/core/repository"
)

// Overlay represents the metadata for an overlay entry.
type Overlay struct {
	RepoPattern     string
	ForInit         bool
	RelativePath    string
	ContentLocation string
}

func (ov *Overlay) String() string {
	return ov.RelativePath + "@" + ov.RepoPattern
}

func (ov Overlay) Match(ref repository.Reference) (bool, error) {
	return doublestar.Match(ov.RepoPattern, ref.String())
}
