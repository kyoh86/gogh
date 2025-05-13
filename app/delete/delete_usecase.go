package delete

import (
	"context"
	"os"

	"github.com/kyoh86/gogh/v3/core/hosting"
	"github.com/kyoh86/gogh/v3/core/repository"
	"github.com/kyoh86/gogh/v3/core/workspace"
)

// UseCase defines the use case for deleting repositories
type UseCase struct {
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
	hostingService   hosting.HostingService
	referenceParser  repository.ReferenceParser
}

// NewUseCase creates a new instance of UseCase
func NewUseCase(
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	hostingService hosting.HostingService,
	referenceParser repository.ReferenceParser,
) *UseCase {
	return &UseCase{
		workspaceService: workspaceService,
		finderService:    finderService,
		hostingService:   hostingService,
		referenceParser:  referenceParser,
	}
}

// Options defines the options for deleting repositories
type Options struct {
	Remote bool
	Local  bool
}

// Execute deletes the specified repository from local and remote
func (u *UseCase) Execute(ctx context.Context, refs string, opts Options) error {
	ref, err := u.referenceParser.Parse(refs)
	if err != nil {
		return err
	}
	if err := u.deleteLocal(ctx, *ref, opts); err != nil {
		return err
	}
	if err := u.deleteRemote(ctx, *ref, opts); err != nil {
		return err
	}
	return nil
}

func (u *UseCase) deleteLocal(ctx context.Context, ref repository.Reference, opts Options) error {
	if !opts.Local {
		return nil
	}
	match, err := u.finderService.FindByReference(ctx, u.workspaceService, ref)
	if err != nil {
		return err
	}
	if match == nil {
		return nil
	}
	return os.RemoveAll(match.FullPath())
}

func (u *UseCase) deleteRemote(ctx context.Context, ref repository.Reference, opts Options) error {
	if !opts.Remote {
		return nil
	}
	return u.hostingService.DeleteRepository(ctx, ref)
}
