// Package overlay provides overlay management for repository file customization.
//
// This package manages overlays that allow customizing repository files by
// overlaying additional content. Overlays are stored with their content and
// can be referenced by hooks or extras.
//
// # Main Interfaces
//
//   - OverlayService: Manages overlay storage and retrieval
//
// # Main Types
//
//   - Overlay: Represents overlay metadata with UUID, name, and relative path
//   - Entry: Input structure for creating overlays with content
//
// # Architecture
//
// Overlays store file content that can be applied to repositories. The OverlayService
// implements store.Content for persistence and provides methods to open overlay
// content as io.ReadCloser. Overlays are referenced by UUID from hooks and extras.
package overlay
