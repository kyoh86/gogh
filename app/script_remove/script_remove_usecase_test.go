package script_remove_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kyoh86/gogh/v4/app/script_remove"
	"github.com/kyoh86/gogh/v4/core/script_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("Success: Remove script by ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := script_mock.NewMockScriptService(ctrl)
		scriptID := "test-script-id"

		mockService.EXPECT().
			Remove(ctx, scriptID).
			Return(nil)

		uc := script_remove.NewUseCase(mockService)
		err := uc.Execute(ctx, scriptID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Success: Remove script by UUID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := script_mock.NewMockScriptService(ctrl)
		scriptID := "550e8400-e29b-41d4-a716-446655440000"

		mockService.EXPECT().
			Remove(ctx, scriptID).
			Return(nil)

		uc := script_remove.NewUseCase(mockService)
		err := uc.Execute(ctx, scriptID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Error: Script not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := script_mock.NewMockScriptService(ctrl)
		scriptID := "non-existent"
		expectedErr := errors.New("script not found")

		mockService.EXPECT().
			Remove(ctx, scriptID).
			Return(expectedErr)

		uc := script_remove.NewUseCase(mockService)
		err := uc.Execute(ctx, scriptID)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("Error: Remove fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := script_mock.NewMockScriptService(ctrl)
		scriptID := "test-script-id"
		expectedErr := errors.New("remove failed")

		mockService.EXPECT().
			Remove(ctx, scriptID).
			Return(expectedErr)

		uc := script_remove.NewUseCase(mockService)
		err := uc.Execute(ctx, scriptID)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
}