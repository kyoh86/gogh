package config_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/script_mock"
	"go.uber.org/mock/gomock"
)

func TestScriptStore_Source(t *testing.T) {
	store := config.NewScriptStore()

	// Test with no environment variable
	source, err := store.Source()
	if err != nil {
		t.Fatalf("Source() error = %v", err)
	}

	// Should contain "script.v4.toml"
	if !contains(source, "script.v4.toml") {
		t.Errorf("Expected source to contain 'script.v4.toml', got %s", source)
	}

	// Test with environment variable
	tempDir := t.TempDir()
	customPath := filepath.Join(tempDir, "custom-script.toml")
	t.Setenv("GOGH_SCRIPT_PATH", customPath)

	source, err = store.Source()
	if err != nil {
		t.Fatalf("Source() with env var error = %v", err)
	}

	if source != customPath {
		t.Errorf("Expected source to be %s, got %s", customPath, source)
	}
}

func TestScriptStore_Load(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		setupFile func(t *testing.T) string
		setupMock func(*gomock.Controller) *script_mock.MockScriptService
		wantErr   bool
		wantEmpty bool
	}{
		{
			name: "Load from non-existent file",
			setupFile: func(t *testing.T) string {
				tempDir := t.TempDir()
				return filepath.Join(tempDir, "non-existent.toml")
			},
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().MarkSaved()
				return ss
			},
			wantErr:   false,
			wantEmpty: true,
		},
		{
			name: "Load valid scripts",
			setupFile: func(t *testing.T) string {
				tempDir := t.TempDir()
				path := filepath.Join(tempDir, "script.toml")

				content := `[[scripts]]
id = "` + uuid.New().String() + `"
name = "deploy-script"
created-at = 2023-01-01T00:00:00Z
updated-at = 2023-01-02T00:00:00Z

[[scripts]]
id = "` + uuid.New().String() + `"
name = "test-runner"
created-at = 2023-02-01T00:00:00Z
updated-at = 2023-02-02T00:00:00Z
`
				if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				return path
			},
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)

				// Expect Load to be called with proper scripts
				ss.EXPECT().Load(gomock.Any()).DoAndReturn(func(yield func(func(script.Script, error) bool)) error {
					loadedCount := 0
					yield(func(s script.Script, err error) bool {
						if err != nil {
							t.Errorf("Unexpected error in yield: %v", err)
							return false
						}
						loadedCount++

						// Verify loaded scripts
						switch loadedCount {
						case 1:
							if s.Name() != "deploy-script" {
								t.Errorf("Expected name 'deploy-script', got %s", s.Name())
							}
							expectedCreated := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
							if !s.CreatedAt().Equal(expectedCreated) {
								t.Errorf("Expected created at %v, got %v", expectedCreated, s.CreatedAt())
							}
						case 2:
							if s.Name() != "test-runner" {
								t.Errorf("Expected name 'test-runner', got %s", s.Name())
							}
						}
						return true
					})

					if loadedCount != 2 {
						t.Errorf("Expected 2 scripts to be loaded, got %d", loadedCount)
					}
					return nil
				})

				ss.EXPECT().MarkSaved()

				return ss
			},
			wantErr: false,
		},
		{
			name: "Load with invalid TOML",
			setupFile: func(t *testing.T) string {
				tempDir := t.TempDir()
				path := filepath.Join(tempDir, "invalid.toml")

				content := `invalid toml content {[}]`
				if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				return path
			},
			setupMock: script_mock.NewMockScriptService,
			wantErr:   true,
		},
		{
			name: "Load with service error",
			setupFile: func(t *testing.T) string {
				tempDir := t.TempDir()
				path := filepath.Join(tempDir, "script.toml")

				content := `[[scripts]]
id = "` + uuid.New().String() + `"
name = "test-script"
created-at = 2023-01-01T00:00:00Z
updated-at = 2023-01-01T00:00:00Z
`
				if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				return path
			},
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().Load(gomock.Any()).Return(os.ErrPermission)
				return ss
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			path := tc.setupFile(t)
			t.Setenv("GOGH_SCRIPT_PATH", path)

			store := config.NewScriptStore()
			mockService := tc.setupMock(ctrl)

			initial := func() script.ScriptService {
				return mockService
			}

			svc, err := store.Load(ctx, initial)
			if (err != nil) != tc.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && svc == nil {
				t.Error("Expected service to be returned")
			}
		})
	}
}

