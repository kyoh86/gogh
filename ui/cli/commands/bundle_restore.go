package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v4/app/clone"
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/app/try_clone"
	"github.com/kyoh86/gogh/v4/ui/cli/flags"
	"github.com/kyoh86/gogh/v4/ui/cli/view"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewBundleRestoreCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f config.BundleRestoreFlags
	cloneUseCase := clone.NewUseCase(svc.HostingService, svc.WorkspaceService, svc.OverlayService, svc.ReferenceParser, svc.GitService)

	runFunc := func(ctx context.Context) error {
		in := os.Stdin
		if f.File != "" {
			f, err := os.Open(f.File)
			if err != nil {
				return fmt.Errorf("opening file: %w", err)
			}
			defer f.Close()
			in = f
		}
		eg, egCtx := errgroup.WithContext(ctx)
		scan := bufio.NewScanner(in)
		for scan.Scan() {
			ref := scan.Text()
			if f.Dryrun {
				fmt.Printf("git clone %q\n", ref)
			} else {
				eg.Go(func() error {
					return cloneUseCase.Execute(egCtx, ref, clone.Options{
						TryCloneOptions: try_clone.Options{
							Notify: try_clone.RetryLimit(f.CloneRetryLimit, view.TryCloneNotify(egCtx, nil)),
						},
					})
					// TODO: Apply overlays
					// see: ./overlay_apply.go
				})
			}
		}
		return eg.Wait()
	}

	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Get dumped local repositoiries",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runFunc(cmd.Context())
		},
	}
	flags.BoolVarP(cmd, &f.Dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	cmd.Flags().StringVarP(&f.File, "file", "f", svc.Flags.BundleRestore.File, "Read the file as input; if not specified, read from stdin")
	cmd.Flags().DurationVarP(&f.CloneRetryTimeout, "clone-retry-timeout", "", svc.Flags.BundleRestore.CloneRetryTimeout, "Timeout for each clone attempt.")
	cmd.Flags().IntVarP(&f.CloneRetryLimit, "clone-retry-limit", "", svc.Flags.Create.CloneRetryLimit, "The number of retries to clone a repository")
	return cmd, nil
}
