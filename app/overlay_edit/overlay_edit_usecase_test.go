package overlay_edit_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/overlay_edit"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"go.uber.org/mock/gomock"
)

// mockReadCloser implements io.ReadCloser
type mockReadCloser struct {
	io.Reader
	closed bool
}

func (m *mockReadCloser) Close() error {
	m.closed = true
	return nil
}

func TestUseCase_ExtractOverlay(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		overlayID string
		setupMock func(*gomock.Controller) *overlay_mock.MockOverlayService
		wantErr   bool
		wantData  string
	}{
		{
			name:      "Successfully extract overlay",
			overlayID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				content := "This is the overlay content"
				reader := &mockReadCloser{Reader: strings.NewReader(content)}
				os.EXPECT().Open(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string) (io.ReadCloser, error) {
						// Validate UUID format
						if _, err := uuid.Parse(id); err != nil {
							t.Errorf("Expected valid UUID, got %s", id)
						}
						return reader, nil
					},
				)
				return os
			},
			wantErr:  false,
			wantData: "This is the overlay content",
		},
		{
			name:      "Extract overlay with binary content",
			overlayID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				binaryData := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}
				reader := &mockReadCloser{Reader: bytes.NewReader(binaryData)}
				os.EXPECT().Open(ctx, gomock.Any()).Return(reader, nil)
				return os
			},
			wantErr:  false,
			wantData: string([]byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}),
		},
		{
			name:      "Extract large overlay",
			overlayID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				largeContent := strings.Repeat("This is a repeated line.\n", 1000)
				reader := &mockReadCloser{Reader: strings.NewReader(largeContent)}
				os.EXPECT().Open(ctx, gomock.Any()).Return(reader, nil)
				return os
			},
			wantErr:  false,
			wantData: strings.Repeat("This is a repeated line.\n", 1000),
		},
		{
			name:      "Extract empty overlay",
			overlayID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				reader := &mockReadCloser{Reader: strings.NewReader("")}
				os.EXPECT().Open(ctx, gomock.Any()).Return(reader, nil)
				return os
			},
			wantErr:  false,
			wantData: "",
		},
		{
			name:      "Overlay not found",
			overlayID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Open(ctx, gomock.Any()).Return(nil, errors.New("overlay not found"))
				return os
			},
			wantErr: true,
		},
		{
			name:      "Invalid overlay ID",
			overlayID: "invalid-id",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Open(ctx, "invalid-id").Return(nil, errors.New("invalid overlay ID"))
				return os
			},
			wantErr: true,
		},
		{
			name:      "Empty overlay ID",
			overlayID: "",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Open(ctx, "").Return(nil, errors.New("overlay ID is required"))
				return os
			},
			wantErr: true,
		},
		{
			name:      "Service returns unexpected error",
			overlayID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Open(ctx, gomock.Any()).Return(nil, errors.New("storage error"))
				return os
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			os := tc.setupMock(ctrl)
			uc := overlay_edit.NewUseCase(os)

			var buf bytes.Buffer
			err := uc.ExtractOverlay(ctx, tc.overlayID, &buf)
			if (err != nil) != tc.wantErr {
				t.Errorf("ExtractOverlay() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				if buf.String() != tc.wantData {
					t.Errorf("ExtractOverlay() data = %q, want %q", buf.String(), tc.wantData)
				}
			}
		})
	}
}

