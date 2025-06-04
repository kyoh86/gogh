package workspace

import (
	"context"
	"fmt"
	"io"
	"iter"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/typ"
)

// Overlay represents a file to be overlaid onto repositories
type Overlay struct {
	// RepoPattern is a glob pattern that matches repository references
	// Examples:
	// - "*" matches all repositories
	// - "github.com/*" matches all GitHub repositories
	// - "github.com/kyoh86/*" matches all kyoh86's repositories on GitHub
	// - "github.com/kyoh86/gogh" matches only the gogh repository
	RepoPattern string
	// ForInit indicates whether the overlay should be applied only during repository initialization
	ForInit bool
	// RelativePath is the path relative to the repository root where the file should be placed
	// Example: ".envrc", "scripts/setup.sh"
	RelativePath string
}

func (e Overlay) String() string {
	if e.ForInit {
		return fmt.Sprintf("Init(%s): %s", e.RepoPattern, e.RelativePath)
	} else {
		return fmt.Sprintf("Overlay(%s): %s", e.RepoPattern, e.RelativePath)
	}
}

func (e Overlay) Match(ref repository.Reference) (bool, error) {
	return doublestar.Match(e.RepoPattern, ref.String())
}

func FilterOverlayForReference(entries iter.Seq2[*Overlay, error], ref repository.Reference) iter.Seq2[*Overlay, error] {
	return typ.FilterE(entries, func(entry *Overlay) (bool, error) {
		return entry.Match(ref)
	})
}

func FilterOverlayForPattern(entries iter.Seq2[*Overlay, error], repoPattern string) iter.Seq2[*Overlay, error] {
	return typ.FilterE(entries, func(entry *Overlay) (bool, error) {
		return entry.RepoPattern == repoPattern, nil
	})
}

// OverlayStore provides functionality to add files to repositories after they are cloned
type OverlayStore interface {
	// ListOverlays returns all registered overlay entries
	ListOverlays(ctx context.Context) iter.Seq2[*Overlay, error]

	// AddOverlay adds a new overlay file
	AddOverlay(ctx context.Context, entry Overlay, content io.Reader) error

	// RemoveOverlay removes an overlay file
	RemoveOverlay(ctx context.Context, entry Overlay) error

	// OpenOverlay opens an overlay file for reading
	OpenOverlay(ctx context.Context, entry Overlay) (io.ReadCloser, error)
}
