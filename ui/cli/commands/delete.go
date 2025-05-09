package commands

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v3/app/delete"
	"github.com/kyoh86/gogh/v3/app/repos"
	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/spf13/cobra"
)

func NewDeleteCommand(
	conf *config.ConfigStore,
	defaultNameService repository.DefaultNameService,
	tokenService auth.TokenService,
	hostingService hosting.HostingService,
	finderService workspace.FinderService,
	workspaceService workspace.WorkspaceService,
) *cobra.Command {
	var f struct {
		local  bool
		remote bool
		force  bool
		dryrun bool
	}

	reposUseCase := repos.NewUseCase(hostingService)

	checkFlags := func(ctx context.Context, args []string) (string, error) {
		if len(args) != 0 {
			return args[0], nil
		}
		var options []huh.Option[string]
		for repo, err := range reposUseCase.Execute(ctx, repos.Options{}) {
			if err != nil {
				return "", err
			}
			options = append(options, huh.Option[string]{
				Key:   repo.Ref.String(),
				Value: repo.Ref.String(),
			})
		}
		var selected string
		if err := huh.NewForm(huh.NewGroup(
			huh.NewSelect[string]().
				Title("A repository to delete").
				Options(options...).
				Value(&selected),
		)).Run(); err != nil {
			return "", err
		}
		return selected, nil
	}

	prepareFlags := func(ctx context.Context, arg string) (*repository.Reference, error) {
		parser := repository.NewReferenceParser(defaultNameService.GetDefaultHostAndOwner())
		ref, err := parser.Parse(arg)
		if err != nil {
			return nil, err
		}

		if !f.force {
			if f.local {
				var confirmed bool
				if err := huh.NewForm(huh.NewGroup(
					huh.NewConfirm().
						Title(fmt.Sprintf("Are you sure you want to delete local repository %s?", ref.String())).
						Value(&confirmed),
				)).Run(); err != nil {
					return nil, err
				}
				if !confirmed {
					return nil, context.Canceled
				}
			}
			if f.remote {
				var confirmed bool
				if err := huh.NewForm(huh.NewGroup(
					huh.NewConfirm().
						Title(fmt.Sprintf("Are you sure you want to delete remote-repository %s?", ref.String())).
						Value(&confirmed),
				)).Run(); err != nil {
					return nil, err
				}
				if !confirmed {
					return nil, context.Canceled
				}
			}
		}
		return ref, nil
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
			ref, err := prepareFlags(ctx, selected)
			if err != nil {
				return err
			}

			if f.dryrun {
				if f.local {
					fmt.Printf("deleting local %s\n", ref.String())
				}
				if f.remote {
					fmt.Printf("deleting remote %s\n", ref.String())
				}
				return nil
			}
			useCase := delete.NewUseCase(
				workspaceService,
				finderService,
				hostingService,
			)
			return useCase.Execute(ctx, *ref, delete.Options{
				Local:  f.local,
				Remote: f.remote,
			})
		},
	}
	cmd.Flags().BoolVarP(&f.local, "local", "", true, "Delete local repository.")
	cmd.Flags().
		BoolVarP(&f.remote, "remote", "", false, "Delete remote repository.")
	cmd.Flags().
		BoolVarP(&f.force, "force", "", false, "Do NOT confirm to delete.")
	cmd.Flags().
		BoolVarP(&f.dryrun, "dryrun", "", false, "Displays the operations that would be performed using the specified command without actually running them")
	return cmd
}
