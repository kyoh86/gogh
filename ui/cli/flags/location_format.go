package flags

import (
	"fmt"
	"strings"

	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func LocationFormatFlag(cmd *cobra.Command, format *LocationFormat, defaultValue string) error {
	// UNDONE: opt ...Options Accepts NameOption, ShortUsageOption, ShorthandOption
	if defaultValue != "" {
		if err := format.Set(defaultValue); err != nil {
			return fmt.Errorf("failed to set default format: %w", err)
		}
	}
	cmd.Flags().VarP(format, "format", "f", LocationFormatShortUsage)
	if err := cmd.RegisterFlagCompletionFunc("format", CompleteLocationFormat); err != nil {
		return fmt.Errorf("failed to register completion function for format flag: %w", err)
	}
	return nil
}

type LocationFormat string

var _ pflag.Value = (*LocationFormat)(nil)

func (f LocationFormat) String() string {
	return string(f)
}

func (f *LocationFormat) Set(v string) error {
	_, err := formatter(v)
	if err != nil {
		return fmt.Errorf("parse local repo format: %w", err)
	}
	*f = LocationFormat(v)
	return nil
}

func (f LocationFormat) Type() string {
	return "string"
}

func (f LocationFormat) Formatter() (repository.LocationFormat, error) {
	return formatter(string(f))
}

func formatter(v string) (repository.LocationFormat, error) {
	switch v {
	case "", "rel-path", "rel", "path", "rel-file-path":
		return repository.LocationFormatPath, nil
	case "full-file-path", "full":
		return repository.LocationFormatFullPath, nil
	case "json":
		return repository.LocationFormatJSON, nil
	case "fields":
		return repository.LocationFormatFields("\t"), nil
	}
	if strings.HasPrefix(v, "fields:") {
		return repository.LocationFormatFields(v[len("fields:"):]), nil
	}
	return nil, fmt.Errorf("invalid format: %q", v)
}

const LocationFormatShortUsage = `Print local repository in a given format, where [format] can be one of "path", "full-path", "json", "fields" or "fields:[separator]".`

const LocationFormatLongUsage = `
Print local repository in a given format, where [format] can be one of "path",
"full-path", "fields" and "fields:[separator]".

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
