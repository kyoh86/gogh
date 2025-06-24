package update_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/script/update"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/script_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name       string
		scriptID   string
		scriptName string
		content    string
		setupMock  func(*gomock.Controller) *script_mock.MockScriptService
		wantErr    bool
	}{
		{
			name:       "Successfully update script with name and content",
			scriptID:   uuid.New().String(),
			scriptName: "updated-script",
			content:    `print("Updated script content")`,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry script.Entry) error {
						// Validate UUID format
						if _, err := uuid.Parse(id); err != nil {
							t.Errorf("Expected valid UUID, got %s", id)
						}
						// Validate name
						if entry.Name != "updated-script" {
							t.Errorf("Expected name 'updated-script', got %s", entry.Name)
						}
						// Verify content can be read
						buf := new(bytes.Buffer)
						if _, err := buf.ReadFrom(entry.Content); err != nil {
							t.Errorf("Failed to read content: %v", err)
						}
						if buf.String() != `print("Updated script content")` {
							t.Errorf("Expected specific content, got %s", buf.String())
						}
						return nil
					},
				)
				return ss
			},
			wantErr: false,
		},
		{
			name:       "Update script with empty name",
			scriptID:   uuid.New().String(),
			scriptName: "",
			content:    `print("Content without name")`,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry script.Entry) error {
						if entry.Name != "" {
							t.Errorf("Expected empty name, got %s", entry.Name)
						}
						return nil
					},
				)
				return ss
			},
			wantErr: false,
		},
		{
			name:       "Update script with empty content",
			scriptID:   uuid.New().String(),
			scriptName: "empty-content-script",
			content:    "",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry script.Entry) error {
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
				return ss
			},
			wantErr: false,
		},
		{
			name:       "Update script with complex Lua content",
			scriptID:   uuid.New().String(),
			scriptName: "complex-script",
			content: `local gogh = require("gogh")
local json = require("json")
local http = require("http")

-- Complex script with multiple functions
function fetch_data(url)
    local response, err = http.get(url)
    if err then
        error("Failed to fetch: " .. err)
    end
    return response.body
end

function process_repository()
    local repo = gogh.repo
    local metadata = {
        name = repo.name,
        owner = repo.owner,
        host = repo.host,
        timestamp = os.time()
    }
    
    -- Process and log
    print("Processing: " .. json.encode(metadata))
    
    return metadata
end

-- Main execution
local result = process_repository()
print("Completed: " .. json.encode(result))`,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry script.Entry) error {
						// Verify we can read the entire content
						buf := new(bytes.Buffer)
						n, err := buf.ReadFrom(entry.Content)
						if err != nil {
							t.Errorf("Failed to read content: %v", err)
						}
						if n == 0 {
							t.Error("Expected non-empty content")
						}
						// Check that name is correct
						if entry.Name != "complex-script" {
							t.Errorf("Expected name 'complex-script', got %s", entry.Name)
						}
						return nil
					},
				)
				return ss
			},
			wantErr: false,
		},
		{
			name:       "Update script with large content",
			scriptID:   uuid.New().String(),
			scriptName: "large-script",
			content:    strings.Repeat("-- This is a comment line\nprint('Processing item...')\n", 1000),
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry script.Entry) error {
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
				return ss
			},
			wantErr: false,
		},
		{
			name:       "Update non-existent script",
			scriptID:   uuid.New().String(),
			scriptName: "non-existent",
			content:    "content",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(errors.New("script not found"))
				return ss
			},
			wantErr: true,
		},
		{
			name:       "Update with invalid script ID",
			scriptID:   "invalid-id",
			scriptName: "test",
			content:    "content",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, "invalid-id", gomock.Any()).Return(errors.New("invalid script ID"))
				return ss
			},
			wantErr: true,
		},
		{
			name:       "Update with empty script ID",
			scriptID:   "",
			scriptName: "test",
			content:    "content",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, "", gomock.Any()).Return(errors.New("script ID is required"))
				return ss
			},
			wantErr: true,
		},
		{
			name:       "Update with special characters in name",
			scriptID:   uuid.New().String(),
			scriptName: "script-with-special-chars-🚀",
			content:    `print("Script with special name")`,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry script.Entry) error {
						if entry.Name != "script-with-special-chars-🚀" {
							t.Errorf("Expected name with special chars, got %s", entry.Name)
						}
						return nil
					},
				)
				return ss
			},
			wantErr: false,
		},
		{
			name:       "Update with special characters in content",
			scriptID:   uuid.New().String(),
			scriptName: "unicode-script",
			content: `-- Script with unicode characters
print("Hello 世界")
print("Emoji support: 🎉 🚀 ✨")
print("Special chars: \n\t\r\\")`,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(nil)
				return ss
			},
			wantErr: false,
		},
		{
			name:       "Service returns unexpected error",
			scriptID:   uuid.New().String(),
			scriptName: "test",
			content:    "content",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(errors.New("storage error"))
				return ss
			},
			wantErr: true,
		},
		{
			name:       "Update both name and content empty",
			scriptID:   uuid.New().String(),
			scriptName: "",
			content:    "",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry script.Entry) error {
						// Service might validate and return error
						return errors.New("at least name or content must be provided")
					},
				)
				return ss
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ss := tc.setupMock(ctrl)
			uc := testtarget.NewUseCase(ss)

			err := uc.Execute(ctx, tc.scriptID, tc.scriptName, strings.NewReader(tc.content))
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
	customReader := strings.NewReader("test script content")

	ss := script_mock.NewMockScriptService(ctrl)
	ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, id string, entry script.Entry) error {
			// Verify that the same reader instance is passed
			if entry.Content != customReader {
				t.Error("Expected the same reader instance to be passed")
			}
			// Verify name is passed correctly
			if entry.Name != "test-script" {
				t.Errorf("Expected name 'test-script', got %s", entry.Name)
			}
			return nil
		},
	)

	uc := testtarget.NewUseCase(ss)
	err := uc.Execute(ctx, uuid.New().String(), "test-script", customReader)
	if err != nil {
		t.Errorf("Execute() unexpected error = %v", err)
	}
}

func TestUseCase_Execute_MultipleReaders(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ss := script_mock.NewMockScriptService(ctrl)
	uc := testtarget.NewUseCase(ss)

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
		ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, id string, entry script.Entry) error {
				// Verify content can be read from any reader type
				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(entry.Content)
				if err != nil {
					t.Errorf("%s: Failed to read content: %v", r.name, err)
				}
				return nil
			},
		)

		err := uc.Execute(ctx, uuid.New().String(), r.name, r.reader)
		if err != nil {
			t.Errorf("%s: Execute() unexpected error = %v", r.name, err)
		}
	}
}
