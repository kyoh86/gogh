package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3"
	"github.com/kyoh86/gogh/v3/config"
	"github.com/kyoh86/gogh/v3/infra/github"
	"github.com/kyoh86/gogh/v3/view"
	"github.com/kyoh86/gogh/v3/view/repotab"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"golang.org/x/term"
)

func quoteEnums(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, v := range values {
		quoted = append(quoted, strconv.Quote(v))
	}
	return strings.Join(quoted[:len(quoted)-1], ", ") + " or " + quoted[len(quoted)-1]
}

func NewReposCommand(tokens *config.TokenManager, defaults *config.Flags) *cobra.Command {
	var (
		f                  config.ReposFlags
		repoFormatAccept   []string
		repoSortAccept     []string
		repoOrderAccept    []string
		repoRelationAccept []string
	)
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "List remote repositories",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			var listOption gogh.RemoteListOption
			switch f.Limit {
			case 0:
				listOption.Limit = 30
			case -1:
				listOption.Limit = 0 // no limit
			default:
				listOption.Limit = f.Limit
			}
			if f.Private && f.Public {
				return errors.New("specify only one of `--private` or `--public`")
			}
			if f.Private {
				listOption.Private = &f.Private // &true
			}
			if f.Public {
				listOption.Private = &f.Private // &false
			}

			if f.Fork && f.NotFork {
				return errors.New("specify only one of `--fork` or `--no-fork`")
			}
			if f.Fork {
				listOption.IsFork = &f.Fork // &true
			}
			if f.NotFork {
				listOption.IsFork = &f.Fork // &false
			}

			if f.Archived && f.NotArchived {
				return errors.New("specify only one of `--archived` or `--no-archived`")
			}
			if f.Archived {
				listOption.IsArchived = &f.Archived // &true
			}
			if f.NotArchived {
				listOption.IsArchived = &f.Archived // &false
			}
		LOOP_CONVERT_RELATION:
			for _, r := range f.Relation {
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
			switch f.Format {
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
				if term.IsTerminal(int(os.Stdout.Fd())) || f.Color == "always" {
					options = append(options, repotab.Styled())
				}
				format = repotab.NewPrinter(os.Stdout, options...)
			default:
				return fmt.Errorf("invalid format %q; %s", f.Format, repoFormatAccept)
			}
			defer format.Close()
			if f.Sort != "" {
				listOption.Sort = gogh.RepositoryOrderField(f.Sort)
			}
			if f.Order != "" {
				listOption.Order = gogh.OrderDirection(f.Order)
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
				eg.Go(func() error {
					adaptor, err := github.NewAdaptor(ctx, entry.Host, &entry.Token)
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
	cmd.Flags().
		IntVarP(&f.Limit, "limit", "", defaultInt(defaults.Repos.Limit, 30), "Max number of repositories to list. -1 means unlimited")

	cmd.Flags().
		BoolVarP(&f.Public, "public", "", defaults.Repos.Public, "Show only public repositories")
	cmd.Flags().
		BoolVarP(&f.Private, "private", "", defaults.Repos.Private, "Show only private repositories")

	cmd.Flags().
		BoolVarP(&f.Fork, "fork", "", defaults.Repos.Fork, "Show only forks")
	cmd.Flags().
		BoolVarP(&f.NotFork, "no-fork", "", defaults.Repos.NotFork, "Omit forks")

	cmd.Flags().
		BoolVarP(&f.Archived, "archived", "", defaults.Repos.Archived, "Show only archived repositories")
	cmd.Flags().
		BoolVarP(&f.NotArchived, "no-archived", "", defaults.Repos.NotArchived, "Omit archived repositories")

	cmd.Flags().
		StringVarP(&f.Format, "format", "", defaultString(defaults.Repos.Format, "table"), fmt.Sprintf("The formatting style for each repository; it can accept %s", quoteEnums(repoFormatAccept)))
	if err := cmd.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return repoFormatAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	cmd.Flags().
		StringVarP(&f.Color, "color", "", defaultString(defaults.Repos.Color, "auto"), "Colorize the output; It can accept 'auto', 'always' or 'never'")
	if err := cmd.RegisterFlagCompletionFunc("color", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"auto", "always", "never"}, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	cmd.Flags().
		StringSliceVarP(&f.Relation, "relation", "", defaultStringSlice(defaults.Repos.Relation, []string{"owner", "organizationMember"}), fmt.Sprintf("The relation of user to each repository; it can accept %s", quoteEnums(repoRelationAccept)))
	if err := cmd.RegisterFlagCompletionFunc("relation", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return repoRelationAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}

	cmd.Flags().
		StringVarP(&f.Sort, "sort", "", defaults.Repos.Sort, fmt.Sprintf("Property by which repository be ordered; it can accept %s", quoteEnums(repoSortAccept)))
	if err := cmd.RegisterFlagCompletionFunc("sort", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return repoSortAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	cmd.Flags().
		StringVarP(&f.Order, "order", "", defaults.Repos.Order, fmt.Sprintf("Directions in which to order a list of items when provided an `sort` flag; it can accept %s", quoteEnums(repoOrderAccept)))
	if err := cmd.RegisterFlagCompletionFunc("order", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return repoOrderAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	return cmd
}
