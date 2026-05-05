// Package store provides interfaces for loading and saving application content.
//
// This package defines foundational interfaces for content persistence, supporting
// both loading from sources and saving to destinations with change tracking.
//
// # Main Interfaces
//
//   - Content: Interface for types with change tracking
//   - Loader[T]: Interface for loading content from sources
//   - Saver[T]: Interface for saving content to destinations
//   - Store[T]: Combined loader and saver interface
//
// # Architecture
//
// The Content interface enables change tracking with HasChanges() and MarkSaved().
// Generic Loader and Saver interfaces support various persistence backends.
// The Store interface combines both for complete persistence operations.
package store
