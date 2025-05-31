package workspace

import (
	"context"
)

// OverlayFile represents a file to be placed in a repository
type OverlayFile struct {
	// SourcePath is the path to the source file (absolute path)
	SourcePath string
	// TargetPath is the path where the file should be placed (relative to repository root)
	TargetPath string
}

// OverlayPattern defines a pattern for repositories to receive overlay files
type OverlayPattern struct {
	// Pattern is a glob pattern for repository names (e.g. "kyoh86/*", "*/*", "github.com/kyoh86/*")
	Pattern string
	// Files is a list of files to be placed in repositories matching the pattern
	Files []OverlayFile
}

// OverlayService defines the interface for repository overlay files management
type OverlayService interface {
	// AddPattern adds a pattern for repositories to receive overlay files
	AddPattern(pattern string, files []OverlayFile) error
	// RemovePattern removes a pattern
	RemovePattern(pattern string) error
	// GetPatterns returns all patterns
	GetPatterns() []OverlayPattern
	// ApplyToRepository applies overlay files to the given repository
	ApplyToRepository(ctx context.Context, repoPath string, repo string) error
	// MarkSaved marks the service as saved
	MarkSaved()
	// HasChanges returns true if the service has unsaved changes
	HasChanges() bool
}
