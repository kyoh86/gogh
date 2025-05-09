package flags

import (
	"fmt"
	"strings"

	"github.com/kyoh86/gogh/v3/ui/cli/view"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type LocalRepoFormat string

var _ pflag.Value = (*LocalRepoFormat)(nil)

func (f LocalRepoFormat) String() string {
	return string(f)
}

func (f *LocalRepoFormat) Set(v string) error {
	_, err := formatter(v)
	if err != nil {
		return fmt.Errorf("parse local repo format: %w", err)
	}
	*f = LocalRepoFormat(v)
	return nil
}

func (f LocalRepoFormat) Type() string {
	return "string"
}

func (f LocalRepoFormat) Formatter() (view.LocalRepoFormat, error) {
	return formatter(string(f))
}

func formatter(v string) (view.LocalRepoFormat, error) {
	switch v {
	case "", "rel-path", "rel", "path", "rel-file-path":
		return view.LocalRepoFormatPath, nil
	case "full-file-path", "full":
		return view.LocalRepoFormatFullPath, nil
	case "json":
		return view.LocalRepoFormatJSON, nil
	case "fields":
		return view.LocalRepoFormatFields("\t"), nil
	}
	if strings.HasPrefix(v, "fields:") {
		return view.LocalRepoFormatFields(v[len("fields:"):]), nil
	}
	return nil, fmt.Errorf("invalid format: %q", v)
}

const LocalRepoFormatShortUsage = `Print each local repository in a given format, where [format] can be one of "path", "full-path", "json", "fields" or "fields:[separator]".`

func CompleteLocalRepoFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"path", "full-path", "json", "fields", "fields:"}, cobra.ShellCompDirectiveDefault
}
