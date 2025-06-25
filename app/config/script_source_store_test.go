package config_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/config"
)

func TestScriptSourceStore_Source(t *testing.T) {
	store := config.NewScriptSourceStore()

	// Test with no environment variable
	source, err := store.Source()
	if err != nil {
		t.Fatalf("Source() error = %v", err)
	}

	// Should contain "script.v4"
	if !contains(source, "script.v4") {
		t.Errorf("Expected source to contain 'script.v4', got %s", source)
	}

	// Test with environment variable
	tempDir := t.TempDir()
	customPath := filepath.Join(tempDir, "custom-script-content")
	t.Setenv("GOGH_HOOK_CONTENT_PATH", customPath)

	source, err = store.Source()
	if err != nil {
		t.Fatalf("Source() with env var error = %v", err)
	}

	if source != customPath {
		t.Errorf("Expected source to be %s, got %s", customPath, source)
	}

	// Test error path by overriding AppContextPathFunc
	originalFunc := config.AppContextPathFunc
	defer func() {
		config.AppContextPathFunc = originalFunc
	}()

	config.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
		return "", os.ErrPermission
	}

	_, err = store.Source()
	if err == nil {
		t.Error("Expected error when AppContextPathFunc fails")
	}
	if !contains(err.Error(), "search script content path") {
		t.Errorf("Expected error to contain 'search script content path', got %v", err)
	}
}

func TestScriptSourceStore_Save(t *testing.T) {
	ctx := context.Background()

	// Test error when Source() fails
	t.Run("Save with Source error", func(t *testing.T) {
		originalFunc := config.AppContextPathFunc
		defer func() {
			config.AppContextPathFunc = originalFunc
		}()

		config.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
			return "", os.ErrPermission
		}

		store := config.NewScriptSourceStore()
		err := store.Save(ctx, "test-id", bytes.NewReader([]byte("content")))
		if err == nil {
			t.Error("Expected error when Source() fails")
		}
		if !contains(err.Error(), "get content source") {
			t.Errorf("Expected error to contain 'get content source', got %v", err)
		}
	})

	testCases := []struct {
		name     string
		scriptID string
		content  string
		setup    func(t *testing.T) string
		wantErr  bool
		validate func(t *testing.T, path string, scriptID string)
	}{
		{
			name:     "Save script content successfully",
			scriptID: uuid.New().String(),
			content:  "print('Hello from Lua script')\nlocal gogh = require('gogh')\nprint(gogh.repo.name)",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: false,
			validate: func(t *testing.T, path string, scriptID string) {
				filePath := filepath.Join(path, scriptID)
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Failed to read saved file: %v", err)
				}
				expected := "print('Hello from Lua script')\nlocal gogh = require('gogh')\nprint(gogh.repo.name)"
				if string(content) != expected {
					t.Errorf("Expected content %q, got %q", expected, string(content))
				}
			},
		},
		{
			name:     "Save empty content",
			scriptID: uuid.New().String(),
			content:  "",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: false,
			validate: func(t *testing.T, path string, scriptID string) {
				filePath := filepath.Join(path, scriptID)
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Failed to read saved file: %v", err)
				}
				if len(content) != 0 {
					t.Errorf("Expected empty content, got %q", string(content))
				}
			},
		},
		{
			name:     "Save with empty script ID",
			scriptID: "",
			content:  "test content",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: true, // Should fail because empty filename results in directory path
		},
		{
			name:     "Save with binary content",
			scriptID: uuid.New().String(),
			content:  string([]byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}),
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: false,
			validate: func(t *testing.T, path string, scriptID string) {
				filePath := filepath.Join(path, scriptID)
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Failed to read saved file: %v", err)
				}
				expected := []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}
				if !bytes.Equal(content, expected) {
					t.Errorf("Binary content mismatch")
				}
			},
		},
		{
			name:     "Save with special characters",
			scriptID: uuid.New().String(),
			content:  "-- Unicode test 🚀\nprint('Hello 世界')\nprint('Special: \\n\\t\\\\')",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: false,
			validate: func(t *testing.T, path string, scriptID string) {
				filePath := filepath.Join(path, scriptID)
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Failed to read saved file: %v", err)
				}
				if !contains(string(content), "🚀") {
					t.Error("Expected emoji in content")
				}
				if !contains(string(content), "世界") {
					t.Error("Expected Japanese characters in content")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := tc.setup(t)
			t.Setenv("GOGH_HOOK_CONTENT_PATH", tempDir)

			store := config.NewScriptSourceStore()
			reader := bytes.NewReader([]byte(tc.content))

			err := store.Save(ctx, tc.scriptID, reader)
			if (err != nil) != tc.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && tc.validate != nil {
				tc.validate(t, tempDir, tc.scriptID)
			}
		})
	}
}

