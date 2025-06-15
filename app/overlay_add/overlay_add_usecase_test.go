package overlay_add_test

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/overlay_add"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.NotNil(t, useCase)
}

func TestExecute_Success(t *testing.T) {
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
	require.NoError(t, err)

	ctx := context.Background()
	forInit := true
	relativePath := "rel/path"
	repoPattern := "org/*"

	// Set expectations - using gomock matchers
	mockService.EXPECT().
		Add(ctx, overlay.Overlay{
			RepoPattern:  repoPattern,
			ForInit:      forInit,
			RelativePath: relativePath,
		}, gomock.Any()).
		DoAndReturn(func(_ context.Context, _ overlay.Overlay, reader io.Reader) error {
			// Verify the content is being passed correctly
			data, err := io.ReadAll(reader)
			if err != nil {
				return err
			}
			assert.Equal(t, content, string(data))
			return nil
		})

	// Act
	err = useCase.Execute(ctx, forInit, relativePath, repoPattern, tempFile)

	// Assert
	assert.NoError(t, err)
}

func TestExecute_FileOpenError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := overlay_mock.NewMockOverlayService(ctrl)
	useCase := testtarget.NewUseCase(mockService)

	ctx := context.Background()
	forInit := true
	relativePath := "rel/path"
	repoPattern := "org/*"
	nonExistentFile := "/path/to/non-existent-file.txt"

	// No expectations set on mockService since it shouldn't be called

	// Act
	err := useCase.Execute(ctx, forInit, relativePath, repoPattern, nonExistentFile)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "opening source file")
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
	require.NoError(t, err)

	ctx := context.Background()
	forInit := true
	relativePath := "rel/path"
	repoPattern := "org/*"

	// Simulate an error when adding overlay
	expectedErr := errors.New("overlay add error")
	mockService.EXPECT().
		Add(ctx, overlay.Overlay{
			RepoPattern:  repoPattern,
			ForInit:      forInit,
			RelativePath: relativePath,
		}, gomock.Any()).
		Return(expectedErr)

	// Act
	err = useCase.Execute(ctx, forInit, relativePath, repoPattern, tempFile)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "adding repo-pattern")
}
