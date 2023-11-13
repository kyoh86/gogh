package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v2"
	"github.com/spf13/cobra"
)

type bundleDumpFlagsStruct struct {
	File expandedPath `yaml:"file,omitempty"`
}

var (
	bundleDumpFlags   bundleDumpFlagsStruct
	bundleDumpCommand = &cobra.Command{
		Use:     "dump",
		Aliases: []string{"export"},
		Short:   "Export current local projects",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			out := os.Stdout
			if bundleDumpFlags.File.expanded != "" {
				f, err := os.OpenFile(
					bundleDumpFlags.File.expanded,
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
			list := roots()
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
)

func init() {
	bundleDumpFlags.File = defaultFlag.BundleDump.File
	bundleDumpCommand.Flags().VarP(&bundleDumpFlags.File, "file", "", "A file to output")
	bundleCommand.AddCommand(bundleDumpCommand)
}