func TestScriptStore_Save(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		setupMock func(*gomock.Controller) *script_mock.MockScriptService
		force     bool
		wantErr   bool
		validate  func(t *testing.T, path string)
	}{
		{
			name: "Save with no changes and no force",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().HasChanges().Return(false)
				return ss
			},
			force:   false,
			wantErr: false,
			validate: func(t *testing.T, path string) {
				// File should not be created
				if _, err := os.Stat(path); !os.IsNotExist(err) {
					t.Error("Expected file to not exist when no changes")
				}
			},
		},
		{
			name: "Save with changes",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().HasChanges().Return(true)

				// Create test scripts
				script1 := script.ConcreteScript(
					uuid.New(),
					"deploy-script",
					time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				)

				script2 := script.ConcreteScript(
					uuid.New(),
					"test-runner",
					time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2023, 2, 2, 0, 0, 0, 0, time.UTC),
				)

				ss.EXPECT().List().Return(func(yield func(script.Script, error) bool) {
					yield(script1, nil)
					yield(script2, nil)
				})

				ss.EXPECT().MarkSaved()

				return ss
			},
			force:   false,
			wantErr: false,
			validate: func(t *testing.T, path string) {
				// File should be created
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("Failed to read saved file: %v", err)
				}

				// Check content contains expected data
				contentStr := string(content)
				if !contains(contentStr, "deploy-script") {
					t.Error("Expected script name 'deploy-script' in saved content")
				}
				if !contains(contentStr, "test-runner") {
					t.Error("Expected script name 'test-runner' in saved content")
				}
				if !contains(contentStr, "created-at") {
					t.Error("Expected created-at field in saved content")
				}
			},
		},
		{
			name: "Save with force",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				// Even with force, HasChanges is still called
				ss.EXPECT().HasChanges().Return(false)

				testScript := script.ConcreteScript(
					uuid.New(),
					"force-save-script",
					time.Now(),
					time.Now(),
				)

				ss.EXPECT().List().Return(func(yield func(script.Script, error) bool) {
					yield(testScript, nil)
				})

				ss.EXPECT().MarkSaved()

				return ss
			},
			force:   true,
			wantErr: false,
			validate: func(t *testing.T, path string) {
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("Failed to read saved file: %v", err)
				}

				if !contains(string(content), "force-save-script") {
					t.Error("Expected force-saved script in content")
				}
			},
		},
		{
			name: "Save with list error",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().HasChanges().Return(true)

				ss.EXPECT().List().Return(func(yield func(script.Script, error) bool) {
					yield(nil, os.ErrPermission)
				})

				return ss
			},
			force:   false,
			wantErr: true,
		},
		{
			name: "Save with empty script name",
			setupMock: func(ctrl *gomock.Controller) *script_mock.MockScriptService {
				ss := script_mock.NewMockScriptService(ctrl)
				ss.EXPECT().HasChanges().Return(true)

				emptyNameScript := script.ConcreteScript(
					uuid.New(),
					"", // Empty name
					time.Now(),
					time.Now(),
				)

				ss.EXPECT().List().Return(func(yield func(script.Script, error) bool) {
					yield(emptyNameScript, nil)
				})

				ss.EXPECT().MarkSaved()

				return ss
			},
			force:   false,
			wantErr: false,
			validate: func(t *testing.T, path string) {
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("Failed to read saved file: %v", err)
				}

				// Should still save with empty name
				if !contains(string(content), "name = ") {
					t.Error("Expected name field in saved content")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tempDir := t.TempDir()
			path := filepath.Join(tempDir, "script.toml")
			t.Setenv("GOGH_SCRIPT_PATH", path)

			store := config.NewScriptStore()
			mockService := tc.setupMock(ctrl)

			err := store.Save(ctx, mockService, tc.force)
			if (err != nil) != tc.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && tc.validate != nil {
				tc.validate(t, path)
			}
		})
	}
}

func TestScriptDir(t *testing.T) {
	// Test with no environment variable
	dir, err := config.ScriptDir()
	if err != nil {
		t.Fatalf("ScriptDir() error = %v", err)
	}

	if !contains(dir, "script.v4.toml") {
		t.Errorf("Expected dir to contain 'script.v4.toml', got %s", dir)
	}

	// Test with environment variable
	customPath := "/custom/path/script.toml"
	t.Setenv("GOGH_SCRIPT_PATH", customPath)

	dir, err = config.ScriptDir()
	if err != nil {
		t.Fatalf("ScriptDir() with env var error = %v", err)
	}

	if dir != customPath {
		t.Errorf("Expected dir to be %s, got %s", customPath, dir)
	}
}
