package hook_remove_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/hook_remove"
	"github.com/kyoh86/gogh/v4/core/hook_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		hookID    string
		setupMock func(*gomock.Controller) *hook_mock.MockHookService
		wantErr   bool
	}{
		{
			name:   "Successfully remove hook",
			hookID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Remove(ctx, gomock.Any()).DoAndReturn(
					func(ctx context.Context, id string) error {
						// Validate UUID format
						if _, err := uuid.Parse(id); err != nil {
							t.Errorf("Expected valid UUID, got %s", id)
						}
						return nil
					},
				)
				return hs
			},
			wantErr: false,
		},
		{
			name:   "Remove hook with invalid ID",
			hookID: "invalid-id",
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Remove(ctx, "invalid-id").Return(errors.New("invalid hook ID"))
				return hs
			},
			wantErr: true,
		},
		{
			name:   "Remove non-existent hook",
			hookID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Remove(ctx, gomock.Any()).Return(errors.New("hook not found"))
				return hs
			},
			wantErr: true,
		},
		{
			name:   "Remove hook with empty ID",
			hookID: "",
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Remove(ctx, "").Return(errors.New("hook ID is required"))
				return hs
			},
			wantErr: true,
		},
		{
			name:   "Service returns unexpected error",
			hookID: uuid.New().String(),
			setupMock: func(ctrl *gomock.Controller) *hook_mock.MockHookService {
				hs := hook_mock.NewMockHookService(ctrl)
				hs.EXPECT().Remove(ctx, gomock.Any()).Return(errors.New("storage error"))
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
			uc := hook_remove.NewUseCase(hs)

			err := uc.Execute(ctx, tc.hookID)
			if (err != nil) != tc.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestUseCase_Execute_MultipleRemoves(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hs := hook_mock.NewMockHookService(ctrl)
	uc := hook_remove.NewUseCase(hs)

	// Test removing multiple hooks in sequence
	hookIDs := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}

	for i, id := range hookIDs {
		hs.EXPECT().Remove(ctx, id).Return(nil)

		err := uc.Execute(ctx, id)
		if err != nil {
			t.Errorf("Execute() for hook %d unexpected error = %v", i, err)
		}
	}
}

func TestUseCase_Execute_PartialFailure(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hs := hook_mock.NewMockHookService(ctrl)
	uc := hook_remove.NewUseCase(hs)

	// First removal succeeds
	successID := uuid.New().String()
	hs.EXPECT().Remove(ctx, successID).Return(nil)

	err := uc.Execute(ctx, successID)
	if err != nil {
		t.Errorf("First Execute() unexpected error = %v", err)
	}

	// Second removal fails
	failID := uuid.New().String()
	expectedErr := errors.New("hook in use")
	hs.EXPECT().Remove(ctx, failID).Return(expectedErr)

	err = uc.Execute(ctx, failID)
	if err == nil {
		t.Error("Second Execute() expected error, got nil")
	}
	if err != expectedErr {
		t.Errorf("Second Execute() error = %v, want %v", err, expectedErr)
	}

	// Third removal succeeds (showing the service is still functional)
	anotherSuccessID := uuid.New().String()
	hs.EXPECT().Remove(ctx, anotherSuccessID).Return(nil)

	err = uc.Execute(ctx, anotherSuccessID)
	if err != nil {
		t.Errorf("Third Execute() unexpected error = %v", err)
	}
}
