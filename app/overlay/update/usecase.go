package update

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/core/overlay"
)

// Usecase is a struct that encapsulates the overlay service for updating overlays.
type Usecase struct {
	overlayService overlay.OverlayService
}

// NewUsecase creates a new instance of Usecase for updating overlays.
func NewUsecase(overlayService overlay.OverlayService) *Usecase {
	return &Usecase{overlayService: overlayService}
}

// Execute applies a new overlay identified by its ID.
func (uc *Usecase) Execute(ctx context.Context, overlayID, name, relativePath string, content io.Reader) error {
	return uc.overlayService.Update(ctx, overlayID, overlay.Entry{
		Name:         name,
		RelativePath: relativePath,
		Content:      content,
	})
}
