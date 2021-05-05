package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/app"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/kyoh86/gogh/v2/view"
	"github.com/kyoh86/gogh/v2/view/repotab"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
)

var reposFlags struct {
	query    string
	format   string
	color    string
	relation []string
	limit    int
	private  bool
	public   bool
	fork     bool
	notFork  bool
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
		listOption := gogh.RemoteListOption{
			Query: reposFlags.query,
			Limit: &reposFlags.limit,
		}
		if reposFlags.private && reposFlags.public {
			return errors.New("specify only one of `--private` or `--public`")
		}
		if reposFlags.private {
			listOption.Private = &reposFlags.private // &true
		}
		if reposFlags.public {
			listOption.Private = &reposFlags.private // &false
		}

		if reposFlags.fork && reposFlags.notFork {
			return errors.New("specify only one of `--fork` or `--no-fork`")
		}
		if reposFlags.fork {
			listOption.IsFork = &reposFlags.fork // &true
		}
		if reposFlags.notFork {
			listOption.IsFork = &reposFlags.fork // &false
		}
		for _, r := range reposFlags.relation {
			rdef := gogh.RepositoryRelation(r)
			if !rdef.IsValid() {
				return errors.New("--relation can accept `owner`, `organizationMember` or `collaborator`")
			}
			listOption.Relation = append(listOption.Relation, rdef)
		}
		var options []repotab.Option
		if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			options = append(options, repotab.Width(width))
		}
		if term.IsTerminal(int(os.Stdout.Fd())) || reposFlags.color == "always" {
			options = append(options, repotab.Styled())
		}
		var format view.RepositoryPrinter
		switch reposFlags.format {
		case "spec":
			format = view.NewRepositorySpecPrinter(os.Stdout)
		case "url":
			format = view.NewRepositoryURLPrinter(os.Stdout)
		case "json":
			format = view.NewRepositoryJSONPrinter(os.Stdout)
		case "table":
			format = repotab.NewPrinter(os.Stdout, options...)
		default:
			return fmt.Errorf("invalid format %q; it can accept %q, %q, %q or %q", reposFlags.format, "spec", "url", "json", "table")
		}
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
				rch, ech := remote.ListAsync(ctx, &listOption)
				for {
					select {
					case repo, more := <-rch:
						if !more {
							return nil
						}
						format.Print(repo)
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
	reposCommand.Flags().IntVarP(&reposFlags.limit, "limit", "", 30, "Max number of repositories to list")

	reposCommand.Flags().BoolVarP(&reposFlags.public, "public", "", false, "Show only public repositories")
	reposCommand.Flags().BoolVarP(&reposFlags.private, "private", "", false, "Show only private repositories")

	reposCommand.Flags().BoolVarP(&reposFlags.fork, "fork", "", false, "Show only forks")
	reposCommand.Flags().BoolVarP(&reposFlags.notFork, "no-fork", "", false, "Omit forks")

	reposCommand.Flags().StringVarP(&reposFlags.query, "query", "", "", "Query for selecting projects")

	reposCommand.Flags().StringVarP(&reposFlags.format, "format", "", "table", "The formatting style for each repository")
	reposCommand.Flags().StringVarP(&reposFlags.color, "color", "", "auto", "Colorize the output; It can accept 'auto', 'always' or 'never'")

	reposCommand.Flags().StringSliceVarP(&reposFlags.relation, "relation", "", []string{"owner", "organizationMember"}, "The relation of user to each repository; It can accept `owner`, `organizationMember` or `collaborator`")

	// TODO: order
	facadeCommand.AddCommand(reposCommand)
}
