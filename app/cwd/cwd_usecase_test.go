package cwd_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kyoh86/gogh/v4/app/cwd"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("Success: Find repository by path", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockWorkspaceService := workspace_mock.NewMockWorkspaceService(ctrl)
		mockFinderService := workspace_mock.NewMockFinderService(ctrl)

		expectedLocation := repository.NewLocation(
			"/home/user/repos/github.com/owner/repo",
			"github.com",
			"owner",
			"repo",
		)

		// os.Getwd() will be called internally, so we expect FindByPath to be called with the actual working directory
		mockFinderService.EXPECT().
			FindByPath(ctx, mockWorkspaceService, gomock.Any()).
			Return(expectedLocation, nil)

		uc := cwd.NewUseCase(mockWorkspaceService, mockFinderService)
		result, err := uc.Execute(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected location, got nil")
		}
		if result.FullPath() != expectedLocation.FullPath() {
			t.Errorf("expected fullPath %s, got %s", expectedLocation.FullPath(), result.FullPath())
		}
		if result.Host() != expectedLocation.Host() {
			t.Errorf("expected host %s, got %s", expectedLocation.Host(), result.Host())
		}
		if result.Owner() != expectedLocation.Owner() {
			t.Errorf("expected owner %s, got %s", expectedLocation.Owner(), result.Owner())
		}
		if result.Name() != expectedLocation.Name() {
			t.Errorf("expected name %s, got %s", expectedLocation.Name(), result.Name())
		}
	})

	t.Run("Success: Find repository by current directory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockWorkspaceService := workspace_mock.NewMockWorkspaceService(ctrl)
		mockFinderService := workspace_mock.NewMockFinderService(ctrl)

		expectedLocation := repository.NewLocation(
			"/home/user/repos/github.com/user/project",
			"github.com",
			"user",
			"project",
		)

		// os.Getwd() will be called internally
		mockFinderService.EXPECT().
			FindByPath(ctx, mockWorkspaceService, gomock.Any()).
			Return(expectedLocation, nil)

		uc := cwd.NewUseCase(mockWorkspaceService, mockFinderService)
		result, err := uc.Execute(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected location, got nil")
		}
	})

	t.Run("Error: Repository not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockWorkspaceService := workspace_mock.NewMockWorkspaceService(ctrl)
		mockFinderService := workspace_mock.NewMockFinderService(ctrl)

		expectedErr := errors.New("not a repository")

		// os.Getwd() will be called internally
		mockFinderService.EXPECT().
			FindByPath(ctx, mockWorkspaceService, gomock.Any()).
			Return(nil, expectedErr)

		uc := cwd.NewUseCase(mockWorkspaceService, mockFinderService)
		result, err := uc.Execute(ctx)

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

	t.Run("Error: FindByPath fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockWorkspaceService := workspace_mock.NewMockWorkspaceService(ctrl)
		mockFinderService := workspace_mock.NewMockFinderService(ctrl)

		expectedErr := errors.New("finder error")

		// os.Getwd() will be called internally
		mockFinderService.EXPECT().
			FindByPath(ctx, mockWorkspaceService, gomock.Any()).
			Return(nil, expectedErr)

		uc := cwd.NewUseCase(mockWorkspaceService, mockFinderService)
		result, err := uc.Execute(ctx)

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
