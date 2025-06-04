package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/clone"
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/overlay_apply"
	"github.com/kyoh86/gogh/v4/app/overlay_find"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/app/try_clone"
	"github.com/kyoh86/gogh/v4/typ"
	"github.com/kyoh86/gogh/v4/ui/cli/flags"
	"github.com/kyoh86/gogh/v4/ui/cli/view"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewBundleRestoreCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f config.BundleRestoreFlags
	cloneUseCase := clone.NewUseCase(svc.HostingService, svc.WorkspaceService, svc.OverlayStore, svc.ReferenceParser, svc.GitService)

	runFunc := func(ctx context.Context) error {
		logger := log.FromContext(ctx)
		in := os.Stdin
		if f.File != "" && f.File != "-" {
			f, err := os.Open(f.File)
			if err != nil {
				return fmt.Errorf("opening file: %w", err)
			}
			defer f.Close()
			in = f
		}
		eg, egCtx := errgroup.WithContext(ctx)
		scan := bufio.NewScanner(in)
		var refs []string
		for scan.Scan() {
			ref := scan.Text()
			if f.DryRun {
				fmt.Printf("git clone %q\n", ref)
			} else {
				eg.Go(func() error {
					if err := cloneUseCase.Execute(egCtx, ref, clone.Options{
						TryCloneOptions: try_clone.Options{
							Notify: try_clone.RetryLimit(f.CloneRetryLimit, view.TryCloneNotify(egCtx, nil)),
						},
					}); err != nil {
						return fmt.Errorf("cloning %s: %w", ref, err)
					}
					refs = append(refs, ref)
					return nil
				})
			}
		}
		if err := eg.Wait(); err != nil {
			return fmt.Errorf("cloning repositories: %w", err)
		}

		overlayFindUseCase := overlay_find.NewUseCase(
			svc.WorkspaceService,
			svc.FinderService,
			svc.ReferenceParser,
			svc.OverlayStore,
		)
		overlayApplyUseCase := overlay_apply.NewUseCase(svc.OverlayStore)
		for _, ref := range refs {
			if f.DryRun {
				fmt.Printf("Apply overlay for %q\n", ref)
			}
			if err := view.ProcessWithConfirmation(
				ctx,
				typ.FilterE(overlayFindUseCase.Execute(ctx, ref), func(entry *overlay_find.OverlayEntry) (bool, error) {
					return !entry.ForInit, nil
				}),
				func(entry *overlay_find.OverlayEntry) string {
					return fmt.Sprintf("Apply overlay for %s (%s)", ref, entry.RelativePath)
				},
				func(entry *overlay_find.OverlayEntry) error {
					return overlayApplyUseCase.Execute(ctx, entry.Location.FullPath(), entry.RepoPattern, entry.ForInit, entry.RelativePath)
				},
			); err != nil {
				if errors.Is(err, view.ErrQuit) {
					return nil
				}
				return err
			}
			logger.Infof("Applied overlay for %s", ref)
		}
		return nil
	}

	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Get dumped local repositoiries",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runFunc(cmd.Context())
		},
	}
	flags.BoolVarP(cmd, &f.DryRun, "dry-run", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	cmd.Flags().StringVarP(&f.File, "file", "f", svc.Flags.BundleRestore.File, `Read the file as input; if it's empty("") or hyphen("-"), read from stdin`)
	cmd.Flags().DurationVarP(&f.CloneRetryTimeout, "clone-retry-timeout", "", svc.Flags.BundleRestore.CloneRetryTimeout, "Timeout for each clone attempt")
	cmd.Flags().IntVarP(&f.CloneRetryLimit, "clone-retry-limit", "", svc.Flags.Create.CloneRetryLimit, "The number of retries to clone a repository")
	return cmd, nil
}
