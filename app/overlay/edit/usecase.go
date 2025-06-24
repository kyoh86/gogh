package edit

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/core/overlay"
)

type UseCase struct {
	overlayService overlay.OverlayService
}

func NewUseCase(overlayService overlay.OverlayService) *UseCase {
	return &UseCase{overlayService: overlayService}
}

// ExtractOverlay extracts the overlay by its ID and writes it to the provided writer.
func (uc *UseCase) ExtractOverlay(ctx context.Context, overlayID string, w io.Writer) error {
	r, err := uc.overlayService.Open(ctx, overlayID)
	if err != nil {
		return err
	}
	defer r.Close()
	_, err = io.Copy(w, r)
	return err
}

// UpdateOverlay applies a new overlay identified by its ID.
func (uc *UseCase) UpdateOverlay(ctx context.Context, overlayID string, r io.Reader) error {
	return uc.overlayService.Update(ctx, overlayID, overlay.Entry{Content: r})
}
