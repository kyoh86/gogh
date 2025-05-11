package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kyoh86/gogh/v3/app/repos"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/ui/cli/flags"
	"github.com/kyoh86/gogh/v3/ui/cli/view"
	"github.com/spf13/cobra"
)

func quoteEnums(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, v := range values {
		quoted = append(quoted, strconv.Quote(v))
	}
	return strings.Join(quoted[:len(quoted)-1], ", ") + " or " + quoted[len(quoted)-1]
}

func NewReposCommand(svc *ServiceSet) *cobra.Command {
	var (
		f config.ReposFlags
	)
	remoteRepoSortAccept := []string{
		"CREATED_AT",
		"created_at",
		"created-at",
		"createdAt",
		"NAME",
		"name",
		"PUSHED_AT",
		"pushed_at",
		"pushed-at",
		"pushedAt",
		"STARGAZERS",
		"stargazers",
		"UPDATED_AT",
		"updated_at",
		"updated-at",
		"updatedAt",
	}
	remoteRepoOrderAccept := []string{
		"asc", "ascending",
		"ASC", "ASCENDING",
		"desc", "descending",
		"DESC", "DESCENDING",
	}
	remoteRepoRelationAccept := []string{
		"owner",
		"organization-member",
		"organization_member",
		"organizationMember",
		"collaborator",
	}
	checkFlags := func() (printer view.RemoteRepoPrinter, opts *repos.Options, err error) {
		opts = &repos.Options{}
		switch f.Limit {
		case 0:
			opts.Limit = 30
		case -1:
			opts.Limit = 0 // no limit
		default:
			opts.Limit = f.Limit
		}
		if f.Private && f.Public {
			return nil, nil, errors.New("specify only one of `--private` or `--public`")
		}
		if f.Private {
			opts.Privacy = hosting.RepositoryPrivacyPrivate
		}
		if f.Public {
			opts.Privacy = hosting.RepositoryPrivacyPublic
		}

		if f.Fork && f.NotFork {
			return nil, nil, errors.New("specify only one of `--fork` or `--no-fork`")
		}
		if f.Fork {
			opts.IsFork = hosting.BooleanFilterTrue
		}
		if f.NotFork {
			opts.IsFork = hosting.BooleanFilterFalse
		}

		if f.Archived && f.NotArchived {
			return nil, nil, errors.New("specify only one of `--archived` or `--no-archived`")
		}
		if f.Archived {
			opts.IsArchived = hosting.BooleanFilterTrue
		}
		if f.NotArchived {
			opts.IsArchived = hosting.BooleanFilterFalse
		}
		for _, r := range f.Relation {
			switch r {
			case "owner":
				opts.OwnerAffiliations = append(opts.OwnerAffiliations, hosting.RepositoryAffiliationOwner)
			case "organizationMember", "organization-member", "organization_member":
				opts.OwnerAffiliations = append(opts.OwnerAffiliations, hosting.RepositoryAffiliationOrganizationMember)
			case "collaborator":
				opts.OwnerAffiliations = append(opts.OwnerAffiliations, hosting.RepositoryAffiliationCollaborator)
			default:
				return nil, nil, fmt.Errorf("invalid relation %q; %s", r, fmt.Sprintf("it can accept %s", quoteEnums(remoteRepoRelationAccept)))
			}
		}

		switch strings.ToLower(f.Sort) {
		case "created-at", "createdAt", "created_at":
			opts.OrderBy.Field = hosting.RepositoryOrderFieldCreatedAt
		case "name":
			opts.OrderBy.Field = hosting.RepositoryOrderFieldName
		case "pushed-at", "pushedAt", "pushed_at":
			opts.OrderBy.Field = hosting.RepositoryOrderFieldPushedAt
		case "stargazers":
			opts.OrderBy.Field = hosting.RepositoryOrderFieldStargazers
		case "updated-at", "updatedAt", "updated_at":
			opts.OrderBy.Field = hosting.RepositoryOrderFieldUpdatedAt
		}

		switch strings.ToLower(f.Order) {
		case "asc", "ascending":
			opts.OrderBy.Direction = hosting.OrderDirectionAsc
		case "desc", "descending":
			opts.OrderBy.Direction = hosting.OrderDirectionDesc
		}

		printer, err = f.Format.Formatter(os.Stdout)
		return printer, opts, err
	}
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "List remote repositories",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			printer, opts, err := checkFlags()
			if err != nil {
				return err
			}
			defer printer.Close()
			useCase := repos.NewUseCase(svc.hostingService)
			for repo := range useCase.Execute(cmd.Context(), *opts) {
				printer.Print(*repo)
			}
			return nil
		},
	}

	cmd.Flags().
		IntVarP(&f.Limit, "limit", "", DefaultValue(svc.defaults.Repos.Limit, 30), "Max number of repositories to list. -1 means unlimited")

	cmd.Flags().
		BoolVarP(&f.Public, "public", "", svc.defaults.Repos.Public, "Show only public repositories")
	cmd.Flags().
		BoolVarP(&f.Private, "private", "", svc.defaults.Repos.Private, "Show only private repositories")

	cmd.Flags().
		BoolVarP(&f.Fork, "fork", "", svc.defaults.Repos.Fork, "Show only forks")
	cmd.Flags().
		BoolVarP(&f.NotFork, "no-fork", "", svc.defaults.Repos.NotFork, "Omit forks")

	cmd.Flags().
		BoolVarP(&f.Archived, "archived", "", svc.defaults.Repos.Archived, "Show only archived repositories")
	cmd.Flags().
		BoolVarP(&f.NotArchived, "no-archived", "", svc.defaults.Repos.NotArchived, "Omit archived repositories")

	cmd.Flags().
		VarP(&f.Format, "format", "", flags.RemoteRepoFormatShortUsage)
	if err := cmd.RegisterFlagCompletionFunc("format", flags.CompleteRemoteRepoFormat); err != nil {
		panic(err)
	}
	cmd.Flags().
		StringVarP(&f.Color, "color", "", DefaultValue(svc.defaults.Repos.Color, "auto"), "Colorize the output; It can accept 'auto', 'always' or 'never'")
	if err := cmd.RegisterFlagCompletionFunc("color", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"auto", "always", "never"}, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	cmd.Flags().
		StringSliceVarP(&f.Relation, "relation", "", DefaultSlice(svc.defaults.Repos.Relation, []string{"owner", "organizationMember"}), fmt.Sprintf("The relation of user to each repository; it can accept %s", quoteEnums(remoteRepoRelationAccept)))
	if err := cmd.RegisterFlagCompletionFunc("relation", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return remoteRepoRelationAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}

	cmd.Flags().
		StringVarP(&f.Sort, "sort", "", svc.defaults.Repos.Sort, fmt.Sprintf("Property by which repository be ordered; it can accept %s", quoteEnums(remoteRepoSortAccept)))
	if err := cmd.RegisterFlagCompletionFunc("sort", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return remoteRepoSortAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	cmd.Flags().
		StringVarP(&f.Order, "order", "", svc.defaults.Repos.Order, fmt.Sprintf("Directions in which to order a list of items when provided an `sort` flag; it can accept %s", quoteEnums(remoteRepoOrderAccept)))
	if err := cmd.RegisterFlagCompletionFunc("order", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return remoteRepoOrderAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	return cmd
}
