// Package filesystem provides an infrastructure layer implementation for filesystem operations.
//
// This package implements the workspace interfaces using the actual filesystem.
// It provides concrete implementations for:
//
//   - WorkspaceService: Manages multiple workspace root directories
//   - FinderService: Searches and locates repositories across workspaces
//   - LayoutService: Handles repository path resolution and directory structure
//
// The implementation interacts directly with the filesystem to:
//   - Manage multiple root directories where repositories are stored
//   - Search for repositories by name or pattern
//   - Resolve repository paths according to the configured layout structure
//   - Handle repository folder creation and path normalization
package filesystem
