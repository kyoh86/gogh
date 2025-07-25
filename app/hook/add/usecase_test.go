package add_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/hook/add"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/hook_mock"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"github.com/kyoh86/gogh/v4/core/script_mock"
	"go.uber.org/mock/gomock"
)

func TestUsecase_Execute(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name       string
		opts       testtarget.Options
		setupMock  func(*gomock.Controller) (*hook_mock.MockHookService, *overlay_mock.MockOverlayService, *script_mock.MockScriptService)
		wantErr    bool
		validateID func(string) error
	}{
		{
			name: "Successfully add overlay hook",
			opts: testtarget.Options{
				Name:          "test-hook",
				RepoPattern:   "github.com/owner/*",
				TriggerEvent:  string(hook.EventPostClone),
				OperationType: string(hook.OperationTypeOverlay),
				OperationID:   "overlay-123",
			},
			setupMock: func(ctrl *gomock.Controller) (*hook_mock.MockHookService, *overlay_mock.MockOverlayService, *script_mock.MockScriptService) {
				overlayID := uuid.New()
				hs := hook_mock.NewMockHookService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				ss := script_mock.NewMockScriptService(ctrl)

				// Mock overlay resolution
				mockOverlay := overlay_mock.NewMockOverlay(ctrl)
				mockOverlay.EXPECT().UUID().Return(overlayID)
				os.EXPECT().Get(ctx, "overlay-123").Return(mockOverlay, nil)

				hs.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry hook.Entry) (string, error) {
						// Validate the entry
						if entry.Name != "test-hook" {
							t.Errorf("Expected name 'test-hook', got %s", entry.Name)
						}
						if entry.RepoPattern != "github.com/owner/*" {
							t.Errorf("Expected repo pattern 'github.com/owner/*', got %s", entry.RepoPattern)
						}
						if entry.TriggerEvent != hook.EventPostClone {
							t.Errorf("Expected trigger event %s, got %s", hook.EventPostClone, entry.TriggerEvent)
						}
						if entry.OperationType != hook.OperationTypeOverlay {
							t.Errorf("Expected operation type %s, got %s", hook.OperationTypeOverlay, entry.OperationType)
						}
						if entry.OperationID != overlayID {
							t.Errorf("Expected operation ID %s, got %s", overlayID, entry.OperationID)
						}
						return uuid.New().String(), nil
					},
				)
				return hs, os, ss
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
			name: "Successfully add script hook",
			opts: testtarget.Options{
				Name:          "script-hook",
				RepoPattern:   "github.com/test/*",
				TriggerEvent:  string(hook.EventPostFork),
				OperationType: string(hook.OperationTypeScript),
				OperationID:   "script-456",
			},
			setupMock: func(ctrl *gomock.Controller) (*hook_mock.MockHookService, *overlay_mock.MockOverlayService, *script_mock.MockScriptService) {
				scriptID := uuid.New()
				hs := hook_mock.NewMockHookService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				ss := script_mock.NewMockScriptService(ctrl)

				// Mock script resolution
				mockScript := script_mock.NewMockScript(ctrl)
				mockScript.EXPECT().UUID().Return(scriptID)
				ss.EXPECT().Get(ctx, "script-456").Return(mockScript, nil)

				hs.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry hook.Entry) (string, error) {
						if entry.Name != "script-hook" {
							t.Errorf("Expected name 'script-hook', got %s", entry.Name)
						}
						if entry.TriggerEvent != hook.EventPostFork {
							t.Errorf("Expected trigger event %s, got %s", hook.EventPostFork, entry.TriggerEvent)
						}
						if entry.OperationType != hook.OperationTypeScript {
							t.Errorf("Expected operation type %s, got %s", hook.OperationTypeScript, entry.OperationType)
						}
						if entry.OperationID != scriptID {
							t.Errorf("Expected operation ID %s, got %s", scriptID, entry.OperationID)
						}
						return uuid.New().String(), nil
					},
				)
				return hs, os, ss
			},
			wantErr: false,
		},
		{
			name: "Add hook with empty pattern (global hook)",
			opts: testtarget.Options{
				Name:          "global-hook",
				RepoPattern:   "",
				TriggerEvent:  string(hook.EventPostCreate),
				OperationType: string(hook.OperationTypeOverlay),
				OperationID:   "overlay-789",
			},
			setupMock: func(ctrl *gomock.Controller) (*hook_mock.MockHookService, *overlay_mock.MockOverlayService, *script_mock.MockScriptService) {
				overlayID := uuid.New()
				hs := hook_mock.NewMockHookService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				ss := script_mock.NewMockScriptService(ctrl)

				// Mock overlay resolution
				mockOverlay := overlay_mock.NewMockOverlay(ctrl)
				mockOverlay.EXPECT().UUID().Return(overlayID)
				os.EXPECT().Get(ctx, "overlay-789").Return(mockOverlay, nil)

				hs.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry hook.Entry) (string, error) {
						if entry.RepoPattern != "" {
							t.Errorf("Expected empty repo pattern, got %s", entry.RepoPattern)
						}
						return uuid.New().String(), nil
					},
				)
				return hs, os, ss
			},
			wantErr: false,
		},
		{
			name: "Hook service returns error",
			opts: testtarget.Options{
				Name:          "error-hook",
				RepoPattern:   "github.com/error/*",
				TriggerEvent:  string(hook.EventPostClone),
				OperationType: string(hook.OperationTypeOverlay),
				OperationID:   "overlay-error",
			},
			setupMock: func(ctrl *gomock.Controller) (*hook_mock.MockHookService, *overlay_mock.MockOverlayService, *script_mock.MockScriptService) {
				hs := hook_mock.NewMockHookService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				ss := script_mock.NewMockScriptService(ctrl)

				// Mock overlay resolution failure
				os.EXPECT().Get(ctx, "overlay-error").Return(nil, errors.New("overlay not found"))

				return hs, os, ss
			},
			wantErr: true,
		},
		{
			name: "Add hook with custom event string",
			opts: testtarget.Options{
				Name:          "custom-event-hook",
				RepoPattern:   "github.com/custom/*",
				TriggerEvent:  "custom-event", // Not a predefined event
				OperationType: string(hook.OperationTypeScript),
				OperationID:   "script-custom",
			},
			setupMock: func(ctrl *gomock.Controller) (*hook_mock.MockHookService, *overlay_mock.MockOverlayService, *script_mock.MockScriptService) {
				scriptID := uuid.New()
				hs := hook_mock.NewMockHookService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				ss := script_mock.NewMockScriptService(ctrl)

				// Mock script resolution
				mockScript := script_mock.NewMockScript(ctrl)
				mockScript.EXPECT().UUID().Return(scriptID)
				ss.EXPECT().Get(ctx, "script-custom").Return(mockScript, nil)

				hs.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry hook.Entry) (string, error) {
						if string(entry.TriggerEvent) != "custom-event" {
							t.Errorf("Expected trigger event 'custom-event', got %s", entry.TriggerEvent)
						}
						return uuid.New().String(), nil
					},
				)
				return hs, os, ss
			},
			wantErr: false,
		},
		{
			name: "Add hook with all empty values",
			opts: testtarget.Options{
				Name:          "",
				RepoPattern:   "",
				TriggerEvent:  "",
				OperationType: "",
				OperationID:   "",
			},
			setupMock: func(ctrl *gomock.Controller) (*hook_mock.MockHookService, *overlay_mock.MockOverlayService, *script_mock.MockScriptService) {
				hs := hook_mock.NewMockHookService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				ss := script_mock.NewMockScriptService(ctrl)
				// Hook service should be called even with empty operation type
				hs.EXPECT().Add(ctx, gomock.Any()).Return("", errors.New("invalid hook configuration"))
				return hs, os, ss
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			hs, os, ss := tc.setupMock(ctrl)
			uc := testtarget.NewUsecase(hs, os, ss)

			id, err := uc.Execute(ctx, tc.opts)
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

func TestUsecase_Execute_AllEventTypes(t *testing.T) {
	ctx := context.Background()

	events := []hook.Event{
		hook.EventPostClone,
		hook.EventPostFork,
		hook.EventPostCreate,
		hook.EventAny,
	}

	for _, event := range events {
		t.Run(string(event), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			overlayID := uuid.New()
			hs := hook_mock.NewMockHookService(ctrl)
			os := overlay_mock.NewMockOverlayService(ctrl)
			ss := script_mock.NewMockScriptService(ctrl)

			// Mock overlay resolution
			mockOverlay := overlay_mock.NewMockOverlay(ctrl)
			mockOverlay.EXPECT().UUID().Return(overlayID)
			os.EXPECT().Get(ctx, "overlay-123").Return(mockOverlay, nil)

			hs.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
				func(ctx context.Context, entry hook.Entry) (string, error) {
					if entry.TriggerEvent != event {
						t.Errorf("Expected trigger event %s, got %s", event, entry.TriggerEvent)
					}
					return uuid.New().String(), nil
				},
			)

			uc := testtarget.NewUsecase(hs, os, ss)
			opts := testtarget.Options{
				Name:          "test-hook",
				RepoPattern:   "github.com/test/*",
				TriggerEvent:  string(event),
				OperationType: string(hook.OperationTypeOverlay),
				OperationID:   "overlay-123",
			}

			_, err := uc.Execute(ctx, opts)
			if err != nil {
				t.Errorf("Execute() unexpected error = %v", err)
			}
		})
	}
}

