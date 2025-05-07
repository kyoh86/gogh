package commands

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v3/app/repos"
	"github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/domain/local"
	"github.com/kyoh86/gogh/v3/domain/remote"
	"github.com/kyoh86/gogh/v3/domain/reporef"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/infra/github"
	"github.com/spf13/cobra"
)

func NewDeleteCommand(
	conf *config.ConfigStore,
	defaultNameService repository.DefaultNameService,
	tokenService auth.TokenService,
	hostingService hosting.HostingService,
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

			parser := reporef.NewRepoRefParser(defaultNameService.GetDefaultHostAndOwner())
			ref, err := parser.Parse(selected)
			if err != nil {
				return err
			}

			if !f.force {
				if f.local {
					var confirmed bool
					if err := huh.NewForm(huh.NewGroup(
						huh.NewConfirm().
							Title(fmt.Sprintf("Are you sure you want to delete local repository %s?", ref.String())).
							Value(&confirmed),
					)).Run(); err != nil {
						return err
					}
					if !confirmed {
						return nil
					}
				}
				if f.remote {
					var confirmed bool
					if err := huh.NewForm(huh.NewGroup(
						huh.NewConfirm().
							Title(fmt.Sprintf("Are you sure you want to delete remote-repository %s?", ref.String())).
							Value(&confirmed),
					)).Run(); err != nil {
						return err
					}
					if !confirmed {
						return nil
					}
				}
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
			shouldReturn, err := newFunction(ctx, f.local, f.remote, f.force, f.dryrun, conf, ref, tokenService)
			if shouldReturn {
				return err
			}
			return nil
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

func newFunction(
	ctx context.Context,
	_local bool, _remote bool, force bool, dryrun bool,
	conf *config.ConfigStore,
	ref reporef.RepoRef,
	tokenService auth.TokenService,
) (bool, error) {
	if _local {
		ctrl := local.NewController(conf.DefaultRoot())
		if !force {
			var confirmed bool
			if err := huh.NewForm(huh.NewGroup(
				huh.NewConfirm().
					Title(fmt.Sprintf("Are you sure you want to delete local repository %s?", ref.String())).
					Value(&confirmed),
			)).Run(); err != nil {
				return true, err
			}
			if !confirmed {
				return true, nil
			}
		}
		if dryrun {
			fmt.Printf("deleting local %s\n", ref.String())
		} else if err := ctrl.Delete(ctx, ref, nil); err != nil {
			if !os.IsNotExist(err) {
				return true, fmt.Errorf("delete local: %w", err)
			}
		}
	}

	if _remote {
		if !force {
			var confirmed bool
			if err := huh.NewForm(huh.NewGroup(
				huh.NewConfirm().
					Title(fmt.Sprintf("Are you sure you want to delete remote-repository %s?", ref.String())).
					Value(&confirmed),
			)).Run(); err != nil {
				return true, err
			}
			if !confirmed {
				return true, nil
			}
		}
		adaptor, _, err := RemoteControllerFor(ctx, *tokenService, ref)
		if err != nil {
			return true, fmt.Errorf("failed to get token for %s/%s: %w", ref.Host(), ref.Owner(), err)
		}
		if dryrun {
			fmt.Printf("deleting remote %s\n", ref.String())
		} else if err := remote.NewController(adaptor).Delete(ctx, ref.Owner(), ref.Name(), nil); err != nil {
			var gherr *github.ErrorResponse
			if errors.As(err, &gherr) && gherr.Response.StatusCode == http.StatusForbidden {
				log.FromContext(ctx).Errorf("Failed to delete a remote repository: there is no permission to delete %q", ref.URL())
				log.FromContext(ctx).Error(`Add scope "delete_repo" for the token`)
			} else {
				return true, err
			}
		}
	}
	return false, nil
}
