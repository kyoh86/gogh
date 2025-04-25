package commands

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v3"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/spf13/cobra"
)

func NewBundleDumpCommand(conf *config.ConfigStore, defaults *config.FlagStore) *cobra.Command {
	var f config.BundleDumpFlags
	cmd := &cobra.Command{
		Use:     "dump",
		Aliases: []string{"export"},
		Short:   "Export current local projects",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := os.Stdout
			if f.File.Expand() != "" {
				f, err := os.OpenFile(
					f.File.Expand(),
					os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
					0644,
				)
				if err != nil {
					return err
				}
				defer f.Close()
				out = f
			}
			ctx := cmd.Context()
			list := conf.GetRoots()
			if len(list) == 0 {
				return nil
			}
			local := gogh.NewLocalController(list[0])
			if err := local.Walk(ctx, nil, func(project gogh.Project) error {
				utxt, err := gogh.GetDefaultRemoteURLFromLocalProject(ctx, project)
				if err != nil {
					if errors.Is(err, git.ErrRemoteNotFound) {
						return nil
					}
					return err
				}
				uobj, err := url.Parse(utxt)
				if err != nil {
					return err
				}
				remoteName := strings.Join([]string{uobj.Host, strings.TrimPrefix(strings.TrimSuffix(uobj.Path, ".git"), "/")}, "/")
				localName := project.RelPath()
				if remoteName == localName {
					fmt.Fprintln(out, localName)
					return nil
				}
				fmt.Fprintf(out, "%s=%s\n", remoteName, localName)
				return nil
			}); err != nil {
				return err
			}
			return nil
		},
	}

	f.File = defaults.BundleDump.File
	cmd.Flags().VarP(&f.File, "file", "", "A file to output; if not specified, output to stdout")
	return cmd
}
