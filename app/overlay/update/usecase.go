package update

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/core/overlay"
)

// UseCase is a struct that encapsulates the overlay service for updating overlays.
type UseCase struct {
	overlayService overlay.OverlayService
}

// NewUseCase creates a new instance of UseCase for updating overlays.
func NewUseCase(overlayService overlay.OverlayService) *UseCase {
	return &UseCase{overlayService: overlayService}
}

// Execute applies a new overlay identified by its ID.
func (uc *UseCase) Execute(ctx context.Context, overlayID, name, relativePath string, content io.Reader) error {
	return uc.overlayService.Update(ctx, overlayID, overlay.Entry{
		Name:         name,
		RelativePath: relativePath,
		Content:      content,
	})
}
