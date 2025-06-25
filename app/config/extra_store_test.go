package config_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/config"
	"github.com/kyoh86/gogh/v4/core/extra"
	"github.com/kyoh86/gogh/v4/core/extra_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"go.uber.org/mock/gomock"
)

func TestExtraStore_Source(t *testing.T) {
	store := config.NewExtraStore()

	// Test with no environment variable
	source, err := store.Source()
	if err != nil {
		t.Fatalf("Source() error = %v", err)
	}

	// Should contain "extra.v4.toml"
	if !contains(source, "extra.v4.toml") {
		t.Errorf("Expected source to contain 'extra.v4.toml', got %s", source)
	}

	// Test with environment variable
	tempDir := t.TempDir()
	customPath := filepath.Join(tempDir, "custom-extra.toml")
	t.Setenv("GOGH_EXTRA_PATH", customPath)

	source, err = store.Source()
	if err != nil {
		t.Fatalf("Source() with env var error = %v", err)
	}

	if source != customPath {
		t.Errorf("Expected source to be %s, got %s", customPath, source)
	}
}

func TestExtraStore_Load(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		setupFile func(t *testing.T) string
		setupMock func(*gomock.Controller) *extra_mock.MockExtraService
		wantErr   bool
		wantEmpty bool
	}{
		{
			name: "Load from non-existent file",
			setupFile: func(t *testing.T) string {
				tempDir := t.TempDir()
				return filepath.Join(tempDir, "non-existent.toml")
			},
			setupMock: func(ctrl *gomock.Controller) *extra_mock.MockExtraService {
				es := extra_mock.NewMockExtraService(ctrl)
				es.EXPECT().MarkSaved()
				return es
			},
			wantErr:   false,
			wantEmpty: true,
		},
		{
			name: "Load valid extras",
			setupFile: func(t *testing.T) string {
				tempDir := t.TempDir()
				path := filepath.Join(tempDir, "extra.toml")

				content := `[[extra]]
id = "auto-extra-id"
type = "auto"
repository = "github.com/owner/repo"
source = "github.com/source/repo"
created_at = 2023-01-01T00:00:00Z

[[extra.items]]
overlay_id = "overlay-1"
hook_id = "hook-1"

[[extra]]
id = "named-extra-id"
type = "named"
name = "my-config"
source = "github.com/config/source"
created_at = 2023-01-02T00:00:00Z

[[extra.items]]
overlay_id = "overlay-2"
hook_id = "hook-2"
`
				if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				return path
			},
			setupMock: func(ctrl *gomock.Controller) *extra_mock.MockExtraService {
				es := extra_mock.NewMockExtraService(ctrl)

				// Expect Load to be called with proper extras
				es.EXPECT().Load(gomock.Any()).DoAndReturn(func(yield func(func(*extra.Extra, error) bool)) error {
					loadedCount := 0
					yield(func(e *extra.Extra, err error) bool {
						if err != nil {
							t.Errorf("Unexpected error in yield: %v", err)
							return false
						}
						loadedCount++

						// Verify loaded extras
						switch loadedCount {
						case 1:
							if e.ID() != "auto-extra-id" {
								t.Errorf("Expected ID 'auto-extra-id', got %s", e.ID())
							}
							if e.Type() != extra.TypeAuto {
								t.Errorf("Expected type auto, got %s", e.Type())
							}
							if e.Repository() == nil || e.Repository().String() != "github.com/owner/repo" {
								t.Error("Expected repository github.com/owner/repo")
							}
						case 2:
							if e.ID() != "named-extra-id" {
								t.Errorf("Expected ID 'named-extra-id', got %s", e.ID())
							}
							if e.Type() != extra.TypeNamed {
								t.Errorf("Expected type named, got %s", e.Type())
							}
							if e.Name() != "my-config" {
								t.Errorf("Expected name 'my-config', got %s", e.Name())
							}
						}
						return true
					})

					if loadedCount != 2 {
						t.Errorf("Expected 2 extras to be loaded, got %d", loadedCount)
					}
					return nil
				})

				es.EXPECT().MarkSaved()

				return es
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
			setupMock: extra_mock.NewMockExtraService,
			wantErr:   true,
		},
		{
			name: "Load with invalid reference",
			setupFile: func(t *testing.T) string {
				tempDir := t.TempDir()
				path := filepath.Join(tempDir, "extra.toml")

				content := `[[extra]]
id = "invalid-ref-id"
type = "auto"
repository = "invalid-reference"
source = "github.com/source/repo"
created_at = 2023-01-01T00:00:00Z
`
				if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				return path
			},
			setupMock: extra_mock.NewMockExtraService,
			wantErr:   true,
		},
		{
			name: "Load with unknown type",
			setupFile: func(t *testing.T) string {
				tempDir := t.TempDir()
				path := filepath.Join(tempDir, "extra.toml")

				content := `[[extra]]
id = "unknown-type-id"
type = "unknown"
source = "github.com/source/repo"
created_at = 2023-01-01T00:00:00Z
`
				if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				return path
			},
			setupMock: extra_mock.NewMockExtraService,
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			path := tc.setupFile(t)
			t.Setenv("GOGH_EXTRA_PATH", path)

			store := config.NewExtraStore()
			mockService := tc.setupMock(ctrl)

			initial := func() extra.ExtraService {
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

func TestExtraStore_Save(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		setupMock func(*gomock.Controller) *extra_mock.MockExtraService
		force     bool
		wantErr   bool
		validate  func(t *testing.T, path string)
	}{
		{
			name: "Save with no changes and no force",
			setupMock: func(ctrl *gomock.Controller) *extra_mock.MockExtraService {
				es := extra_mock.NewMockExtraService(ctrl)
				es.EXPECT().HasChanges().Return(false)
				return es
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
			setupMock: func(ctrl *gomock.Controller) *extra_mock.MockExtraService {
				es := extra_mock.NewMockExtraService(ctrl)
				es.EXPECT().HasChanges().Return(true)

				// Create test extras
				autoRef := repository.NewReference("github.com", "owner", "repo")
				sourceRef := repository.NewReference("github.com", "source", "repo")
				autoExtra := extra.NewAutoExtra(
					uuid.New().String(),
					autoRef,
					sourceRef,
					[]extra.Item{
						{OverlayID: "overlay-1", HookID: "hook-1"},
					},
					time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				)

				namedExtra := extra.NewNamedExtra(
					uuid.New().String(),
					"my-config",
					sourceRef,
					[]extra.Item{
						{OverlayID: "overlay-2", HookID: "hook-2"},
					},
					time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				)

				es.EXPECT().List(ctx).Return(func(yield func(*extra.Extra, error) bool) {
					yield(autoExtra, nil)
					yield(namedExtra, nil)
				})

				es.EXPECT().MarkSaved()

				return es
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
				t.Logf("Saved content:\n%s", contentStr)
				if !contains(contentStr, "type = \"auto\"") && !contains(contentStr, "type = 'auto'") {
					t.Error("Expected auto extra in saved content")
				}
				if !contains(contentStr, "type = \"named\"") && !contains(contentStr, "type = 'named'") {
					t.Error("Expected named extra in saved content")
				}
				if !contains(contentStr, "name = \"my-config\"") && !contains(contentStr, "name = 'my-config'") {
					t.Error("Expected named extra name in saved content")
				}
			},
		},
		{
			name: "Save with force",
			setupMock: func(ctrl *gomock.Controller) *extra_mock.MockExtraService {
				es := extra_mock.NewMockExtraService(ctrl)
				// HasChanges should not be called when force is true

				sourceRef := repository.NewReference("github.com", "source", "repo")
				namedExtra := extra.NewNamedExtra(
					uuid.New().String(),
					"force-save",
					sourceRef,
					[]extra.Item{},
					time.Now(),
				)

				es.EXPECT().List(ctx).Return(func(yield func(*extra.Extra, error) bool) {
					yield(namedExtra, nil)
				})

				es.EXPECT().MarkSaved()

				return es
			},
			force:   true,
			wantErr: false,
			validate: func(t *testing.T, path string) {
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("Failed to read saved file: %v", err)
				}

				contentStr := string(content)
				t.Logf("Force saved content:\n%s", contentStr)
				if !contains(contentStr, "name = \"force-save\"") && !contains(contentStr, "name = 'force-save'") {
					t.Error("Expected force-saved extra in content")
				}
			},
		},
		{
			name: "Save with list error",
			setupMock: func(ctrl *gomock.Controller) *extra_mock.MockExtraService {
				es := extra_mock.NewMockExtraService(ctrl)
				es.EXPECT().HasChanges().Return(true)

				es.EXPECT().List(ctx).Return(func(yield func(*extra.Extra, error) bool) {
					yield(nil, os.ErrPermission)
				})

				return es
			},
			force:   false,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tempDir := t.TempDir()
			path := filepath.Join(tempDir, "extra.toml")
			t.Setenv("GOGH_EXTRA_PATH", path)

			store := config.NewExtraStore()
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

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
