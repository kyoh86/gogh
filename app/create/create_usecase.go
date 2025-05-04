package create

import (
	"context"

	"github.com/kyoh86/gogh/v3/app/clone"
	"github.com/kyoh86/gogh/v3/core/git"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
)

// UseCase represents the create use case
type UseCase struct {
	hostingService   hosting.HostingService
	gitService       git.GitService
	workspaceService workspace.WorkspaceService
}

func NewUseCase(
	hostingService hosting.HostingService,
	gitService git.GitService,
	workspaceService workspace.WorkspaceService,
) *UseCase {
	return &UseCase{
		hostingService:   hostingService,
		gitService:       gitService,
		workspaceService: workspaceService,
	}
}

type CreateOptions struct {
	Local  bool
	Remote bool
	Alias  *repository.Reference
	hosting.CreateRepositoryOptions
}

func (uc *UseCase) Execute(ctx context.Context, ref repository.Reference, options CreateOptions) error {
	if options.Remote {
		uc.hostingService.CreateRepository(ctx, ref, options.CreateRepositoryOptions)
		// adaptor, _, err := RemoteControllerFor(ctx, *tokenService, ref)
		// if err != nil {
		// 	return fmt.Errorf("failed to get token for %s/%s: %w", ref.Host(), ref.Owner(), err)
		// }
		// if err := remote.NewController(adaptor).Delete(ctx, ref.Owner(), ref.Name(), nil); err != nil {
		// 	var gherr *github.ErrorResponse
		// 	if errors.As(err, &gherr) && gherr.Response.StatusCode == http.StatusForbidden {
		// 		log.FromContext(ctx).Errorf("Failed to delete a remote repository: there is no permission to delete %q", ref.URL())
		// 		log.FromContext(ctx).Error(`Add scope "delete_repo" for the token`)
		// 	} else {
		// 		return err
		// 	}
		// }
		if options.Local {
			cloneUseCase := clone.NewUseCase(uc.hostingService, uc.workspaceService)
			return cloneUseCase.Execute(ctx, ref, &clone.CloneOptions{
				Alias: options.Alias,
			})
		}
	} else if options.Local {
		layout := uc.workspaceService.GetDefaultLayout()
		path, err := layout.CreateRepositoryFolder(ref)
		if err != nil {
			return err
		}
		remoteURL, err := uc.hostingService.GetURLOf(ref)
		if err != nil {
			return err
		}
		if err := uc.gitService.Init(remoteURL.String(), path, false); err != nil {
			return err
		}
	}

	return nil
}
