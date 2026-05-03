// Package github provides an infrastructure layer implementation for GitHub operations.
//
// This package implements the hosting and authentication service interfaces using
// the GitHub REST API. It provides concrete implementations for:
//
//   - HostingService: Repository operations (clone, create, fork, delete, etc.)
//   - AuthenticateService: GitHub authentication and token management
//
// The implementation uses the GitHub REST API v3 to interact with GitHub repositories,
// including operations like:
//   - Fetching repository information
//   - Creating and deleting repositories
//   - Forking repositories
//   - Creating repositories from templates
//   - Managing authentication tokens
package github
