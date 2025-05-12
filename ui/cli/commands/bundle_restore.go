package commands

import (
	"bufio"
	"context"
	"os"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/clone"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/ui/cli/view"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewBundleRestoreCommand(svc *ServiceSet) *cobra.Command {
	var f config.BundleRestoreFlags
	cloneUseCase := clone.NewUseCase(svc.hostingService, svc.workspaceService, svc.gitService)

	runFunc := func(ctx context.Context) error {
		in := os.Stdin
		if f.File.Expand() != "" {
			f, err := os.Open(f.File.Expand())
			if err != nil {
				return err
			}
			defer f.Close()
			in = f
		}
		eg, ctx := errgroup.WithContext(ctx)
		scan := bufio.NewScanner(in)
		for scan.Scan() {
			ref, err := svc.referenceParser.ParseWithAlias(scan.Text())
			if err != nil {
				return err
			}
			if f.Dryrun {
				if ref.Alias == nil {
					log.FromContext(ctx).Infof("git clone %q", ref.Reference)
				} else {
					log.FromContext(ctx).Infof("git clone %q into %q", ref.Reference, ref.Alias)
				}
			} else {
				eg.Go(func() error {
					return cloneUseCase.Execute(ctx, ref.Reference, clone.Options{
						Alias:          ref.Alias,
						TryCloneNotify: service.RetryLimit(f.CloneRetryLimit, view.TryCloneNotify(ctx, nil)),
					})
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
	cmd.Flags().
		BoolVarP(&f.Dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	f.File = svc.flags.BundleRestore.File
	cmd.Flags().
		VarP(&f.File, "file", "", "Read the file as input; if not specified, read from stdin")
	cmd.Flags().
		IntVarP(&f.CloneRetryLimit, "clone-retry-limit", "", svc.flags.Create.CloneRetryLimit, "")
	return cmd
}
