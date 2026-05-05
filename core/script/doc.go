// Package script provides script management for executable hooks.
//
// This package manages executable scripts that can be triggered by hooks for
// repository operations. Scripts store executable content with timestamps for
// creation and update tracking.
//
// # Main Interfaces
//
//   - ScriptService: Manages script storage and retrieval
//
// # Main Types
//
//   - Script: Represents script metadata with UUID, name, and timestamps
//   - Entry: Input structure for creating scripts with content
//
// # Architecture
//
// Scripts store executable content that can be run by hooks. The ScriptService
// implements store.Content for persistence and provides methods to open script
// content as io.ReadCloser. Scripts track creation and update times and are
// referenced by UUID from hooks.
package script
