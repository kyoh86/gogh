// Package git provides an infrastructure layer implementation for Git operations.
//
// This package implements the core Git interfaces using the go-git library
// and system git commands. It provides concrete implementations for:
//
//   - GitService: Core Git operations like clone, init, fetch, branch management
//
// # Requirements
//
// This package requires Git to be installed on the system and available in PATH.
// While most operations use go-git/v5, some functionality relies on system git commands:
//
//   - git worktree add: go-git does not support worktree operations
//   - git remote set-head: needed for setting remote HEAD references
//   - git symbolic-ref: needed for resolving remote HEAD
//
// These external commands are used as fallbacks where go-git functionality is incomplete.
package git
