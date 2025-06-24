package overlay_show_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/overlay_show"
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

func TestUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name       string
		overlayID  string
		asJSON     bool
		withSource bool
		setupMock  func(*gomock.Controller) *overlay_mock.MockOverlayService
		wantErr    bool
		validate   func(*testing.T, string)
	}{
		{
			name:       "Show overlay as one-line",
			overlayID:  uuid.New().String(),
			asJSON:     false,
			withSource: false,
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				o := overlay.ConcreteOverlay(
					uuid.New(),
					"test-overlay",
					"path/to/file.txt",
				)
				os.EXPECT().Get(ctx, gomock.Any()).Return(o, nil)
				return os
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				// One-line format should contain ID prefix, name, and path
				if !strings.Contains(output, "test-overlay") {
					t.Error("Expected output to contain overlay name")
				}
				if !strings.Contains(output, "path/to/file.txt") {
					t.Error("Expected output to contain relative path")
				}
			},
		},
		{
			name:       "Show overlay as JSON",
			overlayID:  uuid.New().String(),
			asJSON:     true,
			withSource: false,
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				o := overlay.ConcreteOverlay(
					uuid.New(),
					"json-overlay",
					".config/settings.json",
				)
				os.EXPECT().Get(ctx, gomock.Any()).Return(o, nil)
				return os
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				// Should be valid JSON
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(output), &data); err != nil {
					t.Errorf("Expected valid JSON output, got error: %v", err)
				}
				// Check fields
				if data["name"] != "json-overlay" {
					t.Errorf("Expected name 'json-overlay', got %v", data["name"])
				}
				if data["relative_path"] != ".config/settings.json" {
					t.Errorf("Expected relative_path '.config/settings.json', got %v", data["relative_path"])
				}
			},
		},
		{
			name:       "Show overlay with source content (detail)",
			overlayID:  uuid.New().String(),
			asJSON:     false,
			withSource: true,
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				overlayID := uuid.New()
				o := overlay.ConcreteOverlay(
					overlayID,
					"detail-overlay",
					"README.md",
				)
				os.EXPECT().Get(ctx, gomock.Any()).Return(o, nil)

				// Expect Open to be called for content
				content := "# This is the overlay content\nLine 2\nLine 3"
				reader := &mockReadCloser{Reader: strings.NewReader(content)}
				os.EXPECT().Open(ctx, overlayID.String()).Return(reader, nil)

				return os
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				// Detail format should contain overlay info and content
				if !strings.Contains(output, "detail-overlay") {
					t.Error("Expected output to contain overlay name")
				}
				if !strings.Contains(output, "README.md") {
					t.Error("Expected output to contain relative path")
				}
				if !strings.Contains(output, "# This is the overlay content") {
					t.Error("Expected output to contain overlay content")
				}
			},
		},
		{
			name:       "Show overlay with source content as JSON",
			overlayID:  uuid.New().String(),
			asJSON:     true,
			withSource: true,
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				overlayID := uuid.New()
				o := overlay.ConcreteOverlay(
					overlayID,
					"json-with-content",
					"config.yaml",
				)
				os.EXPECT().Get(ctx, gomock.Any()).Return(o, nil)

				// Expect Open to be called for content
				content := "key: value\nlist:\n  - item1\n  - item2"
				reader := &mockReadCloser{Reader: strings.NewReader(content)}
				os.EXPECT().Open(ctx, overlayID.String()).Return(reader, nil)

				return os
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(output), &data); err != nil {
					t.Errorf("Expected valid JSON output, got error: %v", err)
				}
				// Check content field exists
				if content, ok := data["content"].(string); ok {
					if !strings.Contains(content, "key: value") {
						t.Errorf("Expected content to contain 'key: value', got %v", content)
					}
				} else {
					t.Error("Expected 'content' field in JSON output")
				}
			},
		},
		{
			name:       "Overlay not found",
			overlayID:  uuid.New().String(),
			asJSON:     false,
			withSource: false,
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Get(ctx, gomock.Any()).Return(nil, errors.New("overlay not found"))
				return os
			},
			wantErr: true,
		},
		{
			name:       "Invalid overlay ID",
			overlayID:  "invalid-id",
			asJSON:     false,
			withSource: false,
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Get(ctx, "invalid-id").Return(nil, errors.New("invalid overlay ID"))
				return os
			},
			wantErr: true,
		},
		{
			name:       "Empty overlay ID",
			overlayID:  "",
			asJSON:     false,
			withSource: false,
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				os.EXPECT().Get(ctx, "").Return(nil, errors.New("overlay ID is required"))
				return os
			},
			wantErr: true,
		},
		{
			name:       "Error reading overlay content",
			overlayID:  uuid.New().String(),
			asJSON:     false,
			withSource: true,
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				overlayID := uuid.New()
				o := overlay.ConcreteOverlay(
					overlayID,
					"error-overlay",
					"error.txt",
				)
				os.EXPECT().Get(ctx, gomock.Any()).Return(o, nil)

				// Open returns error
				os.EXPECT().Open(ctx, overlayID.String()).Return(nil, errors.New("cannot read content"))

				return os
			},
			wantErr: true,
		},
		{
			name:       "Show overlay with empty name",
			overlayID:  uuid.New().String(),
			asJSON:     true,
			withSource: false,
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				o := overlay.ConcreteOverlay(
					uuid.New(),
					"", // Empty name
					"file.txt",
				)
				os.EXPECT().Get(ctx, gomock.Any()).Return(o, nil)
				return os
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(output), &data); err != nil {
					t.Errorf("Expected valid JSON output, got error: %v", err)
				}
				if data["name"] != "" {
					t.Errorf("Expected empty name, got %v", data["name"])
				}
			},
		},
		{
			name:       "Show overlay with special characters in path",
			overlayID:  uuid.New().String(),
			asJSON:     false,
			withSource: false,
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				o := overlay.ConcreteOverlay(
					uuid.New(),
					"special-overlay",
					"path/with spaces/and-dashes/file_name (copy).txt",
				)
				os.EXPECT().Get(ctx, gomock.Any()).Return(o, nil)
				return os
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "path/with spaces/and-dashes/file_name (copy).txt") {
					t.Error("Expected output to contain full path with special characters")
				}
			},
		},
		{
			name:       "Show overlay with binary content",
			overlayID:  uuid.New().String(),
			asJSON:     true,
			withSource: true,
			setupMock: func(ctrl *gomock.Controller) *overlay_mock.MockOverlayService {
				os := overlay_mock.NewMockOverlayService(ctrl)
				overlayID := uuid.New()
				o := overlay.ConcreteOverlay(
					overlayID,
					"binary-overlay",
					"binary.dat",
				)
				os.EXPECT().Get(ctx, gomock.Any()).Return(o, nil)

				// Binary content
				binaryContent := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}
				reader := &mockReadCloser{Reader: bytes.NewReader(binaryContent)}
				os.EXPECT().Open(ctx, overlayID.String()).Return(reader, nil)

				return os
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				var data map[string]any
				if err := json.Unmarshal([]byte(output), &data); err != nil {
					t.Errorf("Expected valid JSON output, got error: %v", err)
				}
				// Content should be present (even if it's binary)
				if _, ok := data["content"]; !ok {
					t.Error("Expected 'content' field in JSON output")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var buf bytes.Buffer
			os := tc.setupMock(ctrl)
			uc := overlay_show.NewUseCase(os, &buf)

			err := uc.Execute(ctx, tc.overlayID, tc.asJSON, tc.withSource)
			if (err != nil) != tc.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && tc.validate != nil {
				tc.validate(t, buf.String())
			}
		})
	}
}

