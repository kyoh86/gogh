package main

import (
	"context"
	"os"

	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/kyoh86/gogh/v2/view/repotab"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
)

var reposFlags struct {
	query string
}

var reposCommand = &cobra.Command{
	Use:   "repos",
	Short: "List remote repositories",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) error {
		list, err := app.Servers().List()
		if err != nil {
			return err
		}
		var options []repotab.Option
		if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			options = append(options, repotab.Width(width))
		}
		if term.IsTerminal(int(os.Stdout.Fd())) {
			options = append(options, repotab.Styled())
		}
		format := repotab.NewPrinter(os.Stdout, options...)
		defer format.Close()
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		eg, ctx := errgroup.WithContext(ctx)
		for _, server := range list {
			server := server
			eg.Go(func() error {
				adaptor, err := github.NewAdaptor(ctx, server.Host(), server.Token())
				if err != nil {
					return err
				}
				remote := gogh.NewRemoteController(adaptor)
				sch, ech := remote.ListAsync(ctx, &gogh.RemoteListOption{
					Query: reposFlags.query,
				})
				for {
					select {
					case spec, more := <-sch:
						if !more {
							return nil
						}
						format.Print(spec)
					case err := <-ech:
						if err != nil {
							return err
						}
					case <-ctx.Done():
						if err := ctx.Err(); err != nil {
							return err
						}
					}
				}
			})
		}
		return eg.Wait()
	},
}

func init() {
	// TODO: filter flags
	//          --archived          Show only archived repositories
	//          --fork              Show only forks
	//      -l, --language string   Filter by primary coding language
	//      -L, --limit int         Maximum number of repositories to list (default 30)
	//          --no-archived       Omit archived repositories
	//          --private           Show only private repositories
	//          --public            Show only public repositories
	//          --source            Show only non-forks
	//
	// TODO: style flags
	//          --format            table,spec,URL,json
	//          --color             auto,never,always
	reposCommand.Flags().StringVarP(&reposFlags.query, "query", "", "", "Query for selecting projects")
	facadeCommand.AddCommand(reposCommand)
}
