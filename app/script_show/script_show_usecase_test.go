package script_show_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/script_show"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/script_mock"
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
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour)

	testCases := []struct {
		name       string
		scriptID   string
		asJSON     bool
		withSource bool
		setupMock  func(*gomock.Controller) *script_mock.MockScriptService
		wantErr    bool
		validate   func(*testing.T, string)
	}{
		{
			name:       "Show script as one-line",
			scriptID:   uuid.New().String(),
			asJSON:     false,
			withSource: false,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				s := script.ConcreteScript(
					uuid.New(),
					"test-script",
					createdAt,
					updatedAt,
				)
				ss.EXPECT().Get(ctx, gomock.Any()).Return(s, nil)
				return ss
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				// One-line format should contain ID prefix and name
				if !strings.Contains(output, "test-script") {
					t.Error("Expected output to contain script name")
				}
			},
		},
		{
			name:       "Show script as JSON",
			scriptID:   uuid.New().String(),
			asJSON:     true,
			withSource: false,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				s := script.ConcreteScript(
					uuid.New(),
					"json-script",
					createdAt,
					updatedAt,
				)
				ss.EXPECT().Get(ctx, gomock.Any()).Return(s, nil)
				return ss
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				// Should be valid JSON
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(output), &data); err != nil {
					t.Errorf("Expected valid JSON output, got error: %v", err)
				}
				// Check fields
				if data["name"] != "json-script" {
					t.Errorf("Expected name 'json-script', got %v", data["name"])
				}
			},
		},
		{
			name:       "Show script with source content (detail)",
			scriptID:   uuid.New().String(),
			asJSON:     false,
			withSource: true,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				scriptID := uuid.New()
				s := script.ConcreteScript(
					scriptID,
					"detail-script",
					createdAt,
					updatedAt,
				)
				ss.EXPECT().Get(ctx, gomock.Any()).Return(s, nil)

				// Expect Open to be called for content
				content := `print("Hello from Lua script")
local gogh = require("gogh")
print("Repository: " .. gogh.repo.name)`
				reader := &mockReadCloser{Reader: strings.NewReader(content)}
				ss.EXPECT().Open(ctx, scriptID.String()).Return(reader, nil)

				return ss
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				// Detail format should contain script info and content
				if !strings.Contains(output, "detail-script") {
					t.Error("Expected output to contain script name")
				}
				if !strings.Contains(output, "Hello from Lua script") {
					t.Error("Expected output to contain script content")
				}
			},
		},
		{
			name:       "Show script with source content as JSON",
			scriptID:   uuid.New().String(),
			asJSON:     true,
			withSource: true,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				scriptID := uuid.New()
				s := script.ConcreteScript(
					scriptID,
					"json-with-source",
					createdAt,
					updatedAt,
				)
				ss.EXPECT().Get(ctx, gomock.Any()).Return(s, nil)

				// Expect Open to be called for content
				content := `-- Lua script
function main()
    print("Running main function")
end
main()`
				reader := &mockReadCloser{Reader: strings.NewReader(content)}
				ss.EXPECT().Open(ctx, scriptID.String()).Return(reader, nil)

				return ss
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(output), &data); err != nil {
					t.Errorf("Expected valid JSON output, got error: %v", err)
				}
				// Check source field exists
				if source, ok := data["source"].(string); ok {
					if !strings.Contains(source, "Running main function") {
						t.Errorf("Expected source to contain 'Running main function', got %v", source)
					}
				} else {
					t.Error("Expected 'source' field in JSON output")
				}
			},
		},
		{
			name:       "Script not found",
			scriptID:   uuid.New().String(),
			asJSON:     false,
			withSource: false,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Get(ctx, gomock.Any()).Return(nil, errors.New("script not found"))
				return ss
			},
			wantErr: true,
		},
		{
			name:       "Invalid script ID",
			scriptID:   "invalid-id",
			asJSON:     false,
			withSource: false,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Get(ctx, "invalid-id").Return(nil, errors.New("invalid script ID"))
				return ss
			},
			wantErr: true,
		},
		{
			name:       "Empty script ID",
			scriptID:   "",
			asJSON:     false,
			withSource: false,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Get(ctx, "").Return(nil, errors.New("script ID is required"))
				return ss
			},
			wantErr: true,
		},
		{
			name:       "Error reading script content",
			scriptID:   uuid.New().String(),
			asJSON:     false,
			withSource: true,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				scriptID := uuid.New()
				s := script.ConcreteScript(
					scriptID,
					"error-script",
					createdAt,
					updatedAt,
				)
				ss.EXPECT().Get(ctx, gomock.Any()).Return(s, nil)

				// Open returns error
				ss.EXPECT().Open(ctx, scriptID.String()).Return(nil, errors.New("cannot read content"))

				return ss
			},
			wantErr: true,
		},
		{
			name:       "Show script with empty name",
			scriptID:   uuid.New().String(),
			asJSON:     true,
			withSource: false,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				s := script.ConcreteScript(
					uuid.New(),
					"", // Empty name
					createdAt,
					updatedAt,
				)
				ss.EXPECT().Get(ctx, gomock.Any()).Return(s, nil)
				return ss
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
			name:       "Show script with complex Lua content",
			scriptID:   uuid.New().String(),
			asJSON:     true,
			withSource: true,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				scriptID := uuid.New()
				s := script.ConcreteScript(
					scriptID,
					"complex-script",
					createdAt,
					updatedAt,
				)
				ss.EXPECT().Get(ctx, gomock.Any()).Return(s, nil)

				// Complex Lua script content
				content := `local gogh = require("gogh")
local json = require("json")
local http = require("http")

-- Repository information
local repo = gogh.repo

-- HTTP request example
function fetch_api_data()
    local response = http.get("https://api.example.com/data")
    return json.decode(response.body)
end

-- Main processing
local data = fetch_api_data()
print(json.encode({
    repo = repo.name,
    data = data
}))`
				reader := &mockReadCloser{Reader: strings.NewReader(content)}
				ss.EXPECT().Open(ctx, scriptID.String()).Return(reader, nil)

				return ss
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(output), &data); err != nil {
					t.Errorf("Expected valid JSON output, got error: %v", err)
				}
				source := data["source"].(string)
				if !strings.Contains(source, "fetch_api_data") {
					t.Error("Expected source to contain complex Lua function")
				}
			},
		},
		{
			name:       "Show script with special characters",
			scriptID:   uuid.New().String(),
			asJSON:     false,
			withSource: true,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				scriptID := uuid.New()
				s := script.ConcreteScript(
					scriptID,
					"unicode-script",
					createdAt,
					updatedAt,
				)
				ss.EXPECT().Get(ctx, gomock.Any()).Return(s, nil)

				// Content with special characters
				content := `-- Unicode test 🚀
print("Hello 世界")
print("Special chars: \n\t\\")`
				reader := &mockReadCloser{Reader: strings.NewReader(content)}
				ss.EXPECT().Open(ctx, scriptID.String()).Return(reader, nil)

				return ss
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "Unicode test 🚀") {
					t.Error("Expected output to contain unicode emoji")
				}
				if !strings.Contains(output, "Hello 世界") {
					t.Error("Expected output to contain Japanese characters")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var buf bytes.Buffer
			ss := tc.setupMock(ctrl)
			uc := script_show.NewUseCase(ss, &buf)

			err := uc.Execute(ctx, tc.scriptID, tc.asJSON, tc.withSource)
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
	ss := script_mock.NewMockScriptService(ctrl)
	uc := script_show.NewUseCase(ss, &buf)

	// Test service returning unexpected error
	ss.EXPECT().Get(ctx, "test-id").Return(nil, errors.New("database connection error"))

	err := uc.Execute(ctx, "test-id", false, false)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "get script by ID") {
		t.Errorf("Expected error to contain 'get script by ID', got %v", err.Error())
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
		{"JSONWithSource", true, true},
	}

	for _, mode := range modes {
		t.Run(mode.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var buf bytes.Buffer
			ss := script_mock.NewMockScriptService(ctrl)
			uc := script_show.NewUseCase(ss, &buf)

			scriptID := uuid.New()
			s := script.ConcreteScript(
				scriptID,
				"test-script",
				time.Now(),
				time.Now(),
			)
			ss.EXPECT().Get(ctx, gomock.Any()).Return(s, nil)

			if mode.withSource {
				// Expect content to be read
				reader := &mockReadCloser{Reader: strings.NewReader("test script content")}
				ss.EXPECT().Open(ctx, scriptID.String()).Return(reader, nil)
			}

			err := uc.Execute(ctx, scriptID.String(), mode.asJSON, mode.withSource)
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
