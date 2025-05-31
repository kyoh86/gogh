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
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// OverlayService implements workspace.OverlayService using the filesystem
type OverlayService struct {
	overlayDir string
	fsys       FSOperations
}

// FSOperations combines read and write operations for filesystem
type FSOperations interface {
	// fs.FS methods (read operations)
	fs.FS
	fs.ReadDirFS
	fs.ReadFileFS
	fs.StatFS

	// Write operations
	WriteFile(name string, data []byte, perm fs.FileMode) error
	MkdirAll(path string, perm fs.FileMode) error
	Remove(name string) error
	Create(name string) (io.WriteCloser, error)
}

// OSFileSystem implements FSOperations using the os package
type OSFileSystem struct{}

func (fs *OSFileSystem) Open(name string) (fs.File, error) {
	return os.Open(name)
}

func (fs *OSFileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (fs *OSFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (fs *OSFileSystem) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

func (fs *OSFileSystem) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (fs *OSFileSystem) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs *OSFileSystem) Remove(name string) error {
	return os.Remove(name)
}

func (fs *OSFileSystem) Create(name string) (io.WriteCloser, error) {
	return os.Create(name)
}

// NewOverlayService creates a new OverlayService instance
func NewOverlayService(overlayDir string) (*OverlayService, error) {
	service := &OverlayService{
		overlayDir: overlayDir,
		fsys:       &OSFileSystem{},
	}

	// Ensure the overlay directory exists
	if err := service.fsys.MkdirAll(overlayDir, 0755); err != nil {
		return nil, fmt.Errorf("creating overlay directory: %w", err)
	}

	return service, nil
}

// en: separator is a string used to separate the pattern and relative path
// Base64 encoded results do not contain this string
const separator = "--"

// encodeFileName safely encodes a pattern and path into a valid filename
func encodeFileName(pattern, relativePath string) string {
	patternEncoded := base64.URLEncoding.EncodeToString([]byte(pattern))
	pathEncoded := base64.URLEncoding.EncodeToString([]byte(relativePath))

	return patternEncoded + separator + pathEncoded
}

// decodeFileName decodes an encoded filename back to pattern and path
func decodeFileName(encodedName string) (pattern, relativePath string, err error) {
	parts := strings.Split(encodedName, separator)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid encoded filename format")
	}

	patternBytes, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", "", fmt.Errorf("decoding pattern: %w", err)
	}

	pathBytes, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", fmt.Errorf("decoding path: %w", err)
	}

	return string(patternBytes), string(pathBytes), nil
}

// getContentPath returns the path where the content for a given entry is stored
func (s *OverlayService) getContentPath(entry workspace.OverlayEntry) string {
	encodedName := encodeFileName(entry.Pattern, entry.RelativePath)
	return filepath.Join(s.overlayDir, encodedName)
}

// ApplyOverlays implements workspace.OverlayService
func (s *OverlayService) ApplyOverlays(ctx context.Context, ref repository.Reference, repoPath string) error {
	// Convert repository reference to a string format for pattern matching
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
		if err := s.fsys.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("creating directory '%s': %w", targetDir, err)
		}

		// Read content
		data, err := io.ReadAll(content)
		if err != nil {
			return fmt.Errorf("reading overlay content: %w", err)
		}

		// Write to target file
		if err := s.fsys.WriteFile(targetPath, data, 0644); err != nil {
			return fmt.Errorf("writing to file '%s': %w", targetPath, err)
		}
	}

	return nil
}

// ListOverlays implements workspace.OverlayService
func (s *OverlayService) ListOverlays(ctx context.Context) ([]workspace.OverlayEntry, error) {
	entries, err := s.fsys.ReadDir(s.overlayDir)
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

	// Create parent directories if needed
	dir := filepath.Dir(contentPath)
	if err := s.fsys.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

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
func (s *OverlayService) RemoveOverlay(ctx context.Context, pattern, relativePath string) error {
	entry := workspace.OverlayEntry{
		Pattern:      pattern,
		RelativePath: relativePath,
	}

	contentPath := s.getContentPath(entry)

	if err := s.fsys.Remove(contentPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("overlay not found: pattern=%s, path=%s", pattern, relativePath)
		}
		return fmt.Errorf("removing overlay file: %w", err)
	}

	return nil
}

// Ensure OverlayService implements workspace.OverlayService
var _ workspace.OverlayService = (*OverlayService)(nil)
