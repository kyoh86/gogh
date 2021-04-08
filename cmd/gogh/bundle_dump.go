package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/spf13/cobra"
)

var bundleDumpFlags struct{}

var bundleDumpCommand = &cobra.Command{
	Use:     "dump",
	Aliases: []string{"export"},
	Short:   "Export current local projects",
	Args:    cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		roots := app.Roots()
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
			if project.URL() != urls[0] {
				uobj, err := url.Parse(urls[0])
				if err != nil {
					return err
				}
				remoteSpec := strings.TrimSuffix(uobj.Path, ".git")
				fmt.Printf("%s=%s\n", localSpec.String(), strings.TrimPrefix(remoteSpec, "/"))
				return nil
			}
			fmt.Println(localSpec.String())
			return nil
		}); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	bundleCommand.AddCommand(bundleDumpCommand)
}