func TestScriptSourceStore_Open(t *testing.T) {
	ctx := context.Background()

	// Test error when Source() fails
	t.Run("Open with Source error", func(t *testing.T) {
		originalFunc := config.AppContextPathFunc
		defer func() {
			config.AppContextPathFunc = originalFunc
		}()

		config.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
			return "", os.ErrPermission
		}

		store := config.NewScriptSourceStore()
		_, err := store.Open(ctx, "test-id")
		if err == nil {
			t.Error("Expected error when Source() fails")
		}
		if !contains(err.Error(), "get content source") {
			t.Errorf("Expected error to contain 'get content source', got %v", err)
		}
	})

	testCases := []struct {
		name     string
		scriptID string
		setup    func(t *testing.T) string
		wantErr  bool
		validate func(t *testing.T, rc io.ReadCloser)
	}{
		{
			name:     "Open existing script",
			scriptID: uuid.New().String(),
			setup: func(t *testing.T) string {
				tempDir := t.TempDir()
				scriptID := uuid.New().String()
				content := "test script content"
				err := os.WriteFile(filepath.Join(tempDir, scriptID), []byte(content), 0o644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				// Return the scriptID through test context
				t.Cleanup(func() {
					// Store scriptID for use in test
					t.Setenv("_TEST_SCRIPT_ID", scriptID)
				})
				t.Setenv("_TEST_SCRIPT_ID", scriptID)
				return tempDir
			},
			wantErr: false,
			validate: func(t *testing.T, rc io.ReadCloser) {
				defer rc.Close()
				content, err := io.ReadAll(rc)
				if err != nil {
					t.Fatalf("Failed to read content: %v", err)
				}
				if string(content) != "test script content" {
					t.Errorf("Expected content 'test script content', got %q", string(content))
				}
			},
		},
		{
			name:     "Open non-existent script",
			scriptID: uuid.New().String(),
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: true,
		},
		{
			name:     "Open with empty script ID",
			scriptID: "",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			wantErr: false, // os.Open on a directory succeeds, but reading will fail
			validate: func(t *testing.T, rc io.ReadCloser) {
				defer rc.Close()
				// Attempting to read from a directory should fail
				_, err := io.ReadAll(rc)
				if err == nil {
					t.Error("Expected error when reading from directory")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := tc.setup(t)
			t.Setenv("GOGH_HOOK_CONTENT_PATH", tempDir)

			// For the first test case, get the actual scriptID that was created
			scriptID := tc.scriptID
			if envID := os.Getenv("_TEST_SCRIPT_ID"); envID != "" && tc.name == "Open existing script" {
				scriptID = envID
			}

			store := config.NewScriptSourceStore()
			rc, err := store.Open(ctx, scriptID)

			if (err != nil) != tc.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && tc.validate != nil {
				tc.validate(t, rc)
			}
		})
	}
}

func TestScriptSourceStore_Remove(t *testing.T) {
	ctx := context.Background()

	// Test error when Source() fails
	t.Run("Remove with Source error", func(t *testing.T) {
		originalFunc := config.AppContextPathFunc
		defer func() {
			config.AppContextPathFunc = originalFunc
		}()

		config.AppContextPathFunc = func(envar string, getDir func() (string, error), rel ...string) (string, error) {
			return "", os.ErrPermission
		}

		store := config.NewScriptSourceStore()
		err := store.Remove(ctx, "test-id")
		if err == nil {
			t.Error("Expected error when Source() fails")
		}
		if !contains(err.Error(), "get content source") {
			t.Errorf("Expected error to contain 'get content source', got %v", err)
		}
	})

	testCases := []struct {
		name     string
		scriptID string
		setup    func(t *testing.T) (string, string) // returns tempDir and scriptID
		wantErr  bool
		validate func(t *testing.T, path string, scriptID string)
	}{
		{
			name: "Remove existing script",
			setup: func(t *testing.T) (string, string) {
				tempDir := t.TempDir()
				scriptID := uuid.New().String()
				content := "test script to remove"
				err := os.WriteFile(filepath.Join(tempDir, scriptID), []byte(content), 0o644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return tempDir, scriptID
			},
			wantErr: false,
			validate: func(t *testing.T, path string, scriptID string) {
				filePath := filepath.Join(path, scriptID)
				if _, err := os.Stat(filePath); !os.IsNotExist(err) {
					t.Error("Expected file to be removed")
				}
			},
		},
		{
			name: "Remove non-existent script",
			setup: func(t *testing.T) (string, string) {
				return t.TempDir(), uuid.New().String()
			},
			wantErr: true,
		},
		{
			name: "Remove with empty script ID",
			setup: func(t *testing.T) (string, string) {
				// Create a subdirectory to avoid removing the temp dir itself
				tempDir := t.TempDir()
				subDir := filepath.Join(tempDir, "scripts")
				if err := os.MkdirAll(subDir, 0o755); err != nil {
					t.Fatalf("Failed to create subdirectory: %v", err)
				}
				return subDir, ""
			},
			wantErr: false, // os.Remove will succeed on removing the directory
			validate: func(t *testing.T, path string, scriptID string) {
				// The directory should have been removed
				if _, err := os.Stat(path); !os.IsNotExist(err) {
					t.Error("Expected directory to be removed")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pathToUse, scriptID := tc.setup(t)
			t.Setenv("GOGH_HOOK_CONTENT_PATH", pathToUse)

			store := config.NewScriptSourceStore()
			err := store.Remove(ctx, scriptID)

			if (err != nil) != tc.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && tc.validate != nil {
				tc.validate(t, pathToUse, scriptID)
			}
		})
	}
}

func TestScriptSourceStore_Integration(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir()
	t.Setenv("GOGH_HOOK_CONTENT_PATH", tempDir)

	store := config.NewScriptSourceStore()
	scriptID := uuid.New().String()
	content := "-- Integration test\nprint('Hello, World!')"

	// Save
	err := store.Save(ctx, scriptID, bytes.NewReader([]byte(content)))
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Open and verify
	rc, err := store.Open(ctx, scriptID)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	readContent, err := io.ReadAll(rc)
	rc.Close()
	if err != nil {
		t.Fatalf("Failed to read content: %v", err)
	}

	if string(readContent) != content {
		t.Errorf("Content mismatch: expected %q, got %q", content, string(readContent))
	}

	// Remove
	err = store.Remove(ctx, scriptID)
	if err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	// Verify removed
	_, err = store.Open(ctx, scriptID)
	if err == nil {
		t.Error("Expected error when opening removed script")
	}
}