func TestUseCase_Execute_ServiceErrors(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var buf bytes.Buffer
	os := overlay_mock.NewMockOverlayService(ctrl)
	uc := overlay_show.NewUseCase(os, &buf)

	// Test service returning unexpected error
	os.EXPECT().Get(ctx, "test-id").Return(nil, errors.New("database connection error"))

	err := uc.Execute(ctx, "test-id", false, false)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "get overlay by ID") {
		t.Errorf("Expected error to contain 'get overlay by ID', got %v", err.Error())
	}
}

func TestUseCase_Execute_AllDisplayModes(t *testing.T) {
	ctx := context.Background()

	modes := []struct {
		name       string
		asJSON     bool
		withSource bool
	}{
		{"OneLine", false, false},
		{"Detail", false, true},
		{"JSON", true, false},
		{"JSONWithContent", true, true},
	}

	for _, mode := range modes {
		t.Run(mode.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var buf bytes.Buffer
			os := overlay_mock.NewMockOverlayService(ctrl)
			uc := overlay_show.NewUseCase(os, &buf)

			overlayID := uuid.New()
			o := overlay.ConcreteOverlay(
				overlayID,
				"test-overlay",
				"test.txt",
			)
			os.EXPECT().Get(ctx, gomock.Any()).Return(o, nil)

			if mode.withSource {
				// Expect content to be read
				reader := &mockReadCloser{Reader: strings.NewReader("test content")}
				os.EXPECT().Open(ctx, overlayID.String()).Return(reader, nil)
			}

			err := uc.Execute(ctx, overlayID.String(), mode.asJSON, mode.withSource)
			if err != nil {
				t.Errorf("Execute() unexpected error for mode %s: %v", mode.name, err)
			}

			// Verify some output was produced
			if buf.Len() == 0 {
				t.Errorf("Expected output for mode %s, got empty", mode.name)
			}
		})
	}
}
