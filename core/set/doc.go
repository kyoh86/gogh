// Package set provides a generic thread-safe collection for UUID-indexed items.
//
// This package implements a thread-safe set data structure that stores items
// indexed by UUID. It supports flexible ID lookup using full UUIDs or prefixes,
// with proper duplicate detection and error handling.
//
// # Main Types
//
//   - Set[T]: Generic collection for items with UUID() method
//
// # Main Features
//
//   - Thread-safe operations with mutex protection
//   - Full UUID and prefix-based lookups
//   - Duplicate detection with ErrDuplicated
//   - Multiple match detection with ErrMultipleFound
//   - Iteration support with iter.Seq
//
// # Architecture
//
// The Set is generic over any type implementing UUID() uuid.UUID. It maintains
// both a map for O(1) lookups and a slice for ordered iteration. All operations
// are thread-safe with proper locking.
package set
