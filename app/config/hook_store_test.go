package config_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/hook_mock"
	"go.uber.org/mock/gomock"
)

func TestHookStore_Source(t *testing.T) {
	store := config.NewHookStore()

	// Test with no environment variable
	source, err := store.Source()
	if err != nil {
		t.Fatalf("Source() error = %v", err)
	}

	// Should contain "hook.v4.toml"
	if !contains(source, "hook.v4.toml") {
		t.Errorf("Expected source to contain 'hook.v4.toml', got %s", source)
	}

	// Test with environment variable
	tempDir := t.TempDir()
	customPath := filepath.Join(tempDir, "custom-hook.toml")
	t.Setenv("GOGH_HOOK_PATH", customPath)

	source, err = store.Source()
	if err != nil {
		t.Fatalf("Source() with env var error = %v", err)
	}

	if source != customPath {
		t.Errorf("Expected source to be %s, got %s", customPath, source)
	}
}

func TestHookStore_Load(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		setupFile func(t *testing.T) string
		setupMock func(*gomock.Controller) *hook_mock.MockHookService
		wantErr   bool
		wantEmpty bool
	}{
		{
			name: "Load from non-existent file",
			setupFile: func(t *testing.T) string {
				tempDir := t.TempDir()
				return filepath.Join(tempDir, "non-existent.toml")
			},
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().MarkSaved()
				return hs
			},
			wantErr:   false,
			wantEmpty: true,
		},
		{
			name: "Load valid hooks",
			setupFile: func(t *testing.T) string {
				tempDir := t.TempDir()
				path := filepath.Join(tempDir, "hook.toml")

				content := `[[hooks]]
id = "` + uuid.New().String() + `"
name = "post-clone-hook"
repo-pattern = "github.com/owner/*"
trigger-event = "post-clone"
operation-type = "overlay"
operation-id = "overlay-123"

[[hooks]]
id = "` + uuid.New().String() + `"
name = "post-create-script"
repo-pattern = "github.com/myorg/*"
trigger-event = "post-create"
operation-type = "script"
operation-id = "script-456"
`
				if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				return path
			},
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)

				// Expect Load to be called with proper hooks
				hs.EXPECT().Load(gomock.Any()).DoAndReturn(func(yield func(func(hook.Hook, error) bool)) error {
					loadedCount := 0
					yield(func(h hook.Hook, err error) bool {
						if err != nil {
							t.Errorf("Unexpected error in yield: %v", err)
							return false
						}
						loadedCount++

						// Verify loaded hooks
						switch loadedCount {
						case 1:
							if h.Name() != "post-clone-hook" {
								t.Errorf("Expected name 'post-clone-hook', got %s", h.Name())
							}
							if h.RepoPattern() != "github.com/owner/*" {
								t.Errorf("Expected pattern 'github.com/owner/*', got %s", h.RepoPattern())
							}
							if h.TriggerEvent() != hook.EventPostClone {
								t.Errorf("Expected event post-clone, got %s", h.TriggerEvent())
							}
							if h.OperationType() != hook.OperationTypeOverlay {
								t.Errorf("Expected operation type overlay, got %s", h.OperationType())
							}
						case 2:
							if h.Name() != "post-create-script" {
								t.Errorf("Expected name 'post-create-script', got %s", h.Name())
							}
							if h.TriggerEvent() != hook.EventPostCreate {
								t.Errorf("Expected event post-create, got %s", h.TriggerEvent())
							}
							if h.OperationType() != hook.OperationTypeScript {
								t.Errorf("Expected operation type script, got %s", h.OperationType())
							}
						}
						return true
					})

					if loadedCount != 2 {
						t.Errorf("Expected 2 hooks to be loaded, got %d", loadedCount)
					}
					return nil
				})

				hs.EXPECT().MarkSaved()

				return hs
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
			setupMock: hook_mock.NewMockHookService,
			wantErr:   true,
		},
		{
			name: "Load with service error",
			setupFile: func(t *testing.T) string {
				tempDir := t.TempDir()
				path := filepath.Join(tempDir, "hook.toml")

				content := `[[hooks]]
id = "` + uuid.New().String() + `"
name = "test-hook"
repo-pattern = "*"
trigger-event = "post-clone"
operation-type = "overlay"
operation-id = "test-id"
`
				if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				return path
			},
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Load(gomock.Any()).Return(os.ErrPermission)
				return hs
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			path := tc.setupFile(t)
			t.Setenv("GOGH_HOOK_PATH", path)

			store := config.NewHookStore()
			mockService := tc.setupMock(ctrl)

			initial := func() hook.HookService {
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

func TestHookStore_Save(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		setupMock func(*gomock.Controller) *hook_mock.MockHookService
		force     bool
		wantErr   bool
		validate  func(t *testing.T, path string)
	}{
		{
			name: "Save with no changes and no force",
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().HasChanges().Return(false)
				return hs
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
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().HasChanges().Return(true)

				// Create test hooks
				hook1 := hook.ConcreteHook(
					uuid.New(),
					"post-clone-overlay",
					"github.com/owner/*",
					string(hook.EventPostClone),
					string(hook.OperationTypeOverlay),
					"overlay-id-1",
				)

				hook2 := hook.ConcreteHook(
					uuid.New(),
					"post-fork-script",
					"github.com/test/*",
					string(hook.EventPostFork),
					string(hook.OperationTypeScript),
					"script-id-2",
				)

				hs.EXPECT().List().Return(func(yield func(hook.Hook, error) bool) {
					yield(hook1, nil)
					yield(hook2, nil)
				})

				hs.EXPECT().MarkSaved()

				return hs
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
				if !contains(contentStr, "post-clone-overlay") {
					t.Error("Expected hook name 'post-clone-overlay' in saved content")
				}
				if !contains(contentStr, "post-fork-script") {
					t.Error("Expected hook name 'post-fork-script' in saved content")
				}
				if !contains(contentStr, "github.com/owner/*") {
					t.Error("Expected repo pattern in saved content")
				}
			},
		},
		{
			name: "Save with force",
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				// Even with force, HasChanges is still called (different from extra_store)
				hs.EXPECT().HasChanges().Return(false)

				testHook := hook.ConcreteHook(
					uuid.New(),
					"force-save-hook",
					"*",
					string(hook.EventPostClone),
					string(hook.OperationTypeOverlay),
					"test-id",
				)

				hs.EXPECT().List().Return(func(yield func(hook.Hook, error) bool) {
					yield(testHook, nil)
				})

				hs.EXPECT().MarkSaved()

				return hs
			},
			force:   true,
			wantErr: false,
			validate: func(t *testing.T, path string) {
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("Failed to read saved file: %v", err)
				}

				if !contains(string(content), "force-save-hook") {
					t.Error("Expected force-saved hook in content")
				}
			},
		},
		{
			name: "Save with list error",
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().HasChanges().Return(true)

				hs.EXPECT().List().Return(func(yield func(hook.Hook, error) bool) {
					yield(nil, os.ErrPermission)
				})

				return hs
			},
			force:   false,
			wantErr: true,
		},
		{
			name: "Save with directory creation failure",
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().HasChanges().Return(true)

				testHook := hook.ConcreteHook(
					uuid.New(),
					"test-hook",
					"*",
					string(hook.EventPostClone),
					string(hook.OperationTypeOverlay),
					"test-id",
				)

				hs.EXPECT().List().Return(func(yield func(hook.Hook, error) bool) {
					yield(testHook, nil)
				})

				hs.EXPECT().MarkSaved()

				return hs
			},
			force:   false,
			wantErr: false, // Directory creation will succeed in temp dir
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tempDir := t.TempDir()
			path := filepath.Join(tempDir, "hook.toml")
			t.Setenv("GOGH_HOOK_PATH", path)

			store := config.NewHookStore()
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

func TestHookDir(t *testing.T) {
	// Test with no environment variable
	dir, err := config.HookDir()
	if err != nil {
		t.Fatalf("HookDir() error = %v", err)
	}

	if !contains(dir, "hook.v4.toml") {
		t.Errorf("Expected dir to contain 'hook.v4.toml', got %s", dir)
	}

	// Test with environment variable
	customPath := "/custom/path/hook.toml"
	t.Setenv("GOGH_HOOK_PATH", customPath)

	dir, err = config.HookDir()
	if err != nil {
		t.Fatalf("HookDir() with env var error = %v", err)
	}

	if dir != customPath {
		t.Errorf("Expected dir to be %s, got %s", customPath, dir)
	}
}