func TestUsecase_Execute_AllOperationTypes(t *testing.T) {
	ctx := context.Background()

	operations := []hook.OperationType{
		hook.OperationTypeOverlay,
		hook.OperationTypeScript,
	}

	for _, opType := range operations {
		t.Run(string(opType), func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			hs := hook_mock.NewMockHookService(ctrl)
			os := overlay_mock.NewMockOverlayService(ctrl)
			ss := script_mock.NewMockScriptService(ctrl)

			// Mock appropriate service based on operation type
			if opType == hook.OperationTypeOverlay {
				overlayID := uuid.New()
				mockOverlay := overlay_mock.NewMockOverlay(ctrl)
				mockOverlay.EXPECT().UUID().Return(overlayID)
				os.EXPECT().Get(ctx, "operation-123").Return(mockOverlay, nil)
				hs.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry hook.Entry) (string, error) {
						if entry.OperationType != opType {
							t.Errorf("Expected operation type %s, got %s", opType, entry.OperationType)
						}
						if entry.OperationID != overlayID {
							t.Errorf("Expected operation ID %s, got %s", overlayID, entry.OperationID)
						}
						return uuid.New().String(), nil
					},
				)
			} else {
				scriptID := uuid.New()
				mockScript := script_mock.NewMockScript(ctrl)
				mockScript.EXPECT().UUID().Return(scriptID)
				ss.EXPECT().Get(ctx, "operation-123").Return(mockScript, nil)
				hs.EXPECT().Add(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, entry hook.Entry) (string, error) {
						if entry.OperationType != opType {
							t.Errorf("Expected operation type %s, got %s", opType, entry.OperationType)
						}
						if entry.OperationID != scriptID {
							t.Errorf("Expected operation ID %s, got %s", scriptID, entry.OperationID)
						}
						return uuid.New().String(), nil
					},
				)
			}

			uc := testtarget.NewUsecase(hs, os, ss)
			opts := testtarget.Options{
				Name:          "test-hook",
				RepoPattern:   "github.com/test/*",
				TriggerEvent:  string(hook.EventPostClone),
				OperationType: string(opType),
				OperationID:   "operation-123",
			}

			_, err := uc.Execute(ctx, opts)
			if err != nil {
				t.Errorf("Execute() unexpected error = %v", err)
			}
		})
	}
}
