package filesystem

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
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

// encodeFileName safely encodes a pattern and path into a valid filename
func encodeFileName(pattern, relativePath string) string {
	patternEncoded := base64.URLEncoding.EncodeToString([]byte(pattern))

	return patternEncoded + separator + relativePath
}

// decodeFileName decodes an encoded filename back to pattern and path
func decodeFileName(encodedName string) (pattern, relativePath string, err error) {
	parts := strings.SplitN(encodedName, separator, 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid encoded filename format")
	}

	patternBytes, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", "", fmt.Errorf("decoding pattern: %w", err)
	}

	return string(patternBytes), parts[1], nil
}

// getContentPath returns the path where the content for a given entry is stored
func (s *OverlayService) getContentPath(entry workspace.OverlayEntry) string {
	return encodeFileName(entry.Pattern, entry.RelativePath)
}

// ApplyOverlays implements workspace.OverlayService
func (s *OverlayService) ApplyOverlays(ctx context.Context, ref repository.Reference, repoPath string) error {
	// Convert repository reference to a string for pattern matching
	repoString := ref.String()

	entries, err := s.ListOverlays(ctx)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// Check if this entry should be applied to this repository
		match, err := doublestar.Match(entry.Pattern, repoString)
		if err != nil {
			return fmt.Errorf("matching pattern '%s': %w", entry.Pattern, err)
		}

		if !match {
			continue
		}

		// Get the content
		content, err := s.GetOverlayContent(ctx, entry)
		if err != nil {
			return fmt.Errorf("getting content for '%s': %w", entry.RelativePath, err)
		}
		defer content.Close()

		// Create the target file
		targetPath := filepath.Join(repoPath, entry.RelativePath)

		// Ensure the directory exists
		targetDir := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("creating directory '%s': %w", targetDir, err)
		}

		// Read content
		data, err := io.ReadAll(content)
		if err != nil {
			return fmt.Errorf("reading overlay content: %w", err)
		}

		// Write to target file
		if err := os.WriteFile(targetPath, data, 0644); err != nil {
			return fmt.Errorf("writing to file '%s': %w", targetPath, err)
		}
	}

	return nil
}

// ListOverlays implements workspace.OverlayService
func (s *OverlayService) ListOverlays(ctx context.Context) ([]workspace.OverlayEntry, error) {
	entries, err := s.fsys.ReadDir("")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return []workspace.OverlayEntry{}, nil
		}
		return nil, fmt.Errorf("reading overlay directory: %w", err)
	}

	result := make([]workspace.OverlayEntry, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		pattern, relativePath, err := decodeFileName(entry.Name())
		if err != nil {
			// Skip files that don't follow our encoding format
			continue
		}

		result = append(result, workspace.OverlayEntry{
			Pattern:      pattern,
			RelativePath: relativePath,
		})
	}

	return result, nil
}

// GetOverlayContent implements workspace.OverlayService
func (s *OverlayService) GetOverlayContent(ctx context.Context, entry workspace.OverlayEntry) (io.ReadCloser, error) {
	contentPath := s.getContentPath(entry)
	file, err := s.fsys.Open(contentPath)
	if err != nil {
		return nil, fmt.Errorf("opening overlay file: %w", err)
	}

	return file, nil
}

// AddOverlay implements workspace.OverlayService
func (s *OverlayService) AddOverlay(ctx context.Context, entry workspace.OverlayEntry, content io.Reader) error {
	// Write the content to a file
	contentPath := s.getContentPath(entry)

	// Read content
	data, err := io.ReadAll(content)
	if err != nil {
		return fmt.Errorf("reading content: %w", err)
	}

	// Write to file
	if err := s.fsys.WriteFile(contentPath, data, 0644); err != nil {
		return fmt.Errorf("writing overlay content: %w", err)
	}

	return nil
}

// RemoveOverlay implements workspace.OverlayService
func (s *OverlayService) RemoveOverlay(ctx context.Context, entry workspace.OverlayEntry) error {
	contentPath := s.getContentPath(entry)

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
