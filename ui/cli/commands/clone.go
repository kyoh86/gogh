package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v4/app/clone"
	"github.com/kyoh86/gogh/v4/app/clone/try"
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/repos"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewCloneCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f config.CloneFlags
	cloneUsecase := clone.NewUsecase(
		svc.HostingService,
		svc.WorkspaceService,
		svc.FinderService,
		svc.OverlayService,
		svc.ScriptService,
		svc.HookService,
		svc.ReferenceParser,
		svc.GitService,
	)

	checkFlags := func(ctx context.Context, args []string) ([]string, error) {
		if len(args) != 0 {
			return args, nil
		}
		var opts []huh.Option[string]
		for repo, err := range repos.NewUsecase(svc.HostingService).Execute(ctx, repos.Options{}) {
			if err != nil {
				return nil, fmt.Errorf("listing up repositories: %w", err)
			}
			opts = append(opts, huh.Option[string]{
				Key:   repo.Ref.String(),
				Value: repo.Ref.String(),
			})
		}
		if err := huh.NewForm(huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("A remote repository to clone").
				Options(opts...).
				Value(&args),
		)).Run(); err != nil {
			return nil, err
		}
		return args, nil
	}

	runFunc := func(ctx context.Context, refs []string) error {
		if f.DryRun {
			for _, ref := range refs {
				fmt.Printf("git clone %q\n", ref)
			}
			return nil
		}

		eg, egCtx := errgroup.WithContext(ctx)
		eg.SetLimit(5)
		for _, ref := range refs {
			ref := ref // capture loop variable
			eg.Go(func() error {
				err := cloneUsecase.Execute(egCtx, ref, clone.Options{
					TryCloneOptions: try.Options{
						Notify:  try.RetryLimit(1, nil),
						Timeout: f.CloneRetryTimeout,
					},
				})
				return err
			})
		}
		if err := eg.Wait(); err != nil {
			return fmt.Errorf("cloning repositories: %w", err)
		}
		return nil
	}

	cmd := &cobra.Command{
		Use:     "clone [flags] [[[<host>/]<owner>/]<name>[=<alias>]...]",
		Aliases: []string{"get"},
		Args:    cobra.ArbitraryArgs,
		Short:   "Clone remote repositories to local",
		Example: `  It accepts a short notation for a repository
  (for example, "github.com/kyoh86/example") like below.
    - "<name>": e.g. "example"; 
    - "<owner>/<name>": e.g. "kyoh86/example"
  They'll be completed with the default host and owner set by "config set-default{-host|-owner}".

  It also accepts an alias for each repository.
	The alias is used for a local repository.
  For example:
    - "kyoh86/example=sample"
    - "kyoh86/example=kyoh86-tryouts/tryout"
  For each them will be cloned from "github.com/kyoh86/example" into the local as:
    - "$(gogh root)/github.com/kyoh86/sample"
    - "$(gogh root)/github.com/kyoh86-tryouts/tryout"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			args, err := checkFlags(ctx, args)
			if err != nil {
				return err
			}
			if len(args) == 0 {
				return errors.New("no repository specified")
			}
			if err := runFunc(ctx, args); err != nil {
				return fmt.Errorf("cloning repositories: %v", err)
			}
			log.FromContext(ctx).Infof("Cloning %d repositories completed", len(args))
			return nil
		},
	}

	cmd.Flags().BoolVarP(&f.DryRun, "dry-run", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	cmd.Flags().DurationVarP(&f.CloneRetryTimeout, "clone-retry-timeout", "t", svc.Flags.Clone.CloneRetryTimeout, "Timeout for each clone attempt")
	return cmd, nil
}
