package overlay_add_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kyoh86/gogh/v4/app/overlay_add"
	"github.com/kyoh86/gogh/v4/core/workspace"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestNewUseCase(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOverlayService := workspace_mock.NewMockOverlayService(ctrl)

	// Test that constructor returns a non-nil use case
	useCase := overlay_add.NewUseCase(mockOverlayService)
	if useCase == nil {
		t.Fatal("Expected non-nil use case")
	}
}

func TestExecute_ExistingPattern(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOverlayService := workspace_mock.NewMockOverlayService(ctrl)
	useCase := overlay_add.NewUseCase(mockOverlayService)

	ctx := context.Background()
	pattern := "github.com/*"
	sourcePath := "/src/template.go"
	targetPath := "target/file.go"

	// Existing files for the pattern
	existingFiles := []workspace.OverlayFile{
		{SourcePath: "/src/existing.go", TargetPath: "target/existing.go"},
	}

	// Expected files after adding the new one
	expectedFiles := []workspace.OverlayFile{
		{SourcePath: "/src/existing.go", TargetPath: "target/existing.go"},
		{SourcePath: sourcePath, TargetPath: targetPath},
	}

	// Setup expectations
	mockOverlayService.EXPECT().GetPatterns().Return([]workspace.OverlayPattern{
		{Pattern: pattern, Files: existingFiles},
		{Pattern: "gitlab.com/*", Files: []workspace.OverlayFile{}},
	})

	mockOverlayService.EXPECT().AddPattern(pattern, expectedFiles).Return(nil)

	// Execute the use case
	err := useCase.Execute(ctx, pattern, sourcePath, targetPath)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestExecute_NewPattern(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOverlayService := workspace_mock.NewMockOverlayService(ctrl)
	useCase := overlay_add.NewUseCase(mockOverlayService)

	ctx := context.Background()
	pattern := "github.com/new/*"
	sourcePath := "/src/template.go"
	targetPath := "target/file.go"

	// Expected files (just the new one since pattern is new)
	expectedFiles := []workspace.OverlayFile{
		{SourcePath: sourcePath, TargetPath: targetPath},
	}

	// Setup expectations
	// Return patterns that don't include our new pattern
	mockOverlayService.EXPECT().GetPatterns().Return([]workspace.OverlayPattern{
		{Pattern: "github.com/existing/*", Files: []workspace.OverlayFile{}},
		{Pattern: "gitlab.com/*", Files: []workspace.OverlayFile{}},
	})

	mockOverlayService.EXPECT().AddPattern(pattern, expectedFiles).Return(nil)

	// Execute the use case
	err := useCase.Execute(ctx, pattern, sourcePath, targetPath)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestExecute_EmptyPatterns(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOverlayService := workspace_mock.NewMockOverlayService(ctrl)
	useCase := overlay_add.NewUseCase(mockOverlayService)

	ctx := context.Background()
	pattern := "github.com/*"
	sourcePath := "/src/template.go"
	targetPath := "target/file.go"

	// Expected files (just the new one since there are no existing patterns)
	expectedFiles := []workspace.OverlayFile{
		{SourcePath: sourcePath, TargetPath: targetPath},
	}

	// Setup expectations
	mockOverlayService.EXPECT().GetPatterns().Return([]workspace.OverlayPattern{})
	mockOverlayService.EXPECT().AddPattern(pattern, expectedFiles).Return(nil)

	// Execute the use case
	err := useCase.Execute(ctx, pattern, sourcePath, targetPath)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestExecute_AddPatternError(t *testing.T) {
	// Setup
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOverlayService := workspace_mock.NewMockOverlayService(ctrl)
	useCase := overlay_add.NewUseCase(mockOverlayService)

	ctx := context.Background()
	pattern := "github.com/*"
	sourcePath := "/src/template.go"
	targetPath := "target/file.go"

	expectedError := errors.New("add pattern error")

	// Setup expectations
	mockOverlayService.EXPECT().GetPatterns().Return([]workspace.OverlayPattern{})
	mockOverlayService.EXPECT().AddPattern(pattern, gomock.Any()).Return(expectedError)

	// Execute the use case
	err := useCase.Execute(ctx, pattern, sourcePath, targetPath)

	// Verify
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(errors.Unwrap(err), expectedError) {
		t.Fatalf("Expected error containing %v, got %v", expectedError, err)
	}
}
