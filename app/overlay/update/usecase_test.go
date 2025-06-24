package update_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/overlay/update"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name         string
		overlayID    string
		overlayName  string
		relativePath string
		content      string
		setupMock    func(*gomock.Controller) *overlay_mock.MockOverlayService
		wantErr      bool
	}{
		{
			name:         "Successfully update overlay with all fields",
			overlayID:    uuid.New().String(),
			overlayName:  "updated-overlay",
			relativePath: "path/to/updated/file.txt",
			content:      "Updated overlay content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry overlay.Entry) error {
						// Validate UUID format
						if _, err := uuid.Parse(id); err != nil {
							t.Errorf("Expected valid UUID, got %s", id)
						}
						// Validate name
						if entry.Name != "updated-overlay" {
							t.Errorf("Expected name 'updated-overlay', got %s", entry.Name)
						}
						// Validate relative path
						if entry.RelativePath != "path/to/updated/file.txt" {
							t.Errorf("Expected relative path 'path/to/updated/file.txt', got %s", entry.RelativePath)
						}
						// Verify content can be read
						buf := new(bytes.Buffer)
						if _, err := buf.ReadFrom(entry.Content); err != nil {
							t.Errorf("Failed to read content: %v", err)
						}
						if buf.String() != "Updated overlay content" {
							t.Errorf("Expected specific content, got %s", buf.String())
						}
						return nil
					},
				)
				return os
			},
			wantErr: false,
		},
		{
			name:         "Update overlay with empty name",
			overlayID:    uuid.New().String(),
			overlayName:  "",
			relativePath: "file.txt",
			content:      "content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry overlay.Entry) error {
						if entry.Name != "" {
							t.Errorf("Expected empty name, got %s", entry.Name)
						}
						return nil
					},
				)
				return os
			},
			wantErr: false,
		},
		{
			name:         "Update overlay with empty relative path",
			overlayID:    uuid.New().String(),
			overlayName:  "test-overlay",
			relativePath: "",
			content:      "content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry overlay.Entry) error {
						if entry.RelativePath != "" {
							t.Errorf("Expected empty relative path, got %s", entry.RelativePath)
						}
						return nil
					},
				)
				return os
			},
			wantErr: false,
		},
		{
			name:         "Update overlay with empty content",
			overlayID:    uuid.New().String(),
			overlayName:  "empty-content-overlay",
			relativePath: "empty.txt",
			content:      "",
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
			name:         "Update overlay with complex path",
			overlayID:    uuid.New().String(),
			overlayName:  "config-overlay",
			relativePath: ".config/app/settings/production.json",
			content:      `{"debug": false, "version": "1.0.0"}`,
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry overlay.Entry) error {
						if entry.RelativePath != ".config/app/settings/production.json" {
							t.Errorf("Unexpected relative path: %s", entry.RelativePath)
						}
						return nil
					},
				)
				return os
			},
			wantErr: false,
		},
		{
			name:         "Update overlay with large content",
			overlayID:    uuid.New().String(),
			overlayName:  "large-overlay",
			relativePath: "large-file.txt",
			content:      strings.Repeat("This is a repeated line for testing.\n", 1000),
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
			name:         "Update non-existent overlay",
			overlayID:    uuid.New().String(),
			overlayName:  "non-existent",
			relativePath: "file.txt",
			content:      "content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(errors.New("overlay not found"))
				return os
			},
			wantErr: true,
		},
		{
			name:         "Update with invalid overlay ID",
			overlayID:    "invalid-id",
			overlayName:  "test",
			relativePath: "file.txt",
			content:      "content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, "invalid-id", gomock.Any()).Return(errors.New("invalid overlay ID"))
				return os
			},
			wantErr: true,
		},
		{
			name:         "Update with empty overlay ID",
			overlayID:    "",
			overlayName:  "test",
			relativePath: "file.txt",
			content:      "content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, "", gomock.Any()).Return(errors.New("overlay ID is required"))
				return os
			},
			wantErr: true,
		},
		{
			name:         "Update with special characters in path",
			overlayID:    uuid.New().String(),
			overlayName:  "special-overlay",
			relativePath: "path/with spaces/and-dashes/file_name (copy).txt",
			content:      "content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry overlay.Entry) error {
						if entry.RelativePath != "path/with spaces/and-dashes/file_name (copy).txt" {
							t.Errorf("Unexpected relative path: %s", entry.RelativePath)
						}
						return nil
					},
				)
				return os
			},
			wantErr: false,
		},
		{
			name:         "Update with binary content",
			overlayID:    uuid.New().String(),
			overlayName:  "binary-overlay",
			relativePath: "binary.dat",
			content:      string([]byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}),
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
			name:         "Service returns unexpected error",
			overlayID:    uuid.New().String(),
			overlayName:  "test",
			relativePath: "file.txt",
			content:      "content",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(errors.New("storage error"))
				return os
			},
			wantErr: true,
		},
		{
			name:         "Update all fields empty except ID",
			overlayID:    uuid.New().String(),
			overlayName:  "",
			relativePath: "",
			content:      "",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry overlay.Entry) error {
						// Service might validate and return error
						return errors.New("at least one field must be provided")
					},
				)
				return os
			},
			wantErr: true,
		},
		{
			name:         "Update with JSON content",
			overlayID:    uuid.New().String(),
			overlayName:  "json-overlay",
			relativePath: "config.json",
			content: `{
  "name": "test-app",
  "version": "1.2.3",
  "dependencies": {
    "library-a": "^2.0.0",
    "library-b": "~3.1.0"
  },
  "scripts": {
    "start": "node index.js",
    "test": "jest"
  }
}`,
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(nil)
				return os
			},
			wantErr: false,
		},
		{
			name:         "Update with unicode content",
			overlayID:    uuid.New().String(),
			overlayName:  "unicode-overlay",
			relativePath: "unicode.txt",
			content:      "Hello 世界! 🌍 Unicode test: αβγδε",
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(nil)
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
			uc := testtarget.NewUseCase(os)

			err := uc.Execute(ctx, tc.overlayID, tc.overlayName, tc.relativePath, strings.NewReader(tc.content))
			if (err != nil) != tc.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestUseCase_Execute_ReaderBehavior(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test that the reader is passed through correctly
	customReader := strings.NewReader("test overlay content")

	os := overlay_mock.NewMockOverlayService(ctrl)
	os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, id string, entry overlay.Entry) error {
			// Verify that the same reader instance is passed
			if entry.Content != customReader {
				t.Error("Expected the same reader instance to be passed")
			}
			// Verify other fields are passed correctly
			if entry.Name != "test-overlay" {
				t.Errorf("Expected name 'test-overlay', got %s", entry.Name)
			}
			if entry.RelativePath != "test/path.txt" {
				t.Errorf("Expected relative path 'test/path.txt', got %s", entry.RelativePath)
			}
			return nil
		},
	)

	uc := testtarget.NewUseCase(os)
	err := uc.Execute(ctx, uuid.New().String(), "test-overlay", "test/path.txt", customReader)
	if err != nil {
		t.Errorf("Execute() unexpected error = %v", err)
	}
}

func TestUseCase_Execute_MultipleReaders(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os := overlay_mock.NewMockOverlayService(ctrl)
	uc := testtarget.NewUseCase(os)

	// Test with different reader types
	readers := []struct {
		name   string
		reader io.Reader
	}{
		{"strings.Reader", strings.NewReader("content from strings.Reader")},
		{"bytes.Buffer", bytes.NewBufferString("content from bytes.Buffer")},
		{"bytes.Reader", bytes.NewReader([]byte("content from bytes.Reader"))},
	}

	for _, r := range readers {
		os.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, id string, entry overlay.Entry) error {
				// Verify content can be read from any reader type
				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(entry.Content)
				if err != nil {
					t.Errorf("%s: Failed to read content: %v", r.name, err)
				}
				return nil
			},
		)

		err := uc.Execute(ctx, uuid.New().String(), r.name, "file.txt", r.reader)
		if err != nil {
			t.Errorf("%s: Execute() unexpected error = %v", r.name, err)
		}
	}
}
