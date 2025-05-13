package flags

import (
	"fmt"
	"io"
	"os"

	"github.com/kyoh86/gogh/v3/ui/cli/view"
	"github.com/kyoh86/gogh/v3/ui/cli/view/repotab"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func RepositoryFormatFlag(cmd *cobra.Command, format *RepositoryFormat, defaultValue string) error {
	// UNDONE: opt ...Options Accepts NameOption, ShortUsageOption, ShorthandOption
	if defaultValue != "" {
		if err := format.Set(defaultValue); err != nil {
			return fmt.Errorf("failed to set default format: %w", err)
		}
	}
	cmd.Flags().VarP(format, "format", "f", RepositoryFormatShortUsage)
	if err := cmd.RegisterFlagCompletionFunc("format", CompleteRepositoryFormat); err != nil {
		return fmt.Errorf("failed to register completion function for format flag: %w", err)
	}
	return nil
}

type RepositoryFormat string

var _ pflag.Value = (*RepositoryFormat)(nil)

func (f RepositoryFormat) String() string {
	return string(f)
}

func (f *RepositoryFormat) Set(v string) error {
	_, err := repositoryFormatter(v, os.Stdout)
	if err != nil {
		return fmt.Errorf("parse remote repo format: %w", err)
	}
	*f = RepositoryFormat(v)
	return nil
}

func (f RepositoryFormat) Type() string {
	return "string"
}

func (f RepositoryFormat) Formatter(w io.Writer) (view.RepositoryPrinter, error) {
	return repositoryFormatter(string(f), w)
}

func repositoryFormatter(v string, w io.Writer) (view.RepositoryPrinter, error) {
	switch v {
	case "", "table":
		return repotab.NewPrinter(w, repotab.TermWidth(), repotab.Styled(false)), nil
	case "ref":
		return view.NewRepositoryPrinterRef(w), nil
	case "url":
		return view.NewRepositoryPrinterURL(w), nil
	case "json":
		return view.NewRepositoryPrinterJSON(w), nil
	}
	return nil, fmt.Errorf("invalid format: %q", v)
}

const RepositoryFormatShortUsage = `
Print each repository in a given format, where [format] can be one of "table", "ref",
"url" or "json".
`

func CompleteRepositoryFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"table", "ref", "url", "json"}, cobra.ShellCompDirectiveDefault
}
