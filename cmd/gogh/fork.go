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
		if _, err := local.Clone(ctx, spec, server, nil); err != nil {
			return err
		}
		return local.SetRemoteSpecs(ctx, spec, map[string][]gogh.Spec{
			git.DefaultRemoteName: {forked.Spec},
			"upstream":            {spec},
		})
	},
}

func init() {
	facadeCommand.AddCommand(forkCommand)
}
