package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/domain/remote"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/infra/github"
	"github.com/kyoh86/gogh/v3/ui/cli/flags"
	"github.com/kyoh86/gogh/v3/ui/cli/view"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func quoteEnums(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, v := range values {
		quoted = append(quoted, strconv.Quote(v))
	}
	return strings.Join(quoted[:len(quoted)-1], ", ") + " or " + quoted[len(quoted)-1]
}

func NewReposCommand(tokens *config.TokenStore, defaults *config.FlagStore) *cobra.Command {
	var (
		f                        config.ReposFlags
		remoteRepoSortAccept     []string
		remoteRepoOrderAccept    []string
		remoteRepoRelationAccept []string
	)
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "List remote repositories",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			var listOption remote.RemoteListOption
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
				listOption.Private = &f.Public // &false
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
				rdef := remote.RemoteRepoRelation(r)
				for _, def := range remote.AllRemoteRepoRelation {
					if def == rdef {
						listOption.Relation = append(listOption.Relation, rdef)
						continue LOOP_CONVERT_RELATION
					}
				}
				return fmt.Errorf("invalid relation %q; %s", r, fmt.Sprintf("it can accept %s", quoteEnums(remoteRepoRelationAccept)))
			}
			var format view.RemoteRepoPrinter
			var err error
			format, err = f.Format.Formatter(os.Stdout)
			if err != nil {
				return err
			}
			defer format.Close()

			if f.Sort != "" {
				listOption.Sort = remote.RemoteRepoOrderField(f.Sort)
			}
			if f.Order != "" {
				listOption.Order = remote.OrderDirection(f.Order)
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
					ctrl := remote.NewRemoteController(adaptor)
					rch, ech := ctrl.ListAsync(ctx, &listOption)
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

	for _, v := range remote.AllRemoteRepoOrderField {
		remoteRepoSortAccept = append(remoteRepoSortAccept, string(v))
	}
	for _, v := range remote.AllOrderDirection {
		remoteRepoOrderAccept = append(remoteRepoOrderAccept, string(v))
	}
	for _, v := range remote.AllRemoteRepoRelation {
		remoteRepoRelationAccept = append(remoteRepoRelationAccept, v.String())
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
		VarP(&f.Format, "format", "", flags.RemoteRepoFormatShortUsage)
	if err := cmd.RegisterFlagCompletionFunc("format", flags.CompleteRemoteRepoFormat); err != nil {
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
		StringSliceVarP(&f.Relation, "relation", "", defaultStringSlice(defaults.Repos.Relation, []string{"owner", "organizationMember"}), fmt.Sprintf("The relation of user to each repository; it can accept %s", quoteEnums(remoteRepoRelationAccept)))
	if err := cmd.RegisterFlagCompletionFunc("relation", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return remoteRepoRelationAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}

	cmd.Flags().
		StringVarP(&f.Sort, "sort", "", defaults.Repos.Sort, fmt.Sprintf("Property by which repository be ordered; it can accept %s", quoteEnums(remoteRepoSortAccept)))
	if err := cmd.RegisterFlagCompletionFunc("sort", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return remoteRepoSortAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	cmd.Flags().
		StringVarP(&f.Order, "order", "", defaults.Repos.Order, fmt.Sprintf("Directions in which to order a list of items when provided an `sort` flag; it can accept %s", quoteEnums(remoteRepoOrderAccept)))
	if err := cmd.RegisterFlagCompletionFunc("order", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return remoteRepoOrderAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	return cmd
}
