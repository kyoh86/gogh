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

type RemoteRepoFormat string

var _ pflag.Value = (*RemoteRepoFormat)(nil)

func (f RemoteRepoFormat) String() string {
	return string(f)
}

func (f *RemoteRepoFormat) Set(v string) error {
	_, err := remoteRepoFormatter(v, os.Stdout)
	if err != nil {
		return fmt.Errorf("parse remote repo format: %w", err)
	}
	*f = RemoteRepoFormat(v)
	return nil
}

func (f RemoteRepoFormat) Type() string {
	return "string"
}

func (f RemoteRepoFormat) Formatter(w io.Writer) (view.RemoteRepoPrinter, error) {
	return remoteRepoFormatter(string(f), w)
}

func remoteRepoFormatter(v string, w io.Writer) (view.RemoteRepoPrinter, error) {
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
	case "ref":
		return view.NewRemoteRepoRefPrinter(w), nil
	case "url":
		return view.NewRemoteRepoURLPrinter(w), nil
	case "json":
		return view.NewRemoteRepoJSONPrinter(w), nil
	}
	return nil, fmt.Errorf("invalid format: %q", v)
}

const RemoteRepoFormatShortUsage = `
Print each repository in a given format, where [format] can be one of "table", "ref",
"url" or "json".
`

func CompleteRemoteRepoFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"table", "ref", "url", "json"}, cobra.ShellCompDirectiveDefault
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
