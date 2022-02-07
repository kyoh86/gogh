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
	"github.com/wacul/ptr"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
)

type reposFlagsStruct struct {
	Format   string   `yaml:"format,omitempty"`
	Color    string   `yaml:"color,omitempty"`
	Sort     string   `yaml:"sort,omitempty"`
	Order    string   `yaml:"order,omitempty"`
	Relation []string `yaml:"relation,omitempty"`
	Limit    int      `yaml:"limit,omitempty"`
	Private  bool     `yaml:"private,omitempty"`
	Public   bool     `yaml:"public,omitempty"`
	Fork     bool     `yaml:"fork,omitempty"`
	NotFork  bool     `yaml:"notFork,omitempty"`
}

var (
	reposFlags   reposFlagsStruct
	reposCommand = &cobra.Command{
		Use:   "repos",
		Short: "List remote repositories",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			list, err := servers.List()
			if err != nil {
				return err
			}
			var listOption gogh.RemoteListOption
			switch reposFlags.Limit {
			case 0:
				listOption.Limit = ptr.Int(30)
			case -1:
				listOption.Limit = ptr.Int(0) // no limit
			default:
				listOption.Limit = &reposFlags.Limit
			}
			if reposFlags.Private && reposFlags.Public {
				return errors.New("specify only one of `--private` or `--public`")
			}
			if reposFlags.Private {
				listOption.Private = &reposFlags.Private // &true
			}
			if reposFlags.Public {
				listOption.Private = &reposFlags.Private // &false
			}

			if reposFlags.Fork && reposFlags.NotFork {
				return errors.New("specify only one of `--fork` or `--no-fork`")
			}
			if reposFlags.Fork {
				listOption.IsFork = &reposFlags.Fork // &true
			}
			if reposFlags.NotFork {
				listOption.IsFork = &reposFlags.Fork // &false
			}
		LOOP_CONVERT_RELATION:
			for _, r := range reposFlags.Relation {
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
			switch reposFlags.Format {
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
				if term.IsTerminal(int(os.Stdout.Fd())) || reposFlags.Color == "always" {
					options = append(options, repotab.Styled())
				}
				format = repotab.NewPrinter(os.Stdout, options...)
			default:
				return fmt.Errorf("invalid format %q; %s", reposFlags.Format, repoFormatAccept)
			}
			defer format.Close()
			if reposFlags.Sort != "" {
				sort := gogh.RepositoryOrderField(reposFlags.Sort)
				if !sort.IsValid() {
					return fmt.Errorf("invalid sort %q; %s", reposFlags.Sort, repoSortAccept)
				}
				listOption.Sort = sort
			}
			if reposFlags.Order != "" {
				order := gogh.OrderDirection(reposFlags.Order)
				if !order.IsValid() {
					return fmt.Errorf("invalid order %q; %s", reposFlags.Order, repoOrderAccept)
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
)

var (
	repoFormatAccept   string
	repoSortAccept     string
	repoOrderAccept    string
	repoRelationAccept string
)

func init() {
	setup()
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
	reposCommand.Flags().
		IntVarP(&reposFlags.Limit, "limit", "", defaultInt(defaultFlag.Repos.Limit, 30), "Max number of repositories to list. -1 means unlimited")

	reposCommand.Flags().
		BoolVarP(&reposFlags.Public, "public", "", defaultFlag.Repos.Public, "Show only public repositories")
	reposCommand.Flags().
		BoolVarP(&reposFlags.Private, "private", "", defaultFlag.Repos.Private, "Show only private repositories")

	reposCommand.Flags().
		BoolVarP(&reposFlags.Fork, "fork", "", defaultFlag.Repos.Fork, "Show only forks")
	reposCommand.Flags().
		BoolVarP(&reposFlags.NotFork, "no-fork", "", defaultFlag.Repos.NotFork, "Omit forks")

	reposCommand.Flags().
		StringVarP(&reposFlags.Format, "format", "", defaultString(defaultFlag.Repos.Format, "table"), "The formatting style for each repository; "+repoFormatAccept)
	reposCommand.Flags().
		StringVarP(&reposFlags.Color, "color", "", defaultString(defaultFlag.Repos.Color, "auto"), "Colorize the output; It can accept 'auto', 'always' or 'never'")

	reposCommand.Flags().
		StringSliceVarP(&reposFlags.Relation, "relation", "", defaultStringSlice(defaultFlag.Repos.Relation, []string{"owner", "organizationMember"}), "The relation of user to each repository; "+repoRelationAccept)

	reposCommand.Flags().
		StringVarP(&reposFlags.Sort, "sort", "", defaultFlag.Repos.Sort, "Property by which repository be ordered; "+repoSortAccept)
	reposCommand.Flags().
		StringVarP(&reposFlags.Order, "order", "", defaultFlag.Repos.Order, "Directions in which to order a list of items when provided an `sort` flag; "+repoOrderAccept)

	facadeCommand.AddCommand(reposCommand)
}
