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
		var options []repotab.TableOption
		if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			options = append(options, repotab.OptionWidth(width))
		}
		if term.IsTerminal(int(os.Stdout.Fd())) {
			options = append(options, repotab.OptionStyled())
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
	// TODO:          --archived          Show only archived repositories
	// TODO:          --fork              Show only forks
	// TODO:      -l, --language string   Filter by primary coding language
	// TODO:      -L, --limit int         Maximum number of repositories to list (default 30)
	// TODO:          --no-archived       Omit archived repositories
	// TODO:          --private           Show only private repositories
	// TODO:          --public            Show only public repositories
	// TODO:          --source            Show only non-forks
	reposCommand.Flags().StringVarP(&reposFlags.query, "query", "", "", "Query for selecting projects")
	facadeCommand.AddCommand(reposCommand)
}
