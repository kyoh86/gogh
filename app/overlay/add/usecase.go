package add

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/core/overlay"
)

// UseCase represents the create use case
type UseCase struct {
	overlayService overlay.OverlayService
}

func NewUseCase(overlayService overlay.OverlayService) *UseCase {
	return &UseCase{
		overlayService: overlayService,
	}
}

func (uc *UseCase) Execute(ctx context.Context, name, relativePath string, content io.Reader) (string, error) {
	e := overlay.Entry{
		Name:         name,
		RelativePath: relativePath,
		Content:      content,
	}
	return uc.overlayService.Add(ctx, e)
}
