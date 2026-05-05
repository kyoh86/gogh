// Package hook provides hook management for repository operations.
//
// This package manages hooks that can be triggered by repository events such as
// post-clone, post-fork, or post-create. Hooks can execute overlays or scripts
// based on pattern matching against repository references.
//
// # Main Interfaces
//
//   - HookService: Manages hook storage and retrieval
//
// # Main Types
//
//   - Hook: Represents a hook with trigger conditions and operations
//   - Entry: Input structure for creating hooks
//   - Event: Trigger events (post-clone, post-fork, post-create, any)
//   - OperationType: Type of operation (overlay or script)
//
// # Architecture
//
// Hooks use glob patterns to match repository references and can be triggered
// by specific events. Each hook references an overlay or script by UUID.
// The HookService implements store.Content for persistence.
package hook
