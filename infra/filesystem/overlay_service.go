package filesystem

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	"github.com/gobwas/glob"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// OverlayService implements workspace.OverlayService
type OverlayService struct {
	patterns []workspace.OverlayPattern
	changed  bool
}

// NewOverlayService creates a new DefaultOverlayService
func NewOverlayService() workspace.OverlayService {
	return &OverlayService{
		patterns: []workspace.OverlayPattern{},
		changed:  false,
	}
}

// AddPattern implements workspace.OverlayService
func (s *OverlayService) AddPattern(pattern string, files []workspace.OverlayFile) error {
	// Check if pattern is valid
	_, err := glob.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern %s: %w", pattern, err)
	}

	// Check if pattern already exists
	for i, p := range s.patterns {
		if p.Pattern == pattern {
			s.patterns[i].Files = files
			s.changed = true
			return nil
		}
	}

	s.patterns = append(s.patterns, workspace.OverlayPattern{
		Pattern: pattern,
		Files:   files,
	})
	s.changed = true
	return nil
}

// RemovePattern implements workspace.OverlayService
func (s *OverlayService) RemovePattern(pattern string) error {
	for i, p := range s.patterns {
		if p.Pattern == pattern {
			s.patterns = slices.Delete(s.patterns, i, i+1)
			s.changed = true
			return nil
		}
	}
	return fmt.Errorf("pattern %s not found", pattern)
}

// GetPatterns implements workspace.OverlayService
func (s *OverlayService) GetPatterns() []workspace.OverlayPattern {
	return s.patterns
}

// MatchRepository checks which overlay files match the given repository
func (s *OverlayService) MatchRepository(repo string) []workspace.OverlayFile {
	var result []workspace.OverlayFile
	for _, pattern := range s.patterns {
		g, err := glob.Compile(pattern.Pattern)
		if err != nil {
			// Skip invalid patterns
			continue
		}
		if g.Match(repo) {
			result = append(result, pattern.Files...)
		}
	}
	return result
}

// ApplyToRepository implements workspace.OverlayService
func (s *OverlayService) ApplyToRepository(ctx context.Context, repoPath string, repo string) error {
	files := s.MatchRepository(repo)
	for _, file := range files {
		targetPath := filepath.Join(repoPath, file.TargetPath)

		// Skip if target file already exists
		if _, err := os.Stat(targetPath); err == nil {
			continue
		}

		// Create directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("create directory for %s: %w", targetPath, err)
		}

		// Copy file
		source, err := os.Open(file.SourcePath)
		if err != nil {
			return fmt.Errorf("open source file %s: %w", file.SourcePath, err)
		}
		defer source.Close()

		target, err := os.Create(targetPath)
		if err != nil {
			return fmt.Errorf("create target file %s: %w", targetPath, err)
		}
		defer target.Close()

		if _, err := io.Copy(target, source); err != nil {
			return fmt.Errorf("copy file from %s to %s: %w", file.SourcePath, targetPath, err)
		}
	}
	return nil
}

// MarkSaved implements workspace.OverlayService
func (s *OverlayService) MarkSaved() {
	s.changed = false
}

// HasChanges implements workspace.OverlayService
func (s *OverlayService) HasChanges() bool {
	return s.changed
}

var _ workspace.OverlayService = (*OverlayService)(nil)
