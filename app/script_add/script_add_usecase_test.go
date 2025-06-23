package script_add_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/script_add"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/script_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	testID := uuid.New()
	testTime := time.Now()

	t.Run("Success: Add script with name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := script_mock.NewMockScriptService(ctrl)
		content := strings.NewReader("print('hello world')")
		expectedScript := script.ConcreteScript(testID, "test-script", testTime, testTime)

		mockService.EXPECT().
			Add(ctx, script.Entry{Name: "test-script", Content: content}).
			Return(testID.String(), nil)

		mockService.EXPECT().
			Get(ctx, testID.String()).
			Return(expectedScript, nil)

		uc := script_add.NewUseCase(mockService)
		result, err := uc.Execute(ctx, "test-script", content)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID() != testID.String() {
			t.Errorf("expected ID %s, got %s", testID.String(), result.ID())
		}
		if result.Name() != "test-script" {
			t.Errorf("expected name 'test-script', got %s", result.Name())
		}
	})

	t.Run("Success: Add script without name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := script_mock.NewMockScriptService(ctrl)
		content := strings.NewReader("print('hello world')")
		expectedScript := script.ConcreteScript(testID, "", testTime, testTime)

		mockService.EXPECT().
			Add(ctx, script.Entry{Name: "", Content: content}).
			Return(testID.String(), nil)

		mockService.EXPECT().
			Get(ctx, testID.String()).
			Return(expectedScript, nil)

		uc := script_add.NewUseCase(mockService)
		result, err := uc.Execute(ctx, "", content)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.ID() != testID.String() {
			t.Errorf("expected ID %s, got %s", testID.String(), result.ID())
		}
		if result.Name() != "" {
			t.Errorf("expected empty name, got %s", result.Name())
		}
	})

	t.Run("Error: Add fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := script_mock.NewMockScriptService(ctrl)
		content := strings.NewReader("print('hello world')")
		expectedErr := errors.New("add failed")

		mockService.EXPECT().
			Add(ctx, script.Entry{Name: "test-script", Content: content}).
			Return("", expectedErr)

		uc := script_add.NewUseCase(mockService)
		result, err := uc.Execute(ctx, "test-script", content)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if result != nil {
			t.Errorf("expected nil result, got %v", result)
		}
	})

	t.Run("Error: Get fails after successful add", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := script_mock.NewMockScriptService(ctrl)
		content := strings.NewReader("print('hello world')")
		expectedErr := errors.New("get failed")

		mockService.EXPECT().
			Add(ctx, script.Entry{Name: "test-script", Content: content}).
			Return(testID.String(), nil)

		mockService.EXPECT().
			Get(ctx, testID.String()).
			Return(nil, expectedErr)

		uc := script_add.NewUseCase(mockService)
		result, err := uc.Execute(ctx, "test-script", content)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if result != nil {
			t.Errorf("expected nil result, got %v", result)
		}
	})
}