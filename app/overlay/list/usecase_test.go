package list_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/overlay/list"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"go.uber.org/mock/gomock"
)

// testOverlay implements overlay.Overlay for testing
type testOverlay struct {
	id           uuid.UUID
	name         string
	relativePath string
}

func (t testOverlay) ID() string           { return t.id.String() }
func (t testOverlay) UUID() uuid.UUID      { return t.id }
func (t testOverlay) Name() string         { return t.name }
func (t testOverlay) RelativePath() string { return t.relativePath }

func TestUsecase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("Pass through overlays", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := overlay_mock.NewMockOverlayService(ctrl)
		overlay1 := testOverlay{
			id:           uuid.New(),
			name:         "overlay1",
			relativePath: "path/to/overlay1",
		}
		overlay2 := testOverlay{
			id:           uuid.New(),
			name:         "overlay2",
			relativePath: "path/to/overlay2",
		}

		mockService.EXPECT().
			List().
			Return(func(yield func(overlay.Overlay, error) bool) {
				yield(overlay1, nil)
				yield(overlay2, nil)
			})

		uc := testtarget.NewUsecase(mockService)
		got := make([]overlay.Overlay, 0, 2)
		for o, err := range uc.Execute(ctx) {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got = append(got, o)
		}

		if len(got) != 2 {
			t.Fatalf("got %d overlays, want 2", len(got))
		}
		if got[0].ID() != overlay1.ID() {
			t.Fatalf("overlay[0] id = %s, want %s", got[0].ID(), overlay1.ID())
		}
		if got[1].ID() != overlay2.ID() {
			t.Fatalf("overlay[1] id = %s, want %s", got[1].ID(), overlay2.ID())
		}
	})

	t.Run("Pass through errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := overlay_mock.NewMockOverlayService(ctrl)
		mockService.EXPECT().
			List().
			Return(func(yield func(overlay.Overlay, error) bool) {
				yield(nil, errors.New("list failed"))
			})

		uc := testtarget.NewUsecase(mockService)
		for _, err := range uc.Execute(ctx) {
			if err == nil {
				t.Fatal("expected error")
			}
			return
		}
		t.Fatal("expected iterator result")
	})
}
