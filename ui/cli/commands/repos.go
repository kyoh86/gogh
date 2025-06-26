package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/repos"
	"github.com/kyoh86/gogh/v4/app/repository_print"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/ui/cli/flags"
	"github.com/spf13/cobra"
)

// quoteEnums takes a slice of strings and returns a string that lists the
// values in a human-readable format, with the last value preceded by "or".
func quoteEnums(values []string) string {
	if len(values) == 0 {
		return ""
	}
	quoted := make([]string, 0, len(values))
	for _, v := range values {
		quoted = append(quoted, strconv.Quote(v))
	}
	return strings.Join(quoted[:len(quoted)-1], ", ") + " or " + quoted[len(quoted)-1]
}

// enumFlag registers a string flag with a set of accepted values and adds
// completion for those values.
func enumFlag(cmd *cobra.Command, v *string, name string, defaultValue string, description string, accepts ...string) error {
	cmd.Flags().StringVarP(
		v,
		name,
		"",
		defaultValue,
		fmt.Sprintf("%s; it can accept %s", description, quoteEnums(accepts)),
	)
	return cmd.RegisterFlagCompletionFunc(
		name,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return accepts, cobra.ShellCompDirectiveDefault
		},
	)
}

func NewReposCommand(ctx context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var (
		opts   config.ReposFlags
		format string
	)
	cmd := &cobra.Command{
		Use:   "repos",
		Short: "List remote repositories",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			if err := repository_print.NewUsecase(cmd.OutOrStdout(), format).Execute(ctx, repos.NewUsecase(svc.HostingService).Execute(ctx, repos.Options{
				Limit:    opts.Limit,
				Privacy:  opts.Privacy,
				Fork:     opts.Fork,
				Archive:  opts.Archive,
				Format:   opts.Format,
				Color:    opts.Color,
				Relation: opts.Relation,
				Sort:     opts.Sort,
				Order:    opts.Order,
			})); err != nil {
				return fmt.Errorf("listing up repositories: %w", err)
			}
			return nil
		},
	}

	defs := svc.Flags.Repos

	cmd.Flags().IntVarP(&opts.Limit, "limit", "", defs.Limit, "Max number of repositories to list. -1 means unlimited")

	if err := flags.RepositoryFormatFlag(cmd, &format, defs.Format); err != nil {
		return nil, fmt.Errorf("initializing format flag: %w", err)
	}

	relationAccepts := []string{
		"owner",
		"organization-member",
		"collaborator",
	}
	cmd.Flags().StringSliceVarP(&opts.Relation, "relation", "", defs.Relation, fmt.Sprintf("The relation of user to each repository; it can accept %s", quoteEnums(relationAccepts)))
	if err := cmd.RegisterFlagCompletionFunc("relation", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return relationAccepts, cobra.ShellCompDirectiveDefault
	}); err != nil {
		return nil, fmt.Errorf("registering relation flag: %w", err)
	}

	if err := enumFlag(cmd, &opts.Privacy, "privacy", defs.Privacy, "Show only public/private repositories", "private", "public"); err != nil {
		return nil, fmt.Errorf("registering privacy flag: %w", err)
	}

	if err := enumFlag(cmd, &opts.Fork, "fork", defs.Fork, "Show only forked/not-forked repositories", "forked", "not-forked"); err != nil {
		return nil, fmt.Errorf("registering fork flag: %w", err)
	}

	if err := enumFlag(cmd, &opts.Archive, "archive", defs.Archive, "Show only archived/not-archived repositories", "archived", "not-archived"); err != nil {
		return nil, fmt.Errorf("registering archive flag: %w", err)
	}

	if err := enumFlag(cmd, &opts.Color, "color", defs.Color, "Colorize the output", "auto", "always", "never"); err != nil {
		return nil, fmt.Errorf("registering color flag: %w", err)
	}

	if err := enumFlag(cmd, &opts.Sort, "sort", defs.Sort, "Property by which repository be ordered", "created-at", "name", "pushed-at", "stargazers", "updated-at"); err != nil {
		return nil, fmt.Errorf("registering sort flag: %w", err)
	}

	if err := enumFlag(cmd, &opts.Order, "order", defs.Order, "Directions in which to order a list of items when provided a `sort` flag", "asc", "ascending", "desc", "descending"); err != nil {
		return nil, fmt.Errorf("registering order flag: %w", err)
	}

	return cmd, nil
}
