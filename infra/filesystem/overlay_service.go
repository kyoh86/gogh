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

	"github.com/bmatcuk/doublestar/v4"
	corefs "github.com/kyoh86/gogh/v4/core/fs"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// OverlayService implements workspace.OverlayService using the filesystem
type OverlayService struct {
	fsys corefs.FS
}

// NewOverlayService creates a new OverlayService instance with the given filesystem
func NewOverlayService(fsys corefs.FS) (*OverlayService, error) {
	service := &OverlayService{
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
func EncodeFileName(entry workspace.OverlayEntry) string {
	patternEncoded := encoding.EncodeToString([]byte(entry.RepoPattern))
	t := "overlay"
	if entry.ForInit {
		t = "init"
	}
	return strings.Join([]string{patternEncoded, t, entry.RelativePath}, separator)
}

// DecodeFileName decodes an encoded filename back to repo-pattern and relative path
func DecodeFileName(encodedName string) (*workspace.OverlayEntry, error) {
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

	return &workspace.OverlayEntry{
		RepoPattern:  string(patternBytes),
		ForInit:      t,
		RelativePath: parts[2],
	}, nil
}

// FindOverlays implements workspace.OverlayService
func (s *OverlayService) FindOverlays(ctx context.Context, ref repository.Reference) iter.Seq2[*workspace.OverlayEntry, error] {
	return func(yield func(*workspace.OverlayEntry, error) bool) {
		// Convert repository reference to a string for pattern matching
		repoString := ref.String()
		entries, err := s.ListOverlays(ctx)
		if err != nil {
			yield(nil, fmt.Errorf("listing overlays: %w", err))
			return
		}
		for _, entry := range entries {
			// Check if this entry should be applied to this repository
			match, err := doublestar.Match(entry.RepoPattern, repoString)
			if err != nil {
				yield(nil, fmt.Errorf("matching repo-pattern '%s': %w", entry.RepoPattern, err))
				return
			}
			if !match {
				continue
			}
			e := entry
			if !yield(&e, nil) {
				return
			}
		}
	}
}

// OpenOverlay implements workspace.OverlayService
func (s *OverlayService) OpenOverlay(ctx context.Context, entry workspace.OverlayEntry) (io.ReadCloser, error) {
	contentPath := EncodeFileName(entry)

	file, err := s.fsys.Open(contentPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("overlay not found: repo-pattern=%s, path=%s", entry.RepoPattern, entry.RelativePath)
		}
		return nil, fmt.Errorf("opening overlay file: %w", err)
	}

	return file, nil
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

		entry, err := DecodeFileName(path)
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
	contentPath := EncodeFileName(entry)

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
	contentPath := EncodeFileName(entry)

	if err := s.fsys.Remove(contentPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("overlay not found: repo-pattern=%s, path=%s", entry.RepoPattern, entry.RelativePath)
		}
		return fmt.Errorf("removing overlay file: %w", err)
	}

	return nil
}

// Ensure OverlayService implements workspace.OverlayService
var _ workspace.OverlayService = (*OverlayService)(nil)
