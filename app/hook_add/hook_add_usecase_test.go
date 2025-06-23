package hook_add_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kyoh86/gogh/v4/app/hook_add"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/hook_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("Success: Add hook with all options", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := hook_mock.NewMockHookService(ctrl)
		expectedID := "test-hook-id"
		opts := hook_add.Options{
			Name:          "test-hook",
			RepoPattern:   "owner/repo",
			TriggerEvent:  "post-clone",
			OperationType: "overlay",
			OperationID:   "test-overlay-id",
		}

		expectedEntry := hook.Entry{
			Name:          opts.Name,
			RepoPattern:   opts.RepoPattern,
			TriggerEvent:  hook.Event(opts.TriggerEvent),
			OperationType: hook.OperationType(opts.OperationType),
			OperationID:   opts.OperationID,
		}

		mockService.EXPECT().
			Add(ctx, expectedEntry).
			Return(expectedID, nil)

		uc := hook_add.NewUseCase(mockService)
		result, err := uc.Execute(ctx, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != expectedID {
			t.Errorf("expected ID %s, got %s", expectedID, result)
		}
	})

	t.Run("Success: Add hook with minimal options", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := hook_mock.NewMockHookService(ctrl)
		expectedID := "minimal-hook-id"
		opts := hook_add.Options{
			Name:          "",
			RepoPattern:   "*/*",
			TriggerEvent:  "post-create",
			OperationType: "script",
			OperationID:   "script-id",
		}

		expectedEntry := hook.Entry{
			Name:          opts.Name,
			RepoPattern:   opts.RepoPattern,
			TriggerEvent:  hook.Event(opts.TriggerEvent),
			OperationType: hook.OperationType(opts.OperationType),
			OperationID:   opts.OperationID,
		}

		mockService.EXPECT().
			Add(ctx, expectedEntry).
			Return(expectedID, nil)

		uc := hook_add.NewUseCase(mockService)
		result, err := uc.Execute(ctx, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != expectedID {
			t.Errorf("expected ID %s, got %s", expectedID, result)
		}
	})

	t.Run("Error: Add fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := hook_mock.NewMockHookService(ctrl)
		expectedErr := errors.New("add failed")
		opts := hook_add.Options{
			Name:          "error-hook",
			RepoPattern:   "owner/repo",
			TriggerEvent:  "post-fork",
			OperationType: "overlay",
			OperationID:   "overlay-id",
		}

		expectedEntry := hook.Entry{
			Name:          opts.Name,
			RepoPattern:   opts.RepoPattern,
			TriggerEvent:  hook.Event(opts.TriggerEvent),
			OperationType: hook.OperationType(opts.OperationType),
			OperationID:   opts.OperationID,
		}

		mockService.EXPECT().
			Add(ctx, expectedEntry).
			Return("", expectedErr)

		uc := hook_add.NewUseCase(mockService)
		result, err := uc.Execute(ctx, opts)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if result != "" {
			t.Errorf("expected empty result, got %s", result)
		}
	})
}
