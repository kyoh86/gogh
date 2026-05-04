package flags

import (
	"fmt"

	"github.com/kyoh86/gogh/v4/core/repositorystructure"
	"github.com/spf13/cobra"
)

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
func StructureFlag(cmd *cobra.Command, structure *repositorystructure.RepositoryStructure, defaultValue string) error {
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
