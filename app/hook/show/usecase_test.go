package show_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/hook/show"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/hook_mock"
	"go.uber.org/mock/gomock"
)

func TestUsecase_Execute(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		hookID    string
		asJSON    bool
		setupMock func(*gomock.Controller) *hook_mock.MockHookService
		wantErr   bool
		validate  func(*testing.T, string)
	}{
		{
			name:   "Show overlay hook as one-line",
			hookID: uuid.New().String(),
			asJSON: false,
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				h := hook.ConcreteHook(
					uuid.New(),
					"test-hook",
					"github.com/owner/*",
					string(hook.EventPostClone),
					string(hook.OperationTypeOverlay),
					uuid.New(),
				)
				hs.EXPECT().Get(ctx, gomock.Any()).Return(h, nil)
				return hs
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				// One-line format should contain ID prefix, name, event, etc.
				if !strings.Contains(output, "test-hook") {
					t.Error("Expected output to contain hook name")
				}
				if !strings.Contains(output, "post-clone") {
					t.Error("Expected output to contain trigger event")
				}
				if !strings.Contains(output, "overlay") {
					t.Error("Expected output to contain operation type")
				}
			},
		},
		{
			name:   "Show script hook as JSON",
			hookID: uuid.New().String(),
			asJSON: true,
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				operationID := uuid.New()
				h := hook.ConcreteHook(
					uuid.New(),
					"script-hook",
					"github.com/test/*",
					string(hook.EventPostFork),
					string(hook.OperationTypeScript),
					operationID,
				)
				hs.EXPECT().Get(ctx, gomock.Any()).Return(h, nil)
				return hs
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				// Should be valid JSON
				var data map[string]any
				if err := json.Unmarshal([]byte(output), &data); err != nil {
					t.Errorf("Expected valid JSON output, got error: %v", err)
				}
				// Check fields
				if data["name"] != "script-hook" {
					t.Errorf("Expected name 'script-hook', got %v", data["name"])
				}
				if data["trigger_event"] != "post-fork" {
					t.Errorf("Expected trigger_event 'post-fork', got %v", data["trigger_event"])
				}
				if data["operation_type"] != "script" {
					t.Errorf("Expected operation_type 'script', got %v", data["operation_type"])
				}
				// Just check that operation_id is a valid UUID
				operationIDStr, ok := data["operation_id"].(string)
				if !ok {
					t.Errorf("Expected operation_id to be a string, got %T", data["operation_id"])
				}
				if _, err := uuid.Parse(operationIDStr); err != nil {
					t.Errorf("Expected operation_id to be a valid UUID, got %v", operationIDStr)
				}
			},
		},
		{
			name:   "Show hook with empty pattern",
			hookID: uuid.New().String(),
			asJSON: false,
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				h := hook.ConcreteHook(
					uuid.New(),
					"global-hook",
					"", // Empty pattern means global
					string(hook.EventPostCreate),
					string(hook.OperationTypeOverlay),
					uuid.New(),
				)
				hs.EXPECT().Get(ctx, gomock.Any()).Return(h, nil)
				return hs
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "global-hook") {
					t.Error("Expected output to contain hook name")
				}
				if !strings.Contains(output, "post-create") {
					t.Error("Expected output to contain trigger event")
				}
			},
		},
		{
			name:   "Show hook with 'any' event",
			hookID: uuid.New().String(),
			asJSON: true,
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				h := hook.ConcreteHook(
					uuid.New(),
					"any-event-hook",
					"github.com/any/*",
					string(hook.EventAny),
					string(hook.OperationTypeScript),
					uuid.New(),
				)
				hs.EXPECT().Get(ctx, gomock.Any()).Return(h, nil)
				return hs
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				var data map[string]any
				if err := json.Unmarshal([]byte(output), &data); err != nil {
					t.Errorf("Expected valid JSON output, got error: %v", err)
				}
				// EventAny is represented as empty string
				if data["trigger_event"] != "" {
					t.Errorf("Expected trigger_event '' (for 'any'), got %v", data["trigger_event"])
				}
			},
		},
		{
			name:   "Hook not found",
			hookID: uuid.New().String(),
			asJSON: false,
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Get(ctx, gomock.Any()).Return(nil, errors.New("hook not found"))
				return hs
			},
			wantErr: true,
		},
		{
			name:   "Invalid hook ID",
			hookID: "invalid-id",
			asJSON: false,
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Get(ctx, "invalid-id").Return(nil, errors.New("invalid hook ID"))
				return hs
			},
			wantErr: true,
		},
		{
			name:   "Empty hook ID",
			hookID: "",
			asJSON: false,
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Get(ctx, "").Return(nil, errors.New("hook ID is required"))
				return hs
			},
			wantErr: true,
		},
		{
			name:   "Service returns unexpected error",
			hookID: uuid.New().String(),
			asJSON: true,
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Get(ctx, gomock.Any()).Return(nil, errors.New("database error"))
				return hs
			},
			wantErr: true,
		},
		{
			name:   "Show hook with long name",
			hookID: uuid.New().String(),
			asJSON: false,
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				longName := "this-is-a-very-long-hook-name-that-might-be-truncated-in-display"
				h := hook.ConcreteHook(
					uuid.New(),
					longName,
					"github.com/long/*",
					string(hook.EventPostClone),
					string(hook.OperationTypeOverlay),
					uuid.New(),
				)
				hs.EXPECT().Get(ctx, gomock.Any()).Return(h, nil)
				return hs
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				// Should contain at least part of the name
				if !strings.Contains(output, "this-is-a-very-long-hook-name") {
					t.Error("Expected output to contain at least part of the long name")
				}
			},
		},
		{
			name:   "Show hook with complex pattern",
			hookID: uuid.New().String(),
			asJSON: true,
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				h := hook.ConcreteHook(
					uuid.New(),
					"complex-pattern-hook",
					"github.com/{owner1,owner2}/{repo1,repo2,repo3}",
					string(hook.EventPostFork),
					string(hook.OperationTypeScript),
					uuid.New(),
				)
				hs.EXPECT().Get(ctx, gomock.Any()).Return(h, nil)
				return hs
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				var data map[string]any
				if err := json.Unmarshal([]byte(output), &data); err != nil {
					t.Errorf("Expected valid JSON output, got error: %v", err)
				}
				pattern := data["repo_pattern"].(string)
				if !strings.Contains(pattern, "{owner1,owner2}") {
					t.Errorf("Expected pattern to contain '{owner1,owner2}', got %v", pattern)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var buf bytes.Buffer
			hs := tc.setupMock(ctrl)
			uc := testtarget.NewUsecase(hs, &buf)

			err := uc.Execute(ctx, tc.hookID, tc.asJSON)
			if (err != nil) != tc.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && tc.validate != nil {
				tc.validate(t, buf.String())
			}
		})
	}
}

