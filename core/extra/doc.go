// Package extra provides management of repository extras (overlays and hooks).
//
// This package manages collections of overlay and hook pairs that can be automatically
// applied to repositories or used as named templates. Extras support two types:
//
//   - TypeAuto: Automatically applied to specific repositories
//   - TypeNamed: Manually applied templates that can be reused
//
// # Main Types
//
//   - Extra: Represents a collection of overlay and hook pairs
//   - Item: Pairs an overlay ID with a hook ID
//   - Type: Distinguishes between auto and named extras
//
// # Architecture
//
// The ExtraService implements store.Content for persistence and provides methods
// to add, remove, and retrieve extras. Auto extras are tied to specific repositories
// while named extras serve as reusable templates.
package extra
