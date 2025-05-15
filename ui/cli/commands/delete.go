package commands

import (
	"context"
	"fmt"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v3/app/delete"
	"github.com/kyoh86/gogh/v3/app/repos"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/spf13/cobra"
)

func NewDeleteCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		local  bool
		remote bool
		force  bool
		dryrun bool
	}

	reposUseCase := repos.NewUseCase(svc.HostingService)

	checkFlags := func(ctx context.Context, args []string) (string, error) {
		if len(args) != 0 {
			return args[0], nil
		}
		var opts []huh.Option[string]
		for repo, err := range reposUseCase.Execute(ctx, repos.Options{}) {
			if err != nil {
				return "", fmt.Errorf("failed to list repositories: %w", err)
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
		Use:     "delete [flags] [[OWNER/]NAME]",
		Aliases: []string{"remove"},
		Short:   "Delete local and remote repository",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, refs []string) error {
			ctx := cmd.Context()
			selected, err := checkFlags(ctx, refs)
			if err != nil {
				return err
			}

			if err := prepareFlags(ctx, selected); err != nil {
				return err
			}

			if f.local {
				fmt.Printf("deleting local %s\n", selected)
			}
			if f.remote {
				fmt.Printf("deleting remote %s\n", selected)
			}
			if f.dryrun {
				return nil
			}
			useCase := delete.NewUseCase(
				svc.WorkspaceService,
				svc.FinderService,
				svc.HostingService,
				svc.ReferenceParser,
			)
			if err := useCase.Execute(ctx, selected, delete.Options{
				Local:  f.local,
				Remote: f.remote,
			}); err != nil {
				return fmt.Errorf("failed to delete the repository: %w", err)
			}
			if f.local {
				log.FromContext(ctx).Infof("deleted local %s", selected)
			}
			if f.remote {
				log.FromContext(ctx).Infof("deleted remote %s", selected)
			}
			return nil
		},
	}
	//TODO: --no-local?
	cmd.Flags().BoolVarP(&f.local, "local", "", true, "Delete local repository.")
	cmd.Flags().
		BoolVarP(&f.remote, "remote", "", false, "Delete remote repository.")
	cmd.Flags().
		BoolVarP(&f.force, "force", "", false, "Do NOT confirm to delete.")
	cmd.Flags().
		BoolVarP(&f.dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	return cmd, nil
}
