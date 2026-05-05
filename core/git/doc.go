// Package git provides core interfaces for Git repository operations.
//
// This package defines the GitService interface for performing actual Git operations
// including cloning, initializing, fetching, branch management, and worktree operations.
// It serves as an abstraction layer over Git implementations.
//
// # Main Interfaces
//
//   - GitService: Core interface for Git operations (clone, init, fetch, worktrees)
//
// # Main Types
//
//   - CloneOptions: Options for clone operations
//   - InitOptions: Options for init operations
//
// # Architecture
//
// This package defines only interfaces - actual implementations are provided by
// infrastructure layers (e.g., infra/git). The interface supports both regular
// and bare repositories, with specific operations for worktree management.
package git
