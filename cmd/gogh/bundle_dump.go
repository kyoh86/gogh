package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v2"
	"github.com/spf13/cobra"
)

var bundleDumpFlags struct {
	file string
}

var bundleDumpCommand = &cobra.Command{
	Use:     "dump",
	Aliases: []string{"export"},
	Short:   "Export current local projects",
	Args:    cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) error {
		out := os.Stdout
		if bundleDumpFlags.file != "" {
			f, err := os.OpenFile(bundleDumpFlags.file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer f.Close()
			out = f
		}
		ctx := cmd.Context()
		roots := Roots()
		if len(roots) == 0 {
			return nil
		}
		local := gogh.NewLocalController(roots[0])
		if err := local.Walk(ctx, nil, func(project gogh.Project) error {
			localSpec := project.Spec()
			urls, err := local.GetRemoteURLs(ctx, localSpec, git.DefaultRemoteName)
			if err != nil {
				return err
			}
			uobj, err := url.Parse(urls[0])
			if err != nil {
				return err
			}
			remoteName := strings.Join([]string{uobj.Host, strings.TrimPrefix(strings.TrimSuffix(uobj.Path, ".git"), "/")}, "/")
			localName := localSpec.String()
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

func init() {
	bundleDumpCommand.Flags().StringVarP(&bundleDumpFlags.file, "file", "", "", "A file to output")
	bundleCommand.AddCommand(bundleDumpCommand)
}
