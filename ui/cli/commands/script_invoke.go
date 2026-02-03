package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/cwd"
	"github.com/kyoh86/gogh/v4/app/list"
	"github.com/kyoh86/gogh/v4/app/script/invoke"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewScriptInvokeCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		allRepositories bool
		patterns        []string
	}
	cmd := &cobra.Command{
		Use:   "invoke [flags] <script-id> [[[<host>/]<owner>/]<name>...]",
		Short: "Invoke an script in a repository",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
			if len(args) > 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return completeScripts(cmd.Context(), svc, toComplete)
		},
		Args: cobra.MinimumNArgs(1),
		Example: `  invoke [flags] <script-id> [[[<host>/]<owner>/]<name>...]
  invoke [flags] <script-id> --all
  invoke [flags] <script-id> --pattern <pattern> [--pattern <pattern>]...

  It accepts a short notation for each repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"; 
    - "<owner>/<name>": e.g. "kyoh86/example"
    - "." for the current directory repository
  They'll be completed with the default host and owner set by "config set-default{-host|-owner}".

  It also accepts an alias for each repository.
	The alias is a local name for the remote repository.
  For example:
    - "kyoh86/example=sample"
    - "kyoh86/example=kyoh86-tryouts/tryout"
  For each them will be cloned from "github.com/kyoh86/example" into the local as:
    - "$(gogh root)/github.com/kyoh86/sample"
    - "$(gogh root)/github.com/kyoh86-tryouts/tryout"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)
			scriptID := args[0]
			refs := args[1:]
			scriptInvokeUsecase := invoke.NewUsecase(
				svc.WorkspaceService,
				svc.FinderService,
				svc.ScriptService,
				svc.ReferenceParser,
			)
			if f.allRepositories || len(f.patterns) > 0 {
				if len(refs) > 0 {
					return errors.New("cannot specify repositories when --all or --pattern flag is set")
				}

				// If --all flag is set, apply the script to all repositories in the workspace
				for repo, err := range list.NewUsecase(
					svc.WorkspaceService,
					svc.FinderService,
				).Execute(ctx, list.Options{ListOptions: list.ListOptions{
					Limit:    0,
					Patterns: f.patterns,
				}}) {
					if err != nil {
						return fmt.Errorf("listing repositories: %w", err)
					}
					refs = append(refs, repo.Ref().String())
				}
				if len(refs) == 0 {
					logger := log.FromContext(ctx)
					if len(f.patterns) > 0 {
						logger = logger.WithField("patterns", strings.Join(f.patterns, "|"))
						logger.Info(strings.Join([]string{
							"No entry found.",
							"Patterns should be formed as <host>/<owner>/<name>.",
							`For example, to match any repository of "kyoh86", use "*/kyoh86/*"`,
						}, "\n"))
					} else {
						logger.Info("No entry found")
					}
				}
			}
			for _, ref := range refs {
				resolvedRef := ref

				// Use current directory if reference is "."
				if ref == "." {
					repo, err := cwd.NewUsecase(svc.WorkspaceService, svc.FinderService).Execute(ctx)
					if err != nil {
						return fmt.Errorf("finding repository from current directory: %w", err)
					}
					resolvedRef = repo.Ref().String()
				}

				if err := scriptInvokeUsecase.Execute(ctx, resolvedRef, scriptID, map[string]any{}); err != nil {
					return err
				}
				logger.Infof("Invoked script %s in %s", scriptID, ref)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&f.allRepositories, "all", "", false, "Apply to all repositories in the workspace")
	cmd.Flags().StringSliceVarP(&f.patterns, "pattern", "p", nil, "Patterns for selecting repositories")
	return cmd, nil
}
