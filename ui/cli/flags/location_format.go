package flags

import (
	"fmt"

	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func LocationFormatFlag(cmd *cobra.Command, format *LocationFormat, defaultValue string) error {
	// UNDONE: opt ...Options Accepts NameOption, ShortUsageOption, ShorthandOption
	if defaultValue != "" {
		if err := format.Set(defaultValue); err != nil {
			return fmt.Errorf("setting default format: %w", err)
		}
	}
	cmd.Flags().VarP(format, "format", "f", LocationFormatShortUsage)
	if err := cmd.RegisterFlagCompletionFunc("format", CompleteLocationFormat); err != nil {
		return fmt.Errorf("registering completion function for format flag: %w", err)
	}
	return nil
}

type LocationFormat string

var _ pflag.Value = (*LocationFormat)(nil)

func (f LocationFormat) String() string {
	return string(f)
}

func (f *LocationFormat) Set(v string) error {
	_, err := config.LocationFormatter(v)
	if err != nil {
		return fmt.Errorf("parse local repo format: %w", err)
	}
	*f = LocationFormat(v)
	return nil
}

func (f LocationFormat) Type() string {
	return "string"
}

const LocationFormatShortUsage = `Print local repository in a given format, where [format] can be one of "path", "full-path", "json", "fields" or "fields:[separator]".`

const LocationFormatLongUsage = `
Print local repository in a given format, where [format] can be one of "path",
"full-path", "json", "fields" and "fields:[separator]".

- path:

	A part of the URL to specify a repository.  For example: "github.com/kyoh86/gogh"

- full-path

	A full path of the local repository.  For example:
	"/root/Projects/github.com/kyoh86/gogh".

- fields

	Tab separated all formats and properties of the local repository.
	i.e. [full-path]\t[path]\t[host]\t[owner]\t[name]

- fields:[separator]

	Like "fields" but with the explicit separator.
`

func CompleteLocationFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"path", "full-path", "json", "fields", "fields:"}, cobra.ShellCompDirectiveDefault
}
