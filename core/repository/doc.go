// Package repository provides core types and services for repository identification and management.
//
// This package defines the fundamental types for identifying repositories (Reference),
// locating them in the filesystem (Location), and parsing repository specifications.
// It also provides validation for repository components and default name management.
//
// # Main Types
//
//   - Reference: Uniquely identifies a repository (host, owner, name)
//   - ReferenceWithAlias: Reference with optional local alias
//   - Location: Represents a repository's filesystem location
//   - ReferenceParser: Parses repository strings into References
//
// # Main Interfaces
//
//   - ReferenceParser: Interface for parsing repository specifications
//   - DefaultNameService: Manages default host and owner configuration
//
// # Architecture
//
// References use the host/owner/name format and support validation. The parser
// handles various input formats with intelligent defaulting. Locations provide
// both full filesystem paths and relative paths from workspace roots.
package repository
