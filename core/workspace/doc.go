// Package workspace provides workspace and repository layout management.
//
// This package manages the filesystem layout of repositories across multiple
// root directories. It provides services for workspace configuration, repository
// discovery, and path resolution.
//
// # Main Interfaces
//
//   - WorkspaceService: Manages multiple repository roots
//   - LayoutService: Defines repository layout structure under a root
//   - FinderService: Searches for repositories in workspaces
//
// # Main Types
//
//   - Root: Represents a repository root directory (alias for string)
//   - ListOptions: Options for repository search operations
//
// # Architecture
//
// Workspaces support multiple roots with a designated primary root. The LayoutService
// defines how repositories are organized (host/owner/name structure). The FinderService
// provides search and listing capabilities across all roots with pattern matching.
package workspace
