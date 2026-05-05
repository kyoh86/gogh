// Package hosting provides interfaces for interacting with remote repository hosting services.
//
// This package defines the HostingService interface for managing remote repositories,
// including CRUD operations, URL parsing, and authentication token management. It
// provides a unified abstraction over different hosting service implementations.
//
// # Main Interfaces
//
//   - HostingService: Interface for remote repository operations
//
// # Main Types
//
//   - Repository: Represents a remote repository with metadata
//   - ParentRepository: Represents the parent of a forked repository
//   - ListRepositoryOptions: Options for listing repositories
//   - CreateRepositoryOptions: Options for creating repositories
//   - CreateRepositoryFromTemplateOptions: Options for creating from templates
//   - ForkRepositoryOptions: Options for forking repositories
//
// # Architecture
//
// This package defines only interfaces - actual implementations are provided by
// infrastructure layers (e.g., infra/github). The interface supports operations
// like listing, creating, forking, and deleting repositories, as well as
// converting between repository references and URLs.
package hosting
