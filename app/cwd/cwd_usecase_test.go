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

		expectedPath := "/home/user/repos/owner/repo"
		expectedLocation := repository.NewLocation(
			"/home/user/repos/github.com/owner/repo",
			"github.com",
			"owner",
			"repo",
		)

		mockFinderService.EXPECT().
			FindByPath(ctx, mockWorkspaceService, expectedPath).
			Return(expectedLocation, nil)

		uc := cwd.NewUseCase(mockWorkspaceService, mockFinderService)
		result, err := uc.Execute(ctx, expectedPath)
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

		currentPath := "."
		expectedLocation := repository.NewLocation(
			"/home/user/repos/github.com/user/project",
			"github.com",
			"user",
			"project",
		)

		mockFinderService.EXPECT().
			FindByPath(ctx, mockWorkspaceService, currentPath).
			Return(expectedLocation, nil)

		uc := cwd.NewUseCase(mockWorkspaceService, mockFinderService)
		result, err := uc.Execute(ctx, currentPath)
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

		nonRepoPath := "/home/user/documents"
		expectedErr := errors.New("not a repository")

		mockFinderService.EXPECT().
			FindByPath(ctx, mockWorkspaceService, nonRepoPath).
			Return(nil, expectedErr)

		uc := cwd.NewUseCase(mockWorkspaceService, mockFinderService)
		result, err := uc.Execute(ctx, nonRepoPath)

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

		testPath := "/some/path"
		expectedErr := errors.New("finder error")

		mockFinderService.EXPECT().
			FindByPath(ctx, mockWorkspaceService, testPath).
			Return(nil, expectedErr)

		uc := cwd.NewUseCase(mockWorkspaceService, mockFinderService)
		result, err := uc.Execute(ctx, testPath)

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
