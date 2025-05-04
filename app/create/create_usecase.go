package create

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/core/git"
	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
	"github.com/kyoh86/gogh/v3/infra/github"
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
}

func (uc *UseCase) Execute(ctx context.Context, ref repository.Reference, options *CreateOptions) error {
	if options != nil && options.Local {
		// ctrl := local.NewController(uc.workspaceService.GetDefaultRoot())
		// if err := ctrl.Delete(ctx, ref, nil); err != nil {
		// 	if !os.IsNotExist(err) {
		// 		return fmt.Errorf("delete local: %w", err)
		// 	}
		// }
	}

	if options != nil && options.Remote {
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
	}
	return nil
}
