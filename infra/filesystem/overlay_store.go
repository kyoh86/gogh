package filesystem

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"iter"
	"path/filepath"
	"strings"

	corefs "github.com/kyoh86/gogh/v4/core/fs"
	"github.com/kyoh86/gogh/v4/core/overlay"
)

// OverlayStore implements overlay.OverlayStore using the filesystem
type OverlayStore struct {
	fsys corefs.FS
}

// NewOverlayStore creates a new OverlayStore instance with the given filesystem
func NewOverlayStore(fsys corefs.FS) (*OverlayStore, error) {
	service := &OverlayStore{
		fsys: fsys,
	}

	// Ensure the root directory exists
	if err := service.fsys.MkdirAll("", 0755); err != nil {
		return nil, fmt.Errorf("creating overlay directory: %w", err)
	}

	return service, nil
}

// separator is a string used to separate the repo-pattern and relative path in the encoded filename
// Base64 encoded results will not contain this string
const separator = "/"

var encoding = base64.URLEncoding.WithPadding('.')

// EncodeFileName safely encodes a repo-pattern and relative path into a valid filename
func EncodeFileName(ov overlay.Overlay) string {
	patternEncoded := encoding.EncodeToString([]byte(ov.RepoPattern))
	t := "overlay"
	if ov.ForInit {
		t = "init"
	}
	return strings.Join([]string{patternEncoded, t, ov.RelativePath}, separator)
}

// DecodeFileName decodes an encoded filename back to repo-pattern and relative path
func DecodeFileName(encodedName string) (*overlay.Overlay, error) {
	parts := strings.SplitN(encodedName, separator, 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid encoded filename format")
	}

	patternBytes, err := encoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("decoding repo-pattern: %w", err)
	}
	t := false
	switch parts[1] {
	case "overlay":
	case "init":
		t = true
	default:
		return nil, fmt.Errorf("unknown overlay type: %s", parts[1])
	}

	return &overlay.Overlay{
		RepoPattern:  string(patternBytes),
		ForInit:      t,
		RelativePath: parts[2],
	}, nil
}

// ListOverlays implements overlay.OverlayStore
func (s *OverlayStore) ListOverlays(ctx context.Context) iter.Seq2[*overlay.Overlay, error] {
	return func(yield func(*overlay.Overlay, error) bool) {
		err := fs.WalkDir(s.fsys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Skip directories
			if d.IsDir() {
				return nil
			}

			// Skip files with normalized path "."
			if path == "." {
				return nil
			}

			ov, err := DecodeFileName(path)
			if err != nil {
				// Skip files that don't follow our encoding format
				return nil
			}

			if !yield(ov, nil) {
				return fs.SkipAll
			}

			return nil
		})

		if err != nil && err != fs.SkipAll {
			yield(nil, fmt.Errorf("walking overlay directory: %w", err))
		}
	}
}

// OpenOverlay implements overlay.OverlayStore
func (s *OverlayStore) OpenOverlay(ctx context.Context, ov overlay.Overlay) (io.ReadCloser, error) {
	contentPath := EncodeFileName(ov)

	file, err := s.fsys.Open(contentPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("overlay not found: repo-pattern=%s, path=%s", ov.RepoPattern, ov.RelativePath)
		}
		return nil, fmt.Errorf("opening overlay file: %w", err)
	}

	return file, nil
}

// AddOverlay implements overlay.OverlayStore
func (s *OverlayStore) AddOverlay(ctx context.Context, ov overlay.Overlay, content io.Reader) error {
	// Write the content to a file
	contentPath := EncodeFileName(ov)

	// Read content
	data, err := io.ReadAll(content)
	if err != nil {
		return fmt.Errorf("reading content: %w", err)
	}

	// Ensure the directory exists
	targetDir := filepath.Dir(contentPath)
	if err := s.fsys.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("creating directory '%s': %w", targetDir, err)
	}

	// Write to file
	if err := s.fsys.WriteFile(contentPath, data, 0644); err != nil {
		return fmt.Errorf("writing overlay content: %w", err)
	}

	return nil
}

// RemoveOverlay implements overlay.OverlayStore
func (s *OverlayStore) RemoveOverlay(ctx context.Context, ov overlay.Overlay) error {
	contentPath := EncodeFileName(ov)

	if err := s.fsys.Remove(contentPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("overlay not found: repo-pattern=%s, path=%s", ov.RepoPattern, ov.RelativePath)
		}
		return fmt.Errorf("removing overlay file: %w", err)
	}

	return nil
}

// Ensure OverlayStore implements overlay.OverlayStore
var _ overlay.OverlayStore = (*OverlayStore)(nil)
