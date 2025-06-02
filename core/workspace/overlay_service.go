package workspace

import (
	"context"
	"fmt"
	"io"
	"iter"

	"github.com/kyoh86/gogh/v4/core/repository"
)

// OverlayEntry represents a file to be overlaid onto repositories
type OverlayEntry struct {
	// Pattern is a glob pattern that matches repository references
	// Examples:
	// - "*" matches all repositories
	// - "github.com/*" matches all GitHub repositories
	// - "github.com/kyoh86/*" matches all kyoh86's repositories on GitHub
	// - "github.com/kyoh86/gogh" matches only the gogh repository
	Pattern string
	// ForInit indicates whether the overlay should be applied only during repository initialization
	ForInit bool
	// RelativePath is the path relative to the repository root where the file should be placed
	// Example: ".envrc", "scripts/setup.sh"
	RelativePath string
}

func (e OverlayEntry) String() string {
	if e.ForInit {
		return fmt.Sprintf("Init(%s): %s", e.Pattern, e.RelativePath)
	} else {
		return fmt.Sprintf("Overlay(%s): %s", e.Pattern, e.RelativePath)
	}
}

// OverlayService provides functionality to add files to repositories after they are cloned
type OverlayService interface {
	// FindOverlays finds all overlay entries that match the given repository reference
	FindOverlays(ctx context.Context, ref repository.Reference) iter.Seq2[*OverlayEntry, error]

	// ListOverlays returns all registered overlay entries
	ListOverlays(ctx context.Context) ([]OverlayEntry, error)

	// AddOverlay adds a new overlay file
	AddOverlay(ctx context.Context, entry OverlayEntry, content io.Reader) error

	// RemoveOverlay removes an overlay file
	RemoveOverlay(ctx context.Context, entry OverlayEntry) error

	// OpenOverlay opens an overlay file for reading
	OpenOverlay(ctx context.Context, entry OverlayEntry) (io.ReadCloser, error)
}