func TestUseCase_UpdateOverlay(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		overlayID string
		content   string
		setupMock func(*gomock.Controller) *overlay_mock.MockOverlayService
		wantErr   bool
	}{
		{
			name:      "Successfully update overlay",
			overlayID: uuid.New().String(),
			content:   "Updated overlay content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry overlay.Entry) error {
						// Validate UUID format
						if _, err := uuid.Parse(id); err != nil {
							t.Errorf("Expected valid UUID, got %s", id)
						}
						// Verify content can be read
						buf := new(bytes.Buffer)
						if _, err := buf.ReadFrom(entry.Content); err != nil {
							t.Errorf("Failed to read content: %v", err)
						}
						if buf.String() != "Updated overlay content" {
							t.Errorf("Expected content 'Updated overlay content', got %s", buf.String())
						}
						return nil
					},
				)
				return os
			},
			wantErr: false,
		},
		{
			name:      "Update with empty content",
			overlayID: uuid.New().String(),
			content:   "",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry overlay.Entry) error {
						buf := new(bytes.Buffer)
						if _, err := buf.ReadFrom(entry.Content); err != nil {
							t.Errorf("Failed to read content: %v", err)
						}
						if buf.String() != "" {
							t.Errorf("Expected empty content, got %s", buf.String())
						}
						return nil
					},
				)
				return os
			},
			wantErr: false,
		},
		{
			name:      "Update with binary content",
			overlayID: uuid.New().String(),
			content:   string([]byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}),
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry overlay.Entry) error {
						// Just verify we can read binary content
						buf := new(bytes.Buffer)
						_, err := buf.ReadFrom(entry.Content)
						if err != nil {
							t.Errorf("Failed to read binary content: %v", err)
						}
						return nil
					},
				)
				return os
			},
			wantErr: false,
		},
		{
			name:      "Update with large content",
			overlayID: uuid.New().String(),
			content:   strings.Repeat("This is a repeated line.\n", 1000),
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry overlay.Entry) error {
						// Just verify we can read the content
						buf := new(bytes.Buffer)
						n, err := buf.ReadFrom(entry.Content)
						if err != nil {
							t.Errorf("Failed to read content: %v", err)
						}
						if n == 0 {
							t.Error("Expected non-empty content")
						}
						return nil
					},
				)
				return os
			},
			wantErr: false,
		},
		{
			name:      "Update non-existent overlay",
			overlayID: uuid.New().String(),
			content:   "content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(errors.New("overlay not found"))
				return os
			},
			wantErr: true,
		},
		{
			name:      "Update with invalid overlay ID",
			overlayID: "invalid-id",
			content:   "content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, "invalid-id", gomock.Any()).Return(errors.New("invalid overlay ID"))
				return os
			},
			wantErr: true,
		},
		{
			name:      "Update with empty overlay ID",
			overlayID: "",
			content:   "content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, "", gomock.Any()).Return(errors.New("overlay ID is required"))
				return os
			},
			wantErr: true,
		},
		{
			name:      "Service returns unexpected error",
			overlayID: uuid.New().String(),
			content:   "content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(errors.New("storage error"))
				return os
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			os := tc.setupMock(ctrl)
			uc := overlay_edit.NewUseCase(os)

			err := uc.UpdateOverlay(ctx, tc.overlayID, strings.NewReader(tc.content))
			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateOverlay() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestUseCase_ExtractOverlay_ReaderBehavior(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test that the reader is properly closed
	os := overlay_mock.NewMockOverlayService(ctrl)
	reader := &mockReadCloser{Reader: strings.NewReader("test content")}

	os.EXPECT().Open(ctx, gomock.Any()).Return(reader, nil)

	uc := overlay_edit.NewUseCase(os)
	var buf bytes.Buffer
	err := uc.ExtractOverlay(ctx, uuid.New().String(), &buf)
	if err != nil {
		t.Errorf("ExtractOverlay() unexpected error = %v", err)
	}

	if !reader.closed {
		t.Error("Expected reader to be closed")
	}
}

func TestUseCase_UpdateOverlay_ReaderPassthrough(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test that the reader is passed through correctly
	customReader := strings.NewReader("test content")

	os := overlay_mock.NewMockOverlayService(ctrl)
	os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, id string, entry overlay.Entry) error {
			// Verify that the same reader instance is passed
			if entry.Content != customReader {
				t.Error("Expected the same reader instance to be passed")
			}
			return nil
		},
	)

	uc := overlay_edit.NewUseCase(os)
	err := uc.UpdateOverlay(ctx, uuid.New().String(), customReader)
	if err != nil {
		t.Errorf("UpdateOverlay() unexpected error = %v", err)
	}
}
