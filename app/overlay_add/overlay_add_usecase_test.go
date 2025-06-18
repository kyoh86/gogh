package overlay_add_test

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/overlay_add"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"go.uber.org/mock/gomock"
)

func TestNewUseCase(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := overlay_mock.NewMockOverlayService(ctrl)

	// Act
	useCase := testtarget.NewUseCase(mockService)

	// Assert
	if useCase == nil {
		t.Fatal("expected useCase to be non-nil")
	}
}

func TestExecute_Success(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := overlay_mock.NewMockOverlayService(ctrl)
	useCase := testtarget.NewUseCase(mockService)

	ctx := context.Background()
	name := "test-overlay"
	relativePath := "rel/path"
	content := strings.NewReader("test content")

	// Set expectations - using gomock matchers
	mockService.EXPECT().
		Add(ctx, overlay.Entry{
			Name:         name,
			RelativePath: relativePath,
			Content:      content,
		}).
		DoAndReturn(func(_ context.Context, entry overlay.Entry) (string, error) {
			// Verify the content is being passed correctly
			data, err := io.ReadAll(entry.Content)
			if err != nil {
				return "", err
			}
			if string(data) != "test content" {
				return "", errors.New("content mismatch")
			}
			return "test-id", nil
		})

	// Act
	id, err := useCase.Execute(ctx, name, relativePath, content)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != "test-id" {
		t.Fatalf("expected id to be 'test-id', got %s", id)
	}
}

func TestExecute_AddOverlayError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := overlay_mock.NewMockOverlayService(ctrl)
	useCase := testtarget.NewUseCase(mockService)

	// Create a temporary file for testing
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test-source.txt")
	content := "test content"
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	ctx := context.Background()
	name := "test-overlay"
	relativePath := "rel/path"

	// Simulate an error when adding overlay
	expectedErr := errors.New("overlay add error")
	mockService.EXPECT().
		Add(ctx, overlay.Entry{
			Name:         name,
			RelativePath: relativePath,
		}).
		Return("", expectedErr)

	// Act
	if _, err := useCase.Execute(ctx, name, relativePath, nil); !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}
