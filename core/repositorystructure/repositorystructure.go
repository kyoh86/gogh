package repositorystructure

import (
	"fmt"

	"github.com/spf13/pflag"
)

// RepositoryStructure represents the structure type for repository storage.
type RepositoryStructure string

const (
	// StructureWorktree uses bare repository + .worktree directories.
	StructureWorktree RepositoryStructure = "worktree"
	// StructureNormal uses traditional git repository structure.
	StructureNormal RepositoryStructure = "normal"
)

var _ pflag.Value = (*RepositoryStructure)(nil)

// ParseRepositoryStructure parses a string into RepositoryStructure.
func ParseRepositoryStructure(v string) (RepositoryStructure, error) {
	switch v {
	case string(StructureWorktree), "":
		return StructureWorktree, nil
	case string(StructureNormal):
		return StructureNormal, nil
	default:
		return "", fmt.Errorf("invalid structure: %q (must be %q or %q)", v, StructureWorktree, StructureNormal)
	}
}

// Set implements pflag.Value interface.
func (s *RepositoryStructure) Set(v string) error {
	parsed, err := ParseRepositoryStructure(v)
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}

// String implements pflag.Value interface.
func (s RepositoryStructure) String() string {
	if s == "" {
		return string(StructureWorktree)
	}
	return string(s)
}

// Type implements pflag.Value interface.
func (s RepositoryStructure) Type() string {
	return "string"
}

// IsWorktree returns true if the structure is worktree.
func (s RepositoryStructure) IsWorktree() bool {
	return s == StructureWorktree
}

// IsNormal returns true if the structure is normal.
func (s RepositoryStructure) IsNormal() bool {
	return s == StructureNormal || s == ""
}

// MarshalText implements encoding.TextMarshaler.
func (s RepositoryStructure) MarshalText() ([]byte, error) {
	if s == "" {
		return []byte(StructureWorktree), nil
	}
	return []byte(s), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *RepositoryStructure) UnmarshalText(text []byte) error {
	parsed, err := ParseRepositoryStructure(string(text))
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}
