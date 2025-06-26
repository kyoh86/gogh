package list

import (
	"context"
	"io"

	"github.com/kyoh86/gogh/v4/app/overlay/describe"
	"github.com/kyoh86/gogh/v4/core/overlay"
)

type Usecase struct {
	overlayService overlay.OverlayService
	writer         io.Writer
}

func NewUsecase(
	overlayService overlay.OverlayService,
	writer io.Writer,
) *Usecase {
	return &Usecase{
		overlayService: overlayService,
		writer:         writer,
	}
}

func (uc *Usecase) Execute(ctx context.Context, asJSON, withSource bool) error {
	var usecase interface {
		Execute(ctx context.Context, s describe.Overlay) error
	}
	if asJSON {
		if withSource {
			usecase = describe.NewJSONWithContentUsecase(uc.overlayService, uc.writer)
		} else {
			usecase = describe.NewJSONUsecase(uc.writer)
		}
	} else {
		if withSource {
			usecase = describe.NewDetailUsecase(uc.overlayService, uc.writer)
		} else {
			usecase = describe.NewOnelineUsecase(uc.writer)
		}
	}
	for s, err := range uc.overlayService.List() {
		if err != nil {
			return err
		}
		if s == nil {
			continue
		}
		if err := usecase.Execute(ctx, s); err != nil {
			return err
		}
	}
	return nil
}
