package overlay_add_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/overlay_add"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name         string
		overlayName  string
		relativePath string
		content      string
		setupMock    func(*gomock.Controller) *overlay_mock.MockOverlayService
		wantErr      bool
		validateID   func(string) error
	}{
		{
			name:         "Successfully add overlay",
			overlayName:  "test-overlay",
			relativePath: "path/to/file.txt",
			content:      "This is the overlay content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry overlay.Entry) (string, error) {
						// Validate the entry
						if entry.Name != "test-overlay" {
							t.Errorf("Expected name 'test-overlay', got %s", entry.Name)
						}
						if entry.RelativePath != "path/to/file.txt" {
							t.Errorf("Expected relative path 'path/to/file.txt', got %s", entry.RelativePath)
						}
						// Read content to verify
						buf := new(bytes.Buffer)
						if _, err := buf.ReadFrom(entry.Content); err != nil {
							t.Errorf("Failed to read content: %v", err)
						}
						if buf.String() != "This is the overlay content" {
							t.Errorf("Expected content 'This is the overlay content', got %s", buf.String())
						}
						return uuid.New().String(), nil
					},
				)
				return os
			},
			wantErr: false,
			validateID: func(id string) error {
				if id == "" {
					return errors.New("expected non-empty ID")
				}
				if _, err := uuid.Parse(id); err != nil {
					return errors.New("expected valid UUID")
				}
				return nil
			},
		},
		{
			name:         "Add overlay with empty name",
			overlayName:  "",
			relativePath: "file.txt",
			content:      "content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry overlay.Entry) (string, error) {
						if entry.Name != "" {
							t.Errorf("Expected empty name, got %s", entry.Name)
						}
						// Service might validate and return error
						return "", errors.New("overlay name is required")
					},
				)
				return os
			},
			wantErr: true,
		},
		{
			name:         "Add overlay with empty relative path",
			overlayName:  "test-overlay",
			relativePath: "",
			content:      "content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry overlay.Entry) (string, error) {
						if entry.RelativePath != "" {
							t.Errorf("Expected empty relative path, got %s", entry.RelativePath)
						}
						// Service might validate and return error
						return "", errors.New("relative path is required")
					},
				)
				return os
			},
			wantErr: true,
		},
		{
			name:         "Add overlay with complex path",
			overlayName:  "config-overlay",
			relativePath: ".config/app/settings.json",
			content:      `{"key": "value"}`,
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry overlay.Entry) (string, error) {
						if entry.RelativePath != ".config/app/settings.json" {
							t.Errorf("Expected relative path '.config/app/settings.json', got %s", entry.RelativePath)
						}
						return uuid.New().String(), nil
					},
				)
				return os
			},
			wantErr: false,
		},
		{
			name:         "Add overlay with large content",
			overlayName:  "large-overlay",
			relativePath: "large-file.txt",
			content:      strings.Repeat("This is a repeated line.\n", 1000),
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry overlay.Entry) (string, error) {
						// Just verify we can read the content
						buf := new(bytes.Buffer)
						n, err := buf.ReadFrom(entry.Content)
						if err != nil {
							t.Errorf("Failed to read content: %v", err)
						}
						if n == 0 {
							t.Error("Expected non-empty content")
						}
						return uuid.New().String(), nil
					},
				)
				return os
			},
			wantErr: false,
		},
		{
			name:         "Service returns error",
			overlayName:  "error-overlay",
			relativePath: "error.txt",
			content:      "content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Add(ctx, gomock.Any()).Return("", errors.New("storage full"))
				return os
			},
			wantErr: true,
		},
		{
			name:         "Add overlay with empty content",
			overlayName:  "empty-overlay",
			relativePath: "empty.txt",
			content:      "",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry overlay.Entry) (string, error) {
						buf := new(bytes.Buffer)
						if _, err := buf.ReadFrom(entry.Content); err != nil {
							t.Errorf("Failed to read content: %v", err)
						}
						if buf.String() != "" {
							t.Errorf("Expected empty content, got %s", buf.String())
						}
						return uuid.New().String(), nil
					},
				)
				return os
			},
			wantErr: false,
		},
		{
			name:         "Add overlay with special characters in path",
			overlayName:  "special-overlay",
			relativePath: "path/with spaces/and-dashes/file_name.txt",
			content:      "content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry overlay.Entry) (string, error) {
						if entry.RelativePath != "path/with spaces/and-dashes/file_name.txt" {
							t.Errorf("Unexpected relative path: %s", entry.RelativePath)
						}
						return uuid.New().String(), nil
					},
				)
				return os
			},
			wantErr: false,
		},
		{
			name:         "Add overlay with binary content",
			overlayName:  "binary-overlay",
			relativePath: "binary.dat",
			content:      string([]byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}),
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry overlay.Entry) (string, error) {
						// Just verify we can read binary content
						buf := new(bytes.Buffer)
						_, err := buf.ReadFrom(entry.Content)
						if err != nil {
							t.Errorf("Failed to read binary content: %v", err)
						}
						return uuid.New().String(), nil
					},
				)
				return os
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			os := tc.setupMock(ctrl)
			uc := overlay_add.NewUseCase(os)

			id, err := uc.Execute(ctx, tc.overlayName, tc.relativePath, strings.NewReader(tc.content))
			if (err != nil) != tc.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && tc.validateID != nil {
				if err := tc.validateID(id); err != nil {
					t.Errorf("ID validation failed: %v", err)
				}
			}
		})
	}
}

func TestUseCase_Execute_ReaderBehavior(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test that the reader is passed through correctly
	customReader := &countingReader{Reader: strings.NewReader("test content")}

	os := overlay_mock.NewMockOverlayService(ctrl)
	os.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, entry overlay.Entry) (string, error) {
			// Verify that the same reader instance is passed
			if entry.Content != customReader {
				t.Error("Expected the same reader instance to be passed")
			}
			return uuid.New().String(), nil
		},
	)

	uc := overlay_add.NewUseCase(os)
	_, err := uc.Execute(ctx, "test", "test.txt", customReader)
	if err != nil {
		t.Errorf("Execute() unexpected error = %v", err)
	}
}

// countingReader is a helper to verify reader behavior
type countingReader struct {
	*strings.Reader
	readCount int
}

func (cr *countingReader) Read(p []byte) (n int, err error) {
	cr.readCount++
	return cr.Reader.Read(p)
}
