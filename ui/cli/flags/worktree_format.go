package flags

import (
	"fmt"

	"github.com/kyoh86/gogh/v4/core/worktree"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type WorktreeFormat string

var _ pflag.Value = (*WorktreeFormat)(nil)

func (f WorktreeFormat) String() string {
	return string(f)
}

func (f *WorktreeFormat) Set(v string) error {
	_, err := worktree.ParseFormat(v)
	if err != nil {
		return fmt.Errorf("parse worktree format: %w", err)
	}
	*f = WorktreeFormat(v)
	return nil
}

func (f WorktreeFormat) Type() string {
	return "string"
}

const WorktreeFormatShortUsage = `Print worktree in a given format, where [format] can be one of "default", "full-path", "json", "fields" or "fields:[separator]".`

const WorktreeFormatLongUsage = `
Print worktree in a given format, where [format] can be one of "default",
"full-path", "json", "fields" and "fields:[separator]".

- default

	Repository name with branch name in parentheses for non-main worktrees.
	For example: "github.com/kyoh86/gogh (feature-branch)"

- full-path

	A full path of the worktree.  For example:
	"/root/Projects/github.com/kyoh86/gogh/.worktree/main"

- json

	JSON format with all worktree properties.

- fields

	Tab separated all formats and properties of the worktree.
	i.e. [full-path]\t[repo]\t[branch]\t[commit]

- fields:[separator]

	Like "fields" but with the explicit separator.
`

func WorktreeFormatFlag(cmd *cobra.Command, format *WorktreeFormat, defaultValue string) error {
	// UNDONE: opt ...Options Accepts NameOption, ShortUsageOption, ShorthandOption
	if defaultValue != "" {
		if err := format.Set(defaultValue); err != nil {
			return fmt.Errorf("setting default format: %w", err)
		}
	}
	cmd.Flags().VarP(format, "format", "f", WorktreeFormatShortUsage)
	if err := cmd.RegisterFlagCompletionFunc("format", CompleteWorktreeFormat); err != nil {
		return fmt.Errorf("registering completion function for format flag: %w", err)
	}
	return nil
}

func CompleteWorktreeFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"default", "full-path", "json", "fields", "fields:"}, cobra.ShellCompDirectiveDefault
}
