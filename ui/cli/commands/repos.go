package commands

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/repos"
	"github.com/kyoh86/gogh/v3/app/service"
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

func NewReposCommand(ctx context.Context, svc *service.ServiceSet) *cobra.Command {
	var (
		opts   repos.Options
		format flags.RepositoryFormat
	)
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "List remote repositories",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Setup output formatter
			printer, err := format.Formatter(os.Stdout)
			if err != nil {
				return fmt.Errorf("failed to create output formatter: %w", err)
			}
			defer printer.Close()

			useCase := repos.NewUseCase(svc.HostingService)
			for repo, err := range useCase.Execute(cmd.Context(), opts) {
				if err != nil {
					return fmt.Errorf("failed to list repositories: %w", err)
				}
				printer.Print(*repo)
			}
			return nil
		},
	}

	cmd.Flags().
		IntVarP(&opts.Limit, "limit", "", svc.Flags.Repos.Limit, "Max number of repositories to list. -1 means unlimited")

	cmd.Flags().
		BoolVarP(&opts.Public, "public", "", svc.Flags.Repos.Public, "Show only public repositories")
	cmd.Flags().
		BoolVarP(&opts.Private, "private", "", svc.Flags.Repos.Private, "Show only private repositories")

	cmd.Flags().
		BoolVarP(&opts.Fork, "fork", "", svc.Flags.Repos.Fork, "Show only forks")
	cmd.Flags().
		BoolVarP(&opts.NotFork, "no-fork", "", svc.Flags.Repos.NotFork, "Omit forks")

	cmd.Flags().
		BoolVarP(&opts.Archived, "archived", "", svc.Flags.Repos.Archived, "Show only archived repositories")
	cmd.Flags().
		BoolVarP(&opts.NotArchived, "no-archived", "", svc.Flags.Repos.NotArchived, "Omit archived repositories")

	if err := flags.RepositoryFormatFlag(cmd, &format, svc.Flags.Repos.Format); err != nil {
		log.FromContext(ctx).WithError(err).Error("failed to init format flag")
	}
	cmd.Flags().
		StringVarP(&opts.Color, "color", "", svc.Flags.Repos.Color, "Colorize the output; It can accept 'auto', 'always' or 'never'")
	if err := cmd.RegisterFlagCompletionFunc("color", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"auto", "always", "never"}, cobra.ShellCompDirectiveDefault
	}); err != nil {
		log.FromContext(ctx).WithError(err).Error("failed to register completion function for color flag")
	}

	var relationAccepts = []string{
		"owner",
		"organization-member",
		"collaborator",
	}
	cmd.Flags().
		StringSliceVarP(&opts.Relation, "relation", "", svc.Flags.Repos.Relation, fmt.Sprintf("The relation of user to each repository; it can accept %s", quoteEnums(relationAccepts)))
	if err := cmd.RegisterFlagCompletionFunc("relation", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return relationAccepts, cobra.ShellCompDirectiveDefault
	}); err != nil {
		log.FromContext(ctx).WithError(err).Error("failed to register completion function for relation flag")
	}

	var sortAccepts = []string{
		"created-at",
		"name",
		"pushed-at",
		"stargazers",
		"updated-at",
	}
	cmd.Flags().
		StringVarP(&opts.Sort, "sort", "", svc.Flags.Repos.Sort, fmt.Sprintf("Property by which repository be ordered; it can accept %s", quoteEnums(sortAccepts)))
	if err := cmd.RegisterFlagCompletionFunc("sort", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return sortAccepts, cobra.ShellCompDirectiveDefault
	}); err != nil {
		log.FromContext(ctx).WithError(err).Error("failed to register completion function for sort flag")
	}

	var orderAccept = []string{
		"asc", "ascending",
		"desc", "descending",
	}
	cmd.Flags().
		StringVarP(&opts.Order, "order", "", svc.Flags.Repos.Order, fmt.Sprintf("Directions in which to order a list of items when provided an `sort` flag; it can accept %s", quoteEnums(orderAccept)))
	if err := cmd.RegisterFlagCompletionFunc("order", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return orderAccept, cobra.ShellCompDirectiveDefault
	}); err != nil {
		log.FromContext(ctx).WithError(err).Error("failed to register completion function for order flag")
	}
	return cmd
}
