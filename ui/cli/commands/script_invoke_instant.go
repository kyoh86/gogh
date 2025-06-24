package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/cwd"
	"github.com/kyoh86/gogh/v4/app/list"
	"github.com/kyoh86/gogh/v4/app/script/invoke"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
)

func NewScriptInvokeInstantCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		allRepositories bool
		patterns        []string
		file            string
	}
	cmd := &cobra.Command{
		Use:   "invoke-instant [flags] [[[<host>/]<owner>/]<name>...]",
		Short: "Run a temporary script in a repository without storing it",
		Args:  cobra.ArbitraryArgs,
		Example: `  invoke-instant --file script.lua repo1 repo2
  invoke-instant --file - repo1 < script.lua
  echo 'print(gogh.repo.name)' | gogh script invoke-instant --file - repo1
  invoke-instant --file script.lua .  # Use current directory repository
  invoke-instant --file script.lua --all
  invoke-instant --file script.lua --pattern <pattern>

  It accepts a short notation for each repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"; 
    - "<owner>/<name>": e.g. "kyoh86/example"
    - "." for the current directory repository
  They'll be completed with the default host and owner set by "config set-default{-host|-owner}".`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)
			refs := args

			// Validate flags
			if f.file == "" {
				return errors.New("must specify --file")
			}

			// Read script content
			var scriptContent []byte
			var err error
			if f.file == "-" {
				scriptContent, err = io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("reading from stdin: %w", err)
				}
			} else {
				scriptContent, err = os.ReadFile(f.file)
				if err != nil {
					return fmt.Errorf("reading script file: %w", err)
				}
			}

			// Get repository list
			if f.allRepositories || len(f.patterns) > 0 {
				if len(refs) > 0 {
					return errors.New("cannot specify repositories when --all or --pattern flag is set")
				}

				for repo, err := range list.NewUseCase(
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
			}

			// Execute script for each repository
			for _, ref := range refs {
				// Use current directory if reference is "."
				if ref == "." {
					repo, err := cwd.NewUseCase(svc.WorkspaceService, svc.FinderService).Execute(ctx)
					if err != nil {
						return fmt.Errorf("finding repository from current directory: %w", err)
					}
					ref = repo.Ref().String()
				}

				refWithAlias, err := svc.ReferenceParser.ParseWithAlias(ref)
				if err != nil {
					return fmt.Errorf("parsing repository reference: %w", err)
				}
				match, err := svc.FinderService.FindByReference(ctx, svc.WorkspaceService, refWithAlias.Local())
				if err != nil {
					return fmt.Errorf("find repository location: %w", err)
				}
				if match == nil {
					return fmt.Errorf("repository not found: %s", ref)
				}

				if err := invoke.InvokeInstant(ctx, match, string(scriptContent), map[string]any{}); err != nil {
					return fmt.Errorf("running script in %s: %w", ref, err)
				}
				logger.Infof("Ran instant script in %s", ref)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&f.allRepositories, "all", "", false, "Apply to all repositories in the workspace")
	cmd.Flags().StringSliceVarP(&f.patterns, "pattern", "p", nil, "Patterns for selecting repositories")
	cmd.Flags().StringVarP(&f.file, "file", "f", "", "Path to script file to invoke (use '-' for stdin)")
	if err := cmd.MarkFlagRequired("file"); err != nil {
		return nil, err
	}
	return cmd, nil
}
