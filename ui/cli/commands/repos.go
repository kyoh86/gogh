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
	var relationMap = map[string]hosting.RepositoryAffiliation{
		"owner":               hosting.RepositoryAffiliationOwner,
		"organization-member": hosting.RepositoryAffiliationOrganizationMember,
		"collaborator":        hosting.RepositoryAffiliationCollaborator,
	}
	var relationAccepts = []string{
		"owner",
		"organization-member",
		"collaborator",
	}
	var sortMap = map[string]hosting.RepositoryOrderField{
		"created-at": hosting.RepositoryOrderFieldCreatedAt,
		"name":       hosting.RepositoryOrderFieldName,
		"pushed-at":  hosting.RepositoryOrderFieldPushedAt,
		"stargazers": hosting.RepositoryOrderFieldStargazers,
		"updated-at": hosting.RepositoryOrderFieldUpdatedAt,
	}
	var sortAccepts = []string{
		"created-at",
		"name",
		"pushed-at",
		"stargazers",
		"updated-at",
	}
	var orderMap = map[string]hosting.OrderDirection{
		"asc":        hosting.OrderDirectionAsc,
		"ascending":  hosting.OrderDirectionAsc,
		"desc":       hosting.OrderDirectionDesc,
		"descending": hosting.OrderDirectionDesc,
	}
	var orderAccept = []string{
		"asc", "ascending",
		"desc", "descending",
	}

	// validateReposFlags validates the mutually exclusive flags and flag values
	validateReposFlags := func(f config.ReposFlags) error {
		// Check mutually exclusive flags
		if f.Private && f.Public {
			return errors.New("specify only one of `--private` or `--public`")
		}
		if f.Fork && f.NotFork {
			return errors.New("specify only one of `--fork` or `--no-fork`")
		}
		if f.Archived && f.NotArchived {
			return errors.New("specify only one of `--archived` or `--no-archived`")
		}

		// Validate relation values
		for _, r := range f.Relation {
			if _, exists := relationMap[r]; !exists {
				return fmt.Errorf("invalid relation %q; %s", r, fmt.Sprintf("it can accept %s", quoteEnums(relationAccepts)))
			}
		}
		// Validate sort field
		if _, exists := sortMap[f.Sort]; !exists {
			return fmt.Errorf("invalid sort field %q; %s", f.Sort, fmt.Sprintf("it can accept %s", quoteEnums(sortAccepts)))
		}
		// Validate order field
		if _, exists := orderMap[f.Order]; !exists {
			return fmt.Errorf("invalid order field %q; %s", f.Order, fmt.Sprintf("it can accept %s", quoteEnums(orderAccept)))
		}
		return nil
	}

	// preprocessReposFlags converts command flags to repository query options
	preprocessReposFlags := func(f config.ReposFlags) *repos.Options {
		opts := &repos.Options{}

		// Set limit
		switch f.Limit {
		case 0:
			opts.Limit = 30
		case -1:
			opts.Limit = 0 // no limit
		default:
			opts.Limit = f.Limit
		}

		// Set privacy filter
		if f.Private {
			opts.Privacy = hosting.RepositoryPrivacyPrivate
		}
		if f.Public {
			opts.Privacy = hosting.RepositoryPrivacyPublic
		}

		// Set fork and archive filters
		if f.Fork {
			opts.IsFork = hosting.BooleanFilterTrue
		}
		if f.NotFork {
			opts.IsFork = hosting.BooleanFilterFalse
		}
		if f.Archived {
			opts.IsArchived = hosting.BooleanFilterTrue
		}
		if f.NotArchived {
			opts.IsArchived = hosting.BooleanFilterFalse
		}

		// Set relation filters
		for _, r := range f.Relation {
			if field, exists := relationMap[r]; exists {
				opts.OwnerAffiliations = append(opts.OwnerAffiliations, field)
			}
		}

		// Set sort field
		if field, exists := sortMap[strings.ToLower(f.Sort)]; exists {
			opts.OrderBy.Field = field
		}

		// Set sort direction
		if field, exists := orderMap[strings.ToLower(f.Order)]; exists {
			opts.OrderBy.Direction = field
		}

		return opts
	}

	var (
		f config.ReposFlags
	)
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "List remote repositories",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Validate flags
			if err := validateReposFlags(f); err != nil {
				return fmt.Errorf("invalid flag configuration: %w", err)
			}

			// Convert flags to options
			opts := preprocessReposFlags(f)

			// Setup output formatter
			printer, err := f.Format.Formatter(os.Stdout)
			if err != nil {
				return fmt.Errorf("failed to create output formatter: %w", err)
			}
			defer printer.Close()

			useCase := repos.NewUseCase(svc.hostingService)
			for repo, err := range useCase.Execute(cmd.Context(), *opts) {
				if err != nil {
					return fmt.Errorf("failed to list repositories: %w", err)
				}
				printer.Print(*repo)
			}
			return nil
		},
	}

	cmd.Flags().
		IntVarP(&f.Limit, "limit", "", svc.flags.Repos.Limit, "Max number of repositories to list. -1 means unlimited")

	cmd.Flags().
		BoolVarP(&f.Public, "public", "", svc.flags.Repos.Public, "Show only public repositories")
	cmd.Flags().
		BoolVarP(&f.Private, "private", "", svc.flags.Repos.Private, "Show only private repositories")

	cmd.Flags().
		BoolVarP(&f.Fork, "fork", "", svc.flags.Repos.Fork, "Show only forks")
	cmd.Flags().
		BoolVarP(&f.NotFork, "no-fork", "", svc.flags.Repos.NotFork, "Omit forks")

	cmd.Flags().
		BoolVarP(&f.Archived, "archived", "", svc.flags.Repos.Archived, "Show only archived repositories")
	cmd.Flags().
		BoolVarP(&f.NotArchived, "no-archived", "", svc.flags.Repos.NotArchived, "Omit archived repositories")

	f.Format = svc.flags.Repos.Format
	cmd.Flags().
		VarP(&f.Format, "format", "", flags.RemoteRepoFormatShortUsage)
	if err := cmd.RegisterFlagCompletionFunc("format", flags.CompleteRemoteRepoFormat); err != nil {
		panic(err)
	}
	cmd.Flags().
		StringVarP(&f.Color, "color", "", svc.flags.Repos.Color, "Colorize the output; It can accept 'auto', 'always' or 'never'")
	if err := cmd.RegisterFlagCompletionFunc("color", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"auto", "always", "never"}, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	cmd.Flags().
		StringSliceVarP(&f.Relation, "relation", "", svc.flags.Repos.Relation, fmt.Sprintf("The relation of user to each repository; it can accept %s", quoteEnums(relationAccepts)))
	if err := cmd.RegisterFlagCompletionFunc("relation", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return relationAccepts, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}

	cmd.Flags().
		StringVarP(&f.Sort, "sort", "", svc.flags.Repos.Sort, fmt.Sprintf("Property by which repository be ordered; it can accept %s", quoteEnums(sortAccepts)))
	if err := cmd.RegisterFlagCompletionFunc("sort", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return sortAccepts, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	cmd.Flags().
		StringVarP(&f.Order, "order", "", svc.flags.Repos.Order, fmt.Sprintf("Directions in which to order a list of items when provided an `sort` flag; it can accept %s", quoteEnums(orderAccept)))
	if err := cmd.RegisterFlagCompletionFunc("order", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return orderAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		panic(err)
	}
	return cmd
}
