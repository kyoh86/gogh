package flags

import (
	"fmt"

	"github.com/spf13/cobra"
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

func (s RepositoryStructure) String() string {
	return string(s)
}

func (s *RepositoryStructure) Set(v string) error {
	_, err := ParseRepositoryStructure(v)
	if err != nil {
		return fmt.Errorf("parse repository structure: %w", err)
	}
	*s = RepositoryStructure(v)
	return nil
}

func (s RepositoryStructure) Type() string {
	return "string"
}

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

// IsWorktree returns true if the structure is worktree.
func (s RepositoryStructure) IsWorktree() bool {
	return s == StructureWorktree
}

// IsNormal returns true if the structure is normal.
func (s RepositoryStructure) IsNormal() bool {
	return s == StructureNormal || s == ""
}

const StructureShortUsage = `Repository structure to use (default: "worktree", one of "worktree" or "normal").`

const StructureLongUsage = `
Repository structure to use, where [structure] can be one of "worktree" or "normal".

- worktree

	Use bare repository + .worktree directories structure.
	The git repository is stored as a bare repository, and working trees
	are created in .worktree/<branch> directories.

- normal

	Use traditional git repository structure (v4 compatible).
	Working tree is directly in the repository directory.
`

// StructureFlag registers the structure flag with a command.
func StructureFlag(cmd *cobra.Command, structure *RepositoryStructure, defaultValue string) error {
	if defaultValue != "" {
		if err := structure.Set(defaultValue); err != nil {
			return fmt.Errorf("setting default structure: %w", err)
		}
	}
	cmd.Flags().VarP(structure, "structure", "s", StructureShortUsage)
	if err := cmd.RegisterFlagCompletionFunc("structure", CompleteStructure); err != nil {
		return fmt.Errorf("registering completion function for structure flag: %w", err)
	}
	return nil
}

// CompleteStructure provides completion for structure flag.
func CompleteStructure(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"worktree", "normal"}, cobra.ShellCompDirectiveDefault
}
