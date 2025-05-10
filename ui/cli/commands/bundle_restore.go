package commands

import (
	"bufio"
	"context"
	"os"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/clone"
	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewBundleRestoreCommand(
	defaultNameService repository.DefaultNameService,
	tokenService auth.TokenService,
	defaults *config.FlagStore,
	hostingService hosting.HostingService,
	workspaceService workspace.WorkspaceService,
) *cobra.Command {
	var f config.BundleRestoreFlags
	cloneUseCase := clone.NewUseCase(hostingService, workspaceService)
	parser := repository.NewReferenceParser(defaultNameService.GetDefaultHostAndOwner())

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
			ref, err := parser.ParseWithAlias(scan.Text())
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
						Alias:      ref.Alias,
						RetryLimit: f.CloneRetryLimit,
					})
				})
			}
		}
		return eg.Wait()
	}

	cmd := &cobra.Command{
		Use:   "restore",
		Short: "Get dumped local repositoiries",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runFunc(cmd.Context())
		},
	}
	cmd.Flags().
		BoolVarP(&f.Dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	f.File = defaults.BundleRestore.File
	cmd.Flags().
		VarP(&f.File, "file", "", "Read the file as input; if not specified, read from stdin")
	cmd.Flags().
		IntVarP(&f.CloneRetryLimit, "clone-retry-limit", "", defaults.Create.CloneRetryLimit, "")
	return cmd
}
