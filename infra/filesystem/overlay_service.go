package filesystem

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/wfs"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// OverlayService implements workspace.OverlayService using the filesystem
type OverlayService struct {
	fsys wfs.WFS
}

// NewOverlayService creates a new OverlayService instance with the given filesystem
func NewOverlayService(fsys wfs.WFS) (*OverlayService, error) {
	service := &OverlayService{
		fsys: fsys,
	}

	// Ensure the root directory exists
	if err := service.fsys.MkdirAll("", 0755); err != nil {
		return nil, fmt.Errorf("creating overlay directory: %w", err)
	}

	return service, nil
}

// separator is a string used to separate the pattern and relative path in the encoded filename
// Base64 encoded results will not contain this string
const separator = "/"

var encoding = base64.URLEncoding.WithPadding('.')

// encodeFileName safely encodes a pattern and path into a valid filename
func encodeFileName(entry workspace.OverlayEntry) string {
	patternEncoded := encoding.EncodeToString([]byte(entry.Pattern))
	t := "overlay"
	if entry.ForInit {
		t = "init"
	}
	return strings.Join([]string{patternEncoded, t, entry.RelativePath}, separator)
}

// decodeFileName decodes an encoded filename back to pattern and path
func decodeFileName(encodedName string) (*workspace.OverlayEntry, error) {
	parts := strings.SplitN(encodedName, separator, 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid encoded filename format")
	}

	patternBytes, err := encoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("decoding pattern: %w", err)
	}
	t := false
	switch parts[1] {
	case "overlay":
	case "init":
		t = true
	default:
		return nil, fmt.Errorf("unknown overlay type: %s", parts[1])
	}

	return &workspace.OverlayEntry{
		Pattern:      string(patternBytes),
		ForInit:      t,
		RelativePath: parts[2],
	}, nil
}

// FindOverlays implements workspace.OverlayService
func (s *OverlayService) FindOverlays(ctx context.Context, ref repository.Reference) iter.Seq2[*workspace.Overlay, error] {
	return func(yield func(*workspace.Overlay, error) bool) {
		// Convert repository reference to a string for pattern matching
		repoString := ref.String()
		entries, err := s.ListOverlays(ctx)
		if err != nil {
			yield(nil, fmt.Errorf("listing overlays: %w", err))
			return
		}
		for _, entry := range entries {
			// Check if this entry should be applied to this repository
			match, err := doublestar.Match(entry.Pattern, repoString)
			if err != nil {
				yield(nil, fmt.Errorf("matching pattern '%s': %w", entry.Pattern, err))
				return
			}
			if !match {
				continue
			}

			// Get the content and yield it
			func() {
				file, err := os.Open(encodeFileName(entry))
				if err != nil {
					yield(nil, fmt.Errorf("opening overlay file '%s': %w", entry.RelativePath, err))
					return
				}
				defer file.Close()
				if !yield(&workspace.Overlay{
					Content:      io.NopCloser(file),
					ForInit:      entry.ForInit,
					RelativePath: entry.RelativePath,
				}, nil) {
					return // Stop
				}
			}()
		}
	}
}

// ListOverlays implements workspace.OverlayService
func (s *OverlayService) ListOverlays(ctx context.Context) ([]workspace.OverlayEntry, error) {
	var result []workspace.OverlayEntry

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

		entry, err := decodeFileName(path)
		if err != nil {
			// Skip files that don't follow our encoding format
			return nil
		}

		result = append(result, *entry)

		return nil
	})

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return []workspace.OverlayEntry{}, nil
		}
		return nil, fmt.Errorf("walking overlay directory: %w", err)
	}

	return result, nil
}

// AddOverlay implements workspace.OverlayService
func (s *OverlayService) AddOverlay(ctx context.Context, entry workspace.OverlayEntry, content io.Reader) error {
	// Write the content to a file
	contentPath := encodeFileName(entry)

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

// RemoveOverlay implements workspace.OverlayService
func (s *OverlayService) RemoveOverlay(ctx context.Context, entry workspace.OverlayEntry) error {
	contentPath := encodeFileName(entry)

	if err := s.fsys.Remove(contentPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("overlay not found: pattern=%s, path=%s", entry.Pattern, entry.RelativePath)
		}
		return fmt.Errorf("removing overlay file: %w", err)
	}

	return nil
}

// Ensure OverlayService implements workspace.OverlayService
var _ workspace.OverlayService = (*OverlayService)(nil)
