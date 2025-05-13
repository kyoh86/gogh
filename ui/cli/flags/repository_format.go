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

func (f RemoteRepoFormat) Formatter(w io.Writer) (view.RepositoryPrinter, error) {
	return remoteRepoFormatter(string(f), w)
}

func remoteRepoFormatter(v string, w io.Writer) (view.RepositoryPrinter, error) {
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

const RemoteRepoFormatShortUsage = `
Print each repository in a given format, where [format] can be one of "table", "ref",
"url" or "json".
`

func CompleteRemoteRepoFormat(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"table", "ref", "url", "json"}, cobra.ShellCompDirectiveDefault
}
