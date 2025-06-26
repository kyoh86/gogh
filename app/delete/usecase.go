package delete

import (
	"context"
	"fmt"
	"os"

	"github.com/kyoh86/gogh/v4/core/hosting"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// Usecase defines the use case for deleting repositories
type Usecase struct {
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
	hostingService   hosting.HostingService
	referenceParser  repository.ReferenceParser
}

// NewUsecase creates a new instance of Usecase
func NewUsecase(
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	hostingService hosting.HostingService,
	referenceParser repository.ReferenceParser,
) *Usecase {
	return &Usecase{
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
func (uc *Usecase) Execute(ctx context.Context, refs string, opts Options) error {
	ref, err := uc.referenceParser.Parse(refs)
	if err != nil {
		return err
	}
	if err := uc.deleteLocal(ctx, *ref, opts); err != nil {
		return fmt.Errorf("deleting local: %w", err)
	}
	if err := uc.deleteRemote(ctx, *ref, opts); err != nil {
		return fmt.Errorf("deleting remote: %w", err)
	}
	return nil
}

func (uc *Usecase) deleteLocal(ctx context.Context, ref repository.Reference, opts Options) error {
	if !opts.Local {
		return nil
	}
	match, err := uc.finderService.FindByReference(ctx, uc.workspaceService, ref)
	if err != nil {
		return fmt.Errorf("finding local repository: %w", err)
	}
	if match == nil {
		return nil
	}
	return os.RemoveAll(match.FullPath())
}

func (uc *Usecase) deleteRemote(ctx context.Context, ref repository.Reference, opts Options) error {
	if !opts.Remote {
		return nil
	}
	return uc.hostingService.DeleteRepository(ctx, ref)
}
