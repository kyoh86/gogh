package flags

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RepositoryFormatFlag adds a flag to the command for specifying the format of the repository output.
func RepositoryFormatFlag(cmd *cobra.Command, format *string, defaultValue string) error {
	// UNDONE: opt ...Options Accepts NameOption, ShortUsageOption, ShorthandOption
	cmd.Flags().StringVarP(format, "format", "f", defaultValue, RepositoryFormatShortUsage)
	if err := cmd.RegisterFlagCompletionFunc("format", CompleteRepositoryFormat); err != nil {
		return fmt.Errorf("failed to register completion function for format flag: %w", err)
	}
	return nil
}

// RepositoryFormatShortUsage is the short usage description for the repository format flag.
const RepositoryFormatShortUsage = `
Print each repository in a given format, where [format] can be one of "table", "ref",
"url" or "json".
`

// CompleteRepositoryFormat provides completion options for the repository format flag.
func CompleteRepositoryFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"table", "ref", "url", "json"}, cobra.ShellCompDirectiveDefault
}