func TestUsecase_Execute_AllEventTypes(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	events := []hook.Event{
		hook.EventPostClone,
		hook.EventPostFork,
		hook.EventPostCreate,
		hook.EventAny,
	}

	for _, event := range events {
		t.Run(string(event), func(t *testing.T) {
			var buf bytes.Buffer
			hs := hook_mock.NewMockHookService(ctrl)
			uc := testtarget.NewUsecase(hs, &buf)

			h := hook.ConcreteHook(
				uuid.New(),
				"test-hook",
				"github.com/test/*",
				string(event),
				string(hook.OperationTypeOverlay),
				uuid.New(),
			)
			hs.EXPECT().Get(ctx, gomock.Any()).Return(h, nil)

			err := uc.Execute(ctx, uuid.New().String(), true)
			if err != nil {
				t.Errorf("Execute() unexpected error for event %s: %v", event, err)
			}

			// Verify JSON contains correct event
			var data map[string]any
			if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
				t.Errorf("Failed to parse JSON for event %s: %v", event, err)
			}
			if data["trigger_event"] != string(event) {
				t.Errorf("Expected trigger_event %s, got %v", event, data["trigger_event"])
			}
		})
	}
}

func TestUsecase_Execute_AllOperationTypes(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	operations := []hook.OperationType{
		hook.OperationTypeOverlay,
		hook.OperationTypeScript,
	}

	for _, opType := range operations {
		t.Run(string(opType), func(t *testing.T) {
			var buf bytes.Buffer
			hs := hook_mock.NewMockHookService(ctrl)
			uc := testtarget.NewUsecase(hs, &buf)

			h := hook.ConcreteHook(
				uuid.New(),
				"test-hook",
				"github.com/test/*",
				string(hook.EventPostClone),
				string(opType),
				uuid.New(),
			)
			hs.EXPECT().Get(ctx, gomock.Any()).Return(h, nil)

			err := uc.Execute(ctx, uuid.New().String(), true)
			if err != nil {
				t.Errorf("Execute() unexpected error for operation %s: %v", opType, err)
			}

			// Verify JSON contains correct operation type
			var data map[string]any
			if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
				t.Errorf("Failed to parse JSON for operation %s: %v", opType, err)
			}
			if data["operation_type"] != string(opType) {
				t.Errorf("Expected operation_type %s, got %v", opType, data["operation_type"])
			}
		})
	}
}
