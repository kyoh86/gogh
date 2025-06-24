package script_edit_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/script_edit"
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

func TestUseCase_ExtractScript(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		scriptID  string
		setupMock func(*gomock.Controller) *script_mock.MockScriptService
		wantErr   bool
		wantData  string
	}{
		{
			name:     "Successfully extract script",
			scriptID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				content := `print("Hello from Lua script")`
				reader := &mockReadCloser{Reader: strings.NewReader(content)}
				ss.EXPECT().Open(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string) (io.ReadCloser, error) {
						// Validate UUID format
						if _, err := uuid.Parse(id); err != nil {
							t.Errorf("Expected valid UUID, got %s", id)
						}
						return reader, nil
					},
				)
				return ss
			},
			wantErr:  false,
			wantData: `print("Hello from Lua script")`,
		},
		{
			name:     "Extract script with complex content",
			scriptID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				content := `local gogh = require("gogh")
local repo = gogh.repo

-- Complex script with multiple functions
function setup()
    print("Setting up repository: " .. repo.name)
    -- More complex logic here
end

function cleanup()
    print("Cleaning up")
end

setup()
cleanup()`
				reader := &mockReadCloser{Reader: strings.NewReader(content)}
				ss.EXPECT().Open(ctx, gomock.Any()).Return(reader, nil)
				return ss
			},
			wantErr: false,
			wantData: `local gogh = require("gogh")
local repo = gogh.repo

-- Complex script with multiple functions
function setup()
    print("Setting up repository: " .. repo.name)
    -- More complex logic here
end

function cleanup()
    print("Cleaning up")
end

setup()
cleanup()`,
		},
		{
			name:     "Extract large script",
			scriptID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				// Simulate a large script with repeated lines
				largeContent := strings.Repeat("-- This is a comment line\nprint('Processing...')\n", 500)
				reader := &mockReadCloser{Reader: strings.NewReader(largeContent)}
				ss.EXPECT().Open(ctx, gomock.Any()).Return(reader, nil)
				return ss
			},
			wantErr:  false,
			wantData: strings.Repeat("-- This is a comment line\nprint('Processing...')\n", 500),
		},
		{
			name:     "Extract empty script",
			scriptID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				reader := &mockReadCloser{Reader: strings.NewReader("")}
				ss.EXPECT().Open(ctx, gomock.Any()).Return(reader, nil)
				return ss
			},
			wantErr:  false,
			wantData: "",
		},
		{
			name:     "Script not found",
			scriptID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Open(ctx, gomock.Any()).Return(nil, errors.New("script not found"))
				return ss
			},
			wantErr: true,
		},
		{
			name:     "Invalid script ID",
			scriptID: "invalid-id",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Open(ctx, "invalid-id").Return(nil, errors.New("invalid script ID"))
				return ss
			},
			wantErr: true,
		},
		{
			name:     "Empty script ID",
			scriptID: "",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Open(ctx, "").Return(nil, errors.New("script ID is required"))
				return ss
			},
			wantErr: true,
		},
		{
			name:     "Service returns unexpected error",
			scriptID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Open(ctx, gomock.Any()).Return(nil, errors.New("storage error"))
				return ss
			},
			wantErr: true,
		},
		{
			name:     "Extract script with special characters",
			scriptID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				content := `print("\n\t\r")
print("Special chars: 🚀 ✨ 💻")
print('Single quotes with "nested" double')`
				reader := &mockReadCloser{Reader: strings.NewReader(content)}
				ss.EXPECT().Open(ctx, gomock.Any()).Return(reader, nil)
				return ss
			},
			wantErr: false,
			wantData: `print("\n\t\r")
print("Special chars: 🚀 ✨ 💻")
print('Single quotes with "nested" double')`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ss := tc.setupMock(ctrl)
			uc := script_edit.NewUseCase(ss)

			var buf bytes.Buffer
			err := uc.ExtractScript(ctx, tc.scriptID, &buf)
			if (err != nil) != tc.wantErr {
				t.Errorf("ExtractScript() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr {
				if buf.String() != tc.wantData {
					t.Errorf("ExtractScript() data = %q, want %q", buf.String(), tc.wantData)
				}
			}
		})
	}
}

