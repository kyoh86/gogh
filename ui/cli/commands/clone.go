package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v4/app/clone"
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/overlay_apply"
	"github.com/kyoh86/gogh/v4/app/overlay_find"
	"github.com/kyoh86/gogh/v4/app/repos"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/app/try_clone"
	"github.com/kyoh86/gogh/v4/typ"
	"github.com/kyoh86/gogh/v4/ui/cli/flags"
	"github.com/kyoh86/gogh/v4/ui/cli/view"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewCloneCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f config.CloneFlags
	cloneUseCase := clone.NewUseCase(svc.HostingService, svc.WorkspaceService, svc.OverlayStore, svc.ReferenceParser, svc.GitService)

	checkFlags := func(ctx context.Context, args []string) ([]string, error) {
		if len(args) != 0 {
			return args, nil
		}
		var opts []huh.Option[string]
		for repo, err := range repos.NewUseCase(svc.HostingService).Execute(ctx, repos.Options{}) {
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
		logger := log.FromContext(ctx)
		if f.DryRun {
			for _, ref := range refs {
				fmt.Printf("git clone %q\n", ref)
			}
			return nil
		}

		eg, egCtx := errgroup.WithContext(ctx)
		eg.SetLimit(5)
		for _, ref := range refs {
			eg.Go(func() error {
				return cloneUseCase.Execute(egCtx, ref, clone.Options{
					TryCloneOptions: try_clone.Options{
						Notify:  try_clone.RetryLimit(1, nil),
						Timeout: f.CloneRetryTimeout,
					},
				})
			})
		}
		if err := eg.Wait(); err != nil {
			return fmt.Errorf("cloning repositories: %w", err)
		}
		overlayFindUseCase := overlay_find.NewUseCase(
			svc.ReferenceParser,
			svc.OverlayStore,
		)
		overlayApplyUseCase := overlay_apply.NewUseCase(
			svc.WorkspaceService,
			svc.FinderService,
			svc.ReferenceParser,
			svc.OverlayStore,
		)
		var applied bool
		for _, ref := range refs {
			if f.DryRun {
				fmt.Printf("Apply overlay for %q\n", ref)
			}
			if err := view.ProcessWithConfirmation(
				ctx,
				typ.FilterE(overlayFindUseCase.Execute(ctx, ref), func(ov *overlay_find.Overlay) (bool, error) {
					return !ov.ForInit, nil
				}),
				func(ov *overlay_find.Overlay) string {
					return fmt.Sprintf("Apply overlay for %s (%s)", ref, ov.RelativePath)
				},
				func(ov *overlay_find.Overlay) error {
					if err := overlayApplyUseCase.Execute(ctx, ref, ov.RepoPattern, ov.ForInit, ov.RelativePath); err != nil {
						return err
					}
					applied = true
					return nil
				},
			); err != nil {
				if errors.Is(err, view.ErrQuit) {
					return nil
				}
				return err
			}
			if applied {
				logger.Infof("Applied overlay for %s", ref)
			}
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

	flags.BoolVarP(cmd, &f.DryRun, "dry-run", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	cmd.Flags().DurationVarP(&f.CloneRetryTimeout, "clone-retry-timeout", "t", svc.Flags.Clone.CloneRetryTimeout, "Timeout for each clone attempt")
	return cmd, nil
}
