package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/spf13/cobra"
)

var forkCommand = &cobra.Command{
	Use:   "fork [flags] OWNER/NAME",
	Short: "Fork a repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, specs []string) error {
		ctx := cmd.Context()
		servers := app.Servers()
		parser := gogh.NewSpecParser(servers)
		spec, server, err := parser.Parse(specs[0])
		if err != nil {
			return err
		}
		adaptor, err := github.NewAdaptor(ctx, server.Host(), server.Token())
		if err != nil {
			return err
		}
		remote := gogh.NewRemoteController(adaptor)
		forked, err := remote.Fork(ctx, spec.Owner(), spec.Name(), nil)
		if err != nil {
			return err
		}

		root := app.DefaultRoot()
		local := gogh.NewLocalController(root)
		upstream, err := local.Clone(ctx, spec, server, nil)
		if err != nil {
			return err
		}
		origin := gogh.NewProject(root, forked)
		return local.SetRemoteURLs(ctx, spec, map[string][]string{
			git.DefaultRemoteName: {origin.URL()},
			"upstream":            {upstream.URL()},
		})
	},
}

func init() {
	facadeCommand.AddCommand(forkCommand)
}
