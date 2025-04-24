package flags

import (
	"fmt"
	"io"
	"os"

	"github.com/kyoh86/gogh/v3/view"
	"github.com/kyoh86/gogh/v3/view/repotab"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

type RepoFormat string

var _ pflag.Value = (*RepoFormat)(nil)

func (f RepoFormat) String() string {
	return string(f)
}

func (f *RepoFormat) Set(v string) error {
	_, err := repoFormatter(v, os.Stdout)
	if err != nil {
		return fmt.Errorf("parse repo format: %w", err)
	}
	*f = RepoFormat(v)
	return nil
}

func (f RepoFormat) Type() string {
	return "string"
}

func (f RepoFormat) Formatter(w io.Writer) (view.RepositoryPrinter, error) {
	return repoFormatter(string(f), w)
}

func repoFormatter(v string, w io.Writer) (view.RepositoryPrinter, error) {
	switch v {
	case "", "table":
		var options []repotab.Option
		if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			options = append(options, repotab.Width(width))
		}
		if term.IsTerminal(int(os.Stdout.Fd())) {
			options = append(options, repotab.Styled())
		}
		return repotab.NewPrinter(w, options...), nil
	case "spec":
		return view.NewRepositorySpecPrinter(w), nil
	case "url":
		return view.NewRepositoryURLPrinter(w), nil
	case "json":
		return view.NewRepositoryJSONPrinter(w), nil
	}
	return nil, fmt.Errorf("invalid format: %q", v)
}

const RepoFormatShortUsage = `
Print each repository in a given format, where [format] can be one of "table", "spec",
"url" or "json".
`

func CompleteRepoFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"table", "spec", "url", "json"}, cobra.ShellCompDirectiveDefault
}

func GetColorOption(colorOpt string) string {
	if colorOpt != "" {
		return colorOpt
	}
	if term.IsTerminal(int(os.Stdout.Fd())) {
		return "auto"
	}
	return "never"
}
