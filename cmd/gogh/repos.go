package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/kyoh86/gogh/v2/view"
	"github.com/kyoh86/gogh/v2/view/repotab"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
)

var reposFlags struct {
	format   string
	color    string
	relation []string
	limit    int
	private  bool
	public   bool
	fork     bool
	notFork  bool
	sort     string
	order    string
}

var reposCommand = &cobra.Command{
	Use:   "repos",
	Short: "List remote repositories",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, _ []string) error {
		list, err := Servers().List()
		if err != nil {
			return err
		}
		listOption := gogh.RemoteListOption{
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
	LOOP_CONVERT_RELATION:
		for _, r := range reposFlags.relation {
			rdef := gogh.RepositoryRelation(r)
			for _, def := range gogh.AllRepositoryRelation {
				if def == rdef {
					listOption.Relation = append(listOption.Relation, rdef)
					continue LOOP_CONVERT_RELATION
				}
			}
			return fmt.Errorf("invalid relation %q; %s", r, repoRelationAccept)
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
			var options []repotab.Option
			if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
				options = append(options, repotab.Width(width))
			}
			if term.IsTerminal(int(os.Stdout.Fd())) || reposFlags.color == "always" {
				options = append(options, repotab.Styled())
			}
			format = repotab.NewPrinter(os.Stdout, options...)
		default:
			return fmt.Errorf("invalid format %q; %s", reposFlags.format, repoFormatAccept)
		}
		defer format.Close()
		if reposFlags.sort != "" {
			sort := gogh.RepositoryOrderField(reposFlags.sort)
			if !sort.IsValid() {
				return fmt.Errorf("invalid sort %q; %s", reposFlags.sort, repoSortAccept)
			}
			listOption.Sort = sort
		}
		if reposFlags.order != "" {
			order := gogh.OrderDirection(reposFlags.order)
			if !order.IsValid() {
				return fmt.Errorf("invalid order %q; %s", reposFlags.order, repoOrderAccept)
			}
			listOption.Order = order
		}
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

var (
	repoFormatAccept   string
	repoSortAccept     string
	repoOrderAccept    string
	repoRelationAccept string
)

func init() {
	repoFormatAccept = fmt.Sprintf("it can accept %q, %q, %q or %q", "spec", "url", "json", "table")
	{
		var valids []string
		for _, v := range gogh.AllRepositoryOrderField {
			valids = append(valids, strconv.Quote(v.String()))
		}
		repoSortAccept = fmt.Sprintf("it can accept %s", strings.Join(valids, ", "))
	}
	{
		var valids []string
		for _, v := range gogh.AllOrderDirection {
			valids = append(valids, strconv.Quote(v.String()))
		}
		repoOrderAccept = fmt.Sprintf("it can accept %s", strings.Join(valids, ", "))
	}
	{
		var valids []string
		for _, v := range gogh.AllRepositoryRelation {
			valids = append(valids, strconv.Quote(string(v)))
		}
		repoRelationAccept = fmt.Sprintf("it can accept %s", strings.Join(valids, ", "))
	}
	reposCommand.Flags().IntVarP(&reposFlags.limit, "limit", "", 30, "Max number of repositories to list. 0 means unlimited")

	reposCommand.Flags().BoolVarP(&reposFlags.public, "public", "", false, "Show only public repositories")
	reposCommand.Flags().BoolVarP(&reposFlags.private, "private", "", false, "Show only private repositories")

	reposCommand.Flags().BoolVarP(&reposFlags.fork, "fork", "", false, "Show only forks")
	reposCommand.Flags().BoolVarP(&reposFlags.notFork, "no-fork", "", false, "Omit forks")

	reposCommand.Flags().StringVarP(&reposFlags.format, "format", "", "table", "The formatting style for each repository; "+repoFormatAccept)
	reposCommand.Flags().StringVarP(&reposFlags.color, "color", "", "auto", "Colorize the output; It can accept 'auto', 'always' or 'never'")

	reposCommand.Flags().StringSliceVarP(&reposFlags.relation, "relation", "", []string{"owner", "organizationMember"}, "The relation of user to each repository; "+repoRelationAccept)

	reposCommand.Flags().StringVarP(&reposFlags.sort, "sort", "", "", "Property by which repository be ordered; "+repoSortAccept)
	reposCommand.Flags().StringVarP(&reposFlags.order, "order", "", "", "Directions in which to order a list of items when provided an `sort` flag; "+repoOrderAccept)

	facadeCommand.AddCommand(reposCommand)
}
