package add

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/core/overlay"
)

// Usecase represents the create use case
type Usecase struct {
	overlayService overlay.OverlayService
}

func NewUsecase(overlayService overlay.OverlayService) *Usecase {
	return &Usecase{
		overlayService: overlayService,
	}
}

func (uc *Usecase) Execute(ctx context.Context, name, relativePath string, content io.Reader) (string, error) {
	e := overlay.Entry{
		Name:         name,
		RelativePath: relativePath,
		Content:      content,
	}
	return uc.overlayService.Add(ctx, e)
}
