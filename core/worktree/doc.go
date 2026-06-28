// Package worktree provides Git worktree management operations.
//
// This package manages Git worktrees, allowing multiple working directories for
// a single repository. It supports listing, adding, and removing worktrees with
// flexible path formatting.
//
// # Main Interfaces
//
//   - Service: Manages worktree operations (list, add, remove, get)
//   - PathBuilder: Generates worktree paths
//   - Format: Formats worktree information for display
//
// # Main Types
//
//   - Worktree: Represents a Git worktree with repository, branch, path, and commit
//
// # Architecture
//
// Worktrees use a .wt/ subdirectory structure. The Service interface provides
// CRUD operations for worktrees. PathBuilder generates paths preserving branch
// hierarchy. Format interface supports multiple output formats (default, full-path,
// json, fields) for displaying worktree information.
package worktree
