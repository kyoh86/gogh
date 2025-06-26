package update_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/hook/update"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/hook_mock"
	"go.uber.org/mock/gomock"
)

func TestUsecase_Execute(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		hookID    string
		opts      testtarget.Options
		setupMock func(*gomock.Controller) *hook_mock.MockHookService
		wantErr   bool
	}{
		{
			name:   "Successfully update hook with all fields",
			hookID: uuid.New().String(),
			opts: testtarget.Options{
				Name:          "updated-hook",
				RepoPattern:   "github.com/updated/*",
				TriggerEvent:  string(hook.EventPostClone),
				OperationType: string(hook.OperationTypeOverlay),
				OperationID:   "overlay-updated",
			},
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry hook.Entry) error {
						// Validate the ID
						if _, err := uuid.Parse(id); err != nil {
							t.Errorf("Expected valid UUID, got %s", id)
						}
						// Validate the entry
						if entry.Name != "updated-hook" {
							t.Errorf("Expected name 'updated-hook', got %s", entry.Name)
						}
						if entry.RepoPattern != "github.com/updated/*" {
							t.Errorf("Expected repo pattern 'github.com/updated/*', got %s", entry.RepoPattern)
						}
						if entry.TriggerEvent != hook.EventPostClone {
							t.Errorf("Expected trigger event %s, got %s", hook.EventPostClone, entry.TriggerEvent)
						}
						if entry.OperationType != hook.OperationTypeOverlay {
							t.Errorf("Expected operation type %s, got %s", hook.OperationTypeOverlay, entry.OperationType)
						}
						if entry.OperationID != "overlay-updated" {
							t.Errorf("Expected operation ID 'overlay-updated', got %s", entry.OperationID)
						}
						return nil
					},
				)
				return hs
			},
			wantErr: false,
		},
		{
			name:   "Update hook with partial fields",
			hookID: uuid.New().String(),
			opts: testtarget.Options{
				Name:          "partial-update",
				RepoPattern:   "",
				TriggerEvent:  string(hook.EventPostFork),
				OperationType: "",
				OperationID:   "",
			},
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry hook.Entry) error {
						if entry.Name != "partial-update" {
							t.Errorf("Expected name 'partial-update', got %s", entry.Name)
						}
						if entry.RepoPattern != "" {
							t.Errorf("Expected empty repo pattern, got %s", entry.RepoPattern)
						}
						if entry.TriggerEvent != hook.EventPostFork {
							t.Errorf("Expected trigger event %s, got %s", hook.EventPostFork, entry.TriggerEvent)
						}
						if entry.OperationType != "" {
							t.Errorf("Expected empty operation type, got %s", entry.OperationType)
						}
						if entry.OperationID != "" {
							t.Errorf("Expected empty operation ID, got %s", entry.OperationID)
						}
						return nil
					},
				)
				return hs
			},
			wantErr: false,
		},
		{
			name:   "Update hook changing operation type",
			hookID: uuid.New().String(),
			opts: testtarget.Options{
				Name:          "script-hook",
				RepoPattern:   "github.com/test/*",
				TriggerEvent:  string(hook.EventPostCreate),
				OperationType: string(hook.OperationTypeScript),
				OperationID:   "script-123",
			},
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry hook.Entry) error {
						if entry.OperationType != hook.OperationTypeScript {
							t.Errorf("Expected operation type %s, got %s", hook.OperationTypeScript, entry.OperationType)
						}
						if entry.OperationID != "script-123" {
							t.Errorf("Expected operation ID 'script-123', got %s", entry.OperationID)
						}
						return nil
					},
				)
				return hs
			},
			wantErr: false,
		},
		{
			name:   "Update non-existent hook",
			hookID: uuid.New().String(),
			opts: testtarget.Options{
				Name:          "non-existent",
				RepoPattern:   "github.com/test/*",
				TriggerEvent:  string(hook.EventPostClone),
				OperationType: string(hook.OperationTypeOverlay),
				OperationID:   "overlay-123",
			},
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(errors.New("hook not found"))
				return hs
			},
			wantErr: true,
		},
		{
			name:   "Update with invalid hook ID",
			hookID: "invalid-id",
			opts: testtarget.Options{
				Name:          "test",
				RepoPattern:   "github.com/test/*",
				TriggerEvent:  string(hook.EventPostClone),
				OperationType: string(hook.OperationTypeOverlay),
				OperationID:   "overlay-123",
			},
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Update(ctx, "invalid-id", gomock.Any()).Return(errors.New("invalid hook ID"))
				return hs
			},
			wantErr: true,
		},
		{
			name:   "Update with empty hook ID",
			hookID: "",
			opts: testtarget.Options{
				Name:          "test",
				RepoPattern:   "github.com/test/*",
				TriggerEvent:  string(hook.EventPostClone),
				OperationType: string(hook.OperationTypeOverlay),
				OperationID:   "overlay-123",
			},
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Update(ctx, "", gomock.Any()).Return(errors.New("hook ID is required"))
				return hs
			},
			wantErr: true,
		},
		{
			name:   "Update with all empty values",
			hookID: uuid.New().String(),
			opts: testtarget.Options{
				Name:          "",
				RepoPattern:   "",
				TriggerEvent:  "",
				OperationType: "",
				OperationID:   "",
			},
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry hook.Entry) error {
						// Service might validate and return error
						return errors.New("at least one field must be provided")
					},
				)
				return hs
			},
			wantErr: true,
		},
		{
			name:   "Update with custom event string",
			hookID: uuid.New().String(),
			opts: testtarget.Options{
				Name:          "custom-event-hook",
				RepoPattern:   "github.com/custom/*",
				TriggerEvent:  "custom-event",
				OperationType: string(hook.OperationTypeScript),
				OperationID:   "script-custom",
			},
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string, entry hook.Entry) error {
						if string(entry.TriggerEvent) != "custom-event" {
							t.Errorf("Expected trigger event 'custom-event', got %s", entry.TriggerEvent)
						}
						return nil
					},
				)
				return hs
			},
			wantErr: false,
		},
		{
			name:   "Service returns unexpected error",
			hookID: uuid.New().String(),
			opts: testtarget.Options{
				Name:          "test",
				RepoPattern:   "github.com/test/*",
				TriggerEvent:  string(hook.EventPostClone),
				OperationType: string(hook.OperationTypeOverlay),
				OperationID:   "overlay-123",
			},
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).Return(errors.New("storage error"))
				return hs
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			hs := tc.setupMock(ctrl)
			uc := testtarget.NewUsecase(hs)

			err := uc.Execute(ctx, tc.hookID, tc.opts)
			if (err != nil) != tc.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tc.wantErr)
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

			hs := hook_mock.NewMockHookService(ctrl)
			hs.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, id string, entry hook.Entry) error {
					if entry.TriggerEvent != event {
						t.Errorf("Expected trigger event %s, got %s", event, entry.TriggerEvent)
					}
					return nil
				},
			)

			uc := testtarget.NewUsecase(hs)
			opts := testtarget.Options{
				Name:          "test-hook",
				RepoPattern:   "github.com/test/*",
				TriggerEvent:  string(event),
				OperationType: string(hook.OperationTypeOverlay),
				OperationID:   "overlay-123",
			}

			err := uc.Execute(ctx, uuid.New().String(), opts)
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
			hs.EXPECT().Update(ctx, gomock.Any(), gomock.Any()).DoAndReturn(
				func(ctx context.Context, id string, entry hook.Entry) error {
					if entry.OperationType != opType {
						t.Errorf("Expected operation type %s, got %s", opType, entry.OperationType)
					}
					return nil
				},
			)

			uc := testtarget.NewUsecase(hs)
			opts := testtarget.Options{
				Name:          "test-hook",
				RepoPattern:   "github.com/test/*",
				TriggerEvent:  string(hook.EventPostClone),
				OperationType: string(opType),
				OperationID:   "operation-123",
			}

			err := uc.Execute(ctx, uuid.New().String(), opts)
			if err != nil {
				t.Errorf("Execute() unexpected error = %v", err)
			}
		})
	}
}