func TestUseCase_UpdateScript(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		scriptID  string
		content   string
		setupMock func(*gomock.Controller) *script_mock.MockScriptService
		wantErr   bool
	}{
		{
			name:     "Successfully update script",
			scriptID: uuid.New().String(),
			content:  `print("Updated Lua script")`,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry script.Entry) error {
						// Validate UUID format
						if _, err := uuid.Parse(id); err != nil {
							t.Errorf("Expected valid UUID, got %s", id)
						}
						// Verify content can be read
						buf := new(bytes.Buffer)
						if _, err := buf.ReadFrom(entry.Content); err != nil {
							t.Errorf("Failed to read content: %v", err)
						}
						if buf.String() != `print("Updated Lua script")` {
							t.Errorf("Expected content 'print(\"Updated Lua script\")', got %s", buf.String())
						}
						return nil
					},
				)
				return ss
			},
			wantErr: false,
		},
		{
			name:     "Update with empty content",
			scriptID: uuid.New().String(),
			content:  "",
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
			name:     "Update with complex script",
			scriptID: uuid.New().String(),
			content: `local gogh = require("gogh")
local json = require("json")

function process_repo()
    local repo = gogh.repo
    local data = {
        name = repo.name,
        owner = repo.owner,
        processed = true
    }
    print(json.encode(data))
end

process_repo()`,
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
			name:     "Update with large script",
			scriptID: uuid.New().String(),
			content:  strings.Repeat("-- Comment line\nprint('Processing...')\n", 1000),
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
			name:     "Update non-existent script",
			scriptID: uuid.New().String(),
			content:  "content",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(errors.New("script not found"))
				return ss
			},
			wantErr: true,
		},
		{
			name:     "Update with invalid script ID",
			scriptID: "invalid-id",
			content:  "content",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, "invalid-id", gomock.Any()).Return(errors.New("invalid script ID"))
				return ss
			},
			wantErr: true,
		},
		{
			name:     "Update with empty script ID",
			scriptID: "",
			content:  "content",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, "", gomock.Any()).Return(errors.New("script ID is required"))
				return ss
			},
			wantErr: true,
		},
		{
			name:     "Service returns unexpected error",
			scriptID: uuid.New().String(),
			content:  "content",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(errors.New("storage error"))
				return ss
			},
			wantErr: true,
		},
		{
			name:     "Update with special characters",
			scriptID: uuid.New().String(),
			content: `-- Script with special chars
print("Tab:\t Newline:\n Backslash:\\ Quote:\"")
print('Unicode: 🚀 ✨ 💻')`,
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(nil)
				return ss
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ss := tc.setupMock(ctrl)
			uc := script_edit.NewUseCase(ss)

			err := uc.UpdateScript(ctx, tc.scriptID, strings.NewReader(tc.content))
			if (err != nil) != tc.wantErr {
				t.Errorf("UpdateScript() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestUseCase_ExtractScript_ReaderBehavior(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test that the reader is properly closed
	ss := script_mock.NewMockScriptService(ctrl)
	reader := &mockReadCloser{Reader: strings.NewReader("test script content")}

	ss.EXPECT().Open(ctx, gomock.Any()).Return(reader, nil)

	uc := script_edit.NewUseCase(ss)
	var buf bytes.Buffer
	err := uc.ExtractScript(ctx, uuid.New().String(), &buf)
	if err != nil {
		t.Errorf("ExtractScript() unexpected error = %v", err)
	}

	if !reader.closed {
		t.Error("Expected reader to be closed")
	}
}

func TestUseCase_UpdateScript_ReaderPassthrough(t *testing.T) {
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
			return nil
		},
	)

	uc := script_edit.NewUseCase(ss)
	err := uc.UpdateScript(ctx, uuid.New().String(), customReader)
	if err != nil {
		t.Errorf("UpdateScript() unexpected error = %v", err)
	}
}
