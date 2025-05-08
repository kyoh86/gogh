package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kyoh86/gogh/v3/app/repos"
	"github.com/kyoh86/gogh/v3/core/auth"
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

func NewReposCommand(tokens auth.TokenService, hostingService hosting.HostingService, defaults *config.FlagStore) *cobra.Command {
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
	checkFlags := func(cmd *cobra.Command, args []string) (printer view.RemoteRepoPrinter, options *repos.Options, err error) {
		switch f.Limit {
		case 0:
			options.Limit = 30
		case -1:
			options.Limit = 0 // no limit
		default:
			options.Limit = f.Limit
		}
		if f.Private && f.Public {
			return nil, nil, errors.New("specify only one of `--private` or `--public`")
		}
		if f.Private {
			options.Privacy = hosting.RepositoryPrivacyPrivate
		}
		if f.Public {
			options.Privacy = hosting.RepositoryPrivacyPublic
		}

		if f.Fork && f.NotFork {
			return nil, nil, errors.New("specify only one of `--fork` or `--no-fork`")
		}
		if f.Fork {
			options.IsFork = hosting.BooleanFilterTrue
		}
		if f.NotFork {
			options.IsFork = hosting.BooleanFilterFalse
		}

		if f.Archived && f.NotArchived {
			return nil, nil, errors.New("specify only one of `--archived` or `--no-archived`")
		}
		if f.Archived {
			options.IsArchived = hosting.BooleanFilterTrue
		}
		if f.NotArchived {
			options.IsArchived = hosting.BooleanFilterFalse
		}
		for _, r := range f.Relation {
			switch r {
			case "owner":
				options.OwnerAffiliations = append(options.OwnerAffiliations, hosting.RepositoryAffiliationOwner)
			case "organizationMember", "organization-member", "organization_member":
				options.OwnerAffiliations = append(options.OwnerAffiliations, hosting.RepositoryAffiliationOrganizationMember)
			case "collaborator":
				options.OwnerAffiliations = append(options.OwnerAffiliations, hosting.RepositoryAffiliationCollaborator)
			default:
				return nil, nil, fmt.Errorf("invalid relation %q; %s", r, fmt.Sprintf("it can accept %s", quoteEnums(remoteRepoRelationAccept)))
			}
		}

		switch strings.ToLower(f.Sort) {
		case "created-at", "createdAt", "created_at":
			options.ListRepositoryOptions.OrderBy.Field = hosting.RepositoryOrderFieldCreatedAt
		case "name":
			options.ListRepositoryOptions.OrderBy.Field = hosting.RepositoryOrderFieldName
		case "pushed-at", "pushedAt", "pushed_at":
			options.ListRepositoryOptions.OrderBy.Field = hosting.RepositoryOrderFieldPushedAt
		case "stargazers":
			options.ListRepositoryOptions.OrderBy.Field = hosting.RepositoryOrderFieldStargazers
		case "updated-at", "updatedAt", "updated_at":
			options.ListRepositoryOptions.OrderBy.Field = hosting.RepositoryOrderFieldUpdatedAt
		}

		switch strings.ToLower(f.Order) {
		case "asc", "ascending":
			options.ListRepositoryOptions.OrderBy.Direction = hosting.OrderDirectionAsc
		case "desc", "descending":
			options.ListRepositoryOptions.OrderBy.Direction = hosting.OrderDirectionDesc
		}

		printer, err = f.Format.Formatter(os.Stdout)
		return printer, options, err
	}
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "List remote repositories",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			printer, options, err := checkFlags(cmd, args)
			if err != nil {
				return err
			}
			defer printer.Close()
			useCase := repos.NewUseCase(hostingService)
			for repo := range useCase.Execute(cmd.Context(), *options) {
				printer.Print(*repo)
			}
			return nil
		},
	}

	cmd.Flags().
		IntVarP(&f.Limit, "limit", "", DefaultValue(defaults.Repos.Limit, 30), "Max number of repositories to list. -1 means unlimited")

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
		StringVarP(&f.Color, "color", "", DefaultValue(defaults.Repos.Color, "auto"), "Colorize the output; It can accept 'auto', 'always' or 'never'")
	if err := cmd.RegisterFlagCompletionFunc("color", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"auto", "always", "never"}, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	cmd.Flags().
		StringSliceVarP(&f.Relation, "relation", "", DefaultSlice(defaults.Repos.Relation, []string{"owner", "organizationMember"}), fmt.Sprintf("The relation of user to each repository; it can accept %s", quoteEnums(remoteRepoRelationAccept)))
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
