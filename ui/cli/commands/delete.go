package commands

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v4/app/delete"
	"github.com/kyoh86/gogh/v4/app/repos"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/ui/cli/flags"
	"github.com/spf13/cobra"
)

func NewDeleteCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		local  bool
		remote bool
		force  bool
		dryrun bool
	}

	checkFlags := func(ctx context.Context, args []string) (string, error) {
		if len(args) != 0 {
			return args[0], nil
		}
		var opts []huh.Option[string]
		for repo, err := range repos.NewUseCase(svc.HostingService).Execute(ctx, repos.Options{}) {
			if err != nil {
				return "", fmt.Errorf("listing up repositories: %w", err)
			}
			opts = append(opts, huh.Option[string]{
				Key:   repo.Ref.String(),
				Value: repo.Ref.String(),
			})
		}
		var selected string
		if err := huh.NewForm(huh.NewGroup(
			huh.NewSelect[string]().
				Title("A repository to delete").
				Options(opts...).
				Value(&selected),
		)).Run(); err != nil {
			return "", err
		}
		return selected, nil
	}

	prepareFlags := func(_ context.Context, arg string) error {
		if !f.force {
			if f.local {
				var confirmed bool
				if err := huh.NewForm(huh.NewGroup(
					huh.NewConfirm().
						Title(fmt.Sprintf("Are you sure you want to delete local repository %s?", arg)).
						Value(&confirmed),
				)).Run(); err != nil {
					return err
				}
				if !confirmed {
					f.local = false
				}
			}
			if f.remote {
				var confirmed bool
				if err := huh.NewForm(huh.NewGroup(
					huh.NewConfirm().
						Title(fmt.Sprintf("Are you sure you want to delete remote-repository %s?", arg)).
						Value(&confirmed),
				)).Run(); err != nil {
					return err
				}
				if !confirmed {
					f.remote = false
				}
			}
		}
		return nil
	}

	cmd := &cobra.Command{
		Use:     "delete [flags] [[<owner>/]<name>]",
		Aliases: []string{"remove", "rm", "del"},
		Short:   "Delete local and remote repository",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, refs []string) error {
			ctx := cmd.Context()
			logger := log.FromContext(ctx)
			selected, err := checkFlags(ctx, refs)
			if err != nil {
				return err
			}

			if err := prepareFlags(ctx, selected); err != nil {
				return err
			}

			if f.local {
				logger.Infof("Delete local %s\n", selected)
			}
			if f.remote {
				logger.Infof("Deleting remote %s\n", selected)
			}
			if f.dryrun {
				return nil
			}
			if err := delete.NewUseCase(
				svc.WorkspaceService,
				svc.FinderService,
				svc.HostingService,
				svc.ReferenceParser,
			).Execute(ctx, selected, delete.Options{
				Local:  f.local,
				Remote: f.remote,
			}); err != nil {
				return fmt.Errorf("deleting the repository: %w", err)
			}
			if f.local {
				logger.Infof("Deleted local %s", selected)
			}
			if f.remote {
				logger.Infof("Deleted remote %s", selected)
			}
			return nil
		},
	}
	flags.BoolVarP(cmd, &f.local, "local", "", true, "Delete local repository.")
	cmd.Flags().Lookup("local").NoOptDefVal = "false"
	flags.BoolVarP(cmd, &f.remote, "remote", "", false, "Delete remote repository.")
	cmd.Flags().Lookup("remote").NoOptDefVal = "false"
	flags.BoolVarP(cmd, &f.force, "force", "", false, "Do NOT confirm to delete.")
	cmd.Flags().Lookup("force").NoOptDefVal = "false"
	flags.BoolVarP(cmd, &f.dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	return cmd, nil
}
