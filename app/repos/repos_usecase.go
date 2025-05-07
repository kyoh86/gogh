package repos

import (
	"context"
	"iter"

	"github.com/kyoh86/gogh/v3/core/hosting"
)

type UseCase struct {
	hostingService hosting.HostingService
}

func NewUseCase(hostingService hosting.HostingService) *UseCase {
	return &UseCase{
		hostingService: hostingService,
	}
}

type Options struct {
	hosting.ListRepositoryOptions
}

func (uc *UseCase) Execute(ctx context.Context, options Options) iter.Seq2[*hosting.Repository, error] {
	return uc.hostingService.ListRepository(ctx, options.ListRepositoryOptions)
}
