package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/kyoh86/gogh/v2/view"
	"github.com/kyoh86/gogh/v2/view/repotab"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
)

type reposFlagsStruct struct {
	Format      string   `yaml:"format,omitempty"`
	Color       string   `yaml:"color,omitempty"`
	Sort        string   `yaml:"sort,omitempty"`
	Order       string   `yaml:"order,omitempty"`
	Relation    []string `yaml:"relation,omitempty"`
	Limit       int      `yaml:"limit,omitempty"`
	Private     bool     `yaml:"private,omitempty"`
	Public      bool     `yaml:"public,omitempty"`
	Fork        bool     `yaml:"fork,omitempty"`
	NotFork     bool     `yaml:"notFork,omitempty"`
	Archived    bool     `yaml:"archived,omitempty"`
	NotArchived bool     `yaml:"notArchived,omitempty"`
}

var (
	reposFlags   reposFlagsStruct
	reposCommand = &cobra.Command{
		Use:   "repos",
		Short: "List remote repositories",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			var listOption gogh.RemoteListOption
			switch reposFlags.Limit {
			case 0:
				listOption.Limit = 30
			case -1:
				listOption.Limit = 0 // no limit
			default:
				listOption.Limit = reposFlags.Limit
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

			if reposFlags.Archived && reposFlags.NotArchived {
				return errors.New("specify only one of `--archived` or `--no-archived`")
			}
			if reposFlags.Archived {
				listOption.IsArchived = &reposFlags.Archived // &true
			}
			if reposFlags.NotArchived {
				listOption.IsArchived = &reposFlags.Archived // &false
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
				return fmt.Errorf("invalid relation %q; %s", r, fmt.Sprintf("it can accept %s", quoteEnums(repoRelationAccept)))
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
				listOption.Sort = gogh.RepositoryOrderField(reposFlags.Sort)
			}
			if reposFlags.Order != "" {
				listOption.Order = gogh.OrderDirection(reposFlags.Order)
			}
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()
			eg, ctx := errgroup.WithContext(ctx)

			entries := tokens.Entries()
			if len(entries) == 0 {
				log.FromContext(ctx).Warn("No valid token found: you need to set token by `gogh auth login`")
				return nil
			}
			for _, entry := range entries {
				entry := entry
				eg.Go(func() error {
					adaptor, err := github.NewAdaptor(ctx, entry.Host, entry.Token)
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
	repoFormatAccept   []string
	repoSortAccept     []string
	repoOrderAccept    []string
	repoRelationAccept []string
)

func quoteEnums(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, v := range values {
		quoted = append(quoted, strconv.Quote(v))
	}
	return strings.Join(quoted[:len(quoted)-1], ", ") + " or " + quoted[len(quoted)-1]
}

func init() {
	repoFormatAccept = []string{"spec", "url", "json", "table"}
	for _, v := range gogh.AllRepositoryOrderField {
		repoSortAccept = append(repoSortAccept, string(v))
	}
	for _, v := range gogh.AllOrderDirection {
		repoOrderAccept = append(repoOrderAccept, string(v))
	}
	for _, v := range gogh.AllRepositoryRelation {
		repoRelationAccept = append(repoRelationAccept, v.String())
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
		BoolVarP(&reposFlags.Archived, "archived", "", defaultFlag.Repos.Archived, "Show only archived repositories")
	reposCommand.Flags().
		BoolVarP(&reposFlags.NotArchived, "no-archived", "", defaultFlag.Repos.NotArchived, "Omit archived repositories")

	reposCommand.Flags().
		StringVarP(&reposFlags.Format, "format", "", defaultString(defaultFlag.Repos.Format, "table"), fmt.Sprintf("The formatting style for each repository; it can accept %s", quoteEnums(repoFormatAccept)))
	if err := reposCommand.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return repoFormatAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	reposCommand.Flags().
		StringVarP(&reposFlags.Color, "color", "", defaultString(defaultFlag.Repos.Color, "auto"), "Colorize the output; It can accept 'auto', 'always' or 'never'")
	if err := reposCommand.RegisterFlagCompletionFunc("color", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"auto", "always", "never"}, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	reposCommand.Flags().
		StringSliceVarP(&reposFlags.Relation, "relation", "", defaultStringSlice(defaultFlag.Repos.Relation, []string{"owner", "organizationMember"}), fmt.Sprintf("The relation of user to each repository; it can accept %s", quoteEnums(repoRelationAccept)))
	if err := reposCommand.RegisterFlagCompletionFunc("relation", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return repoRelationAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}

	reposCommand.Flags().
		StringVarP(&reposFlags.Sort, "sort", "", defaultFlag.Repos.Sort, fmt.Sprintf("Property by which repository be ordered; it can accept %s", quoteEnums(repoSortAccept)))
	if err := reposCommand.RegisterFlagCompletionFunc("sort", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return repoSortAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	reposCommand.Flags().
		StringVarP(&reposFlags.Order, "order", "", defaultFlag.Repos.Order, fmt.Sprintf("Directions in which to order a list of items when provided an `sort` flag; it can accept %s", quoteEnums(repoOrderAccept)))
	if err := reposCommand.RegisterFlagCompletionFunc("order", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return repoOrderAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}

	facadeCommand.AddCommand(reposCommand)
}
