package list_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"iter"
	"strings"
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

func TestUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("Success: List overlays as one-line", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := overlay_mock.NewMockOverlayService(ctrl)
		buf := &bytes.Buffer{}

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

		// Create an iterator that yields overlays
		mockService.EXPECT().
			List().
			Return(func() iter.Seq2[overlay.Overlay, error] {
				return func(yield func(overlay.Overlay, error) bool) {
					yield(overlay1, nil)
					yield(overlay2, nil)
				}
			}())

		uc := testtarget.NewUseCase(mockService, buf)
		err := uc.Execute(ctx, false, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, overlay1.id.String()[:8]) {
			t.Errorf("expected output to contain overlay1 ID prefix, got: %s", output)
		}
		if !strings.Contains(output, overlay2.id.String()[:8]) {
			t.Errorf("expected output to contain overlay2 ID prefix, got: %s", output)
		}
		if !strings.Contains(output, overlay1.name) {
			t.Errorf("expected output to contain overlay1 name, got: %s", output)
		}
		if !strings.Contains(output, overlay2.name) {
			t.Errorf("expected output to contain overlay2 name, got: %s", output)
		}
	})

	t.Run("Success: List overlays as JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := overlay_mock.NewMockOverlayService(ctrl)
		buf := &bytes.Buffer{}

		overlay1 := testOverlay{
			id:           uuid.New(),
			name:         "overlay1",
			relativePath: "path/to/overlay1",
		}

		mockService.EXPECT().
			List().
			Return(func() iter.Seq2[overlay.Overlay, error] {
				return func(yield func(overlay.Overlay, error) bool) {
					yield(overlay1, nil)
				}
			}())

		uc := testtarget.NewUseCase(mockService, buf)
		err := uc.Execute(ctx, true, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, `"id"`) {
			t.Errorf("expected JSON output to contain 'id' field")
		}
		if !strings.Contains(output, overlay1.id.String()) {
			t.Errorf("expected JSON output to contain overlay ID")
		}
	})

	t.Run("Success: List overlays with source as detail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := overlay_mock.NewMockOverlayService(ctrl)
		buf := &bytes.Buffer{}

		overlay1 := testOverlay{
			id:           uuid.New(),
			name:         "overlay1",
			relativePath: "path/to/overlay1",
		}

		mockService.EXPECT().
			List().
			Return(func() iter.Seq2[overlay.Overlay, error] {
				return func(yield func(overlay.Overlay, error) bool) {
					yield(overlay1, nil)
				}
			}())

		// For detail view with source, it will call Open to get content
		mockService.EXPECT().
			Open(ctx, overlay1.id.String()).
			Return(io.NopCloser(strings.NewReader("overlay content")), nil)

		uc := testtarget.NewUseCase(mockService, buf)
		err := uc.Execute(ctx, false, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, overlay1.id.String()) {
			t.Errorf("expected output to contain overlay ID")
		}
	})

	t.Run("Success: Skip nil overlays", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := overlay_mock.NewMockOverlayService(ctrl)
		buf := &bytes.Buffer{}

		overlay1 := testOverlay{
			id:           uuid.New(),
			name:         "overlay1",
			relativePath: "path/to/overlay1",
		}

		mockService.EXPECT().
			List().
			Return(func() iter.Seq2[overlay.Overlay, error] {
				return func(yield func(overlay.Overlay, error) bool) {
					yield(nil, nil) // nil overlay should be skipped
					yield(overlay1, nil)
				}
			}())

		uc := testtarget.NewUseCase(mockService, buf)
		err := uc.Execute(ctx, false, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, overlay1.id.String()[:8]) {
			t.Errorf("expected output to contain overlay1 ID prefix, got: %s", output)
		}
	})

	t.Run("Error: List returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := overlay_mock.NewMockOverlayService(ctrl)
		buf := &bytes.Buffer{}
		expectedErr := errors.New("list failed")

		mockService.EXPECT().
			List().
			Return(func() iter.Seq2[overlay.Overlay, error] {
				return func(yield func(overlay.Overlay, error) bool) {
					yield(nil, expectedErr)
				}
			}())

		uc := testtarget.NewUseCase(mockService, buf)
		err := uc.Execute(ctx, false, false)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("Error: Execute fails for an overlay", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := overlay_mock.NewMockOverlayService(ctrl)
		// Use a writer that fails on write
		failWriter := &failingWriter{err: errors.New("write failed")}

		overlay1 := testOverlay{
			id:           uuid.New(),
			name:         "overlay1",
			relativePath: "path/to/overlay1",
		}

		mockService.EXPECT().
			List().
			Return(func() iter.Seq2[overlay.Overlay, error] {
				return func(yield func(overlay.Overlay, error) bool) {
					yield(overlay1, nil)
				}
			}())

		uc := testtarget.NewUseCase(mockService, failWriter)
		err := uc.Execute(ctx, false, false)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "write failed") {
			t.Errorf("expected write failed error, got %v", err)
		}
	})
}

// failingWriter is a writer that always fails
type failingWriter struct {
	err error
}

func (f *failingWriter) Write(p []byte) (n int, err error) {
	return 0, f.err
}
