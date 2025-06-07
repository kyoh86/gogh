package overlay_find_test

import (
	"context"
	"errors"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/overlay_find"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/repository_mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewUseCase(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParser := repository_mock.NewMockReferenceParser(ctrl)
	mockService := overlay_mock.NewMockOverlayService(ctrl)

	// Act
	useCase := testtarget.NewUseCase(mockParser, mockService)

	// Assert
	assert.NotNil(t, useCase)
}

func TestExecute_ParsingError(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParser := repository_mock.NewMockReferenceParser(ctrl)
	mockService := overlay_mock.NewMockOverlayService(ctrl)
	useCase := testtarget.NewUseCase(mockParser, mockService)

	ctx := context.Background()
	refsInput := "invalid:reference"
	expectedErr := errors.New("parse error")

	// Set up mock to return an error when parsing
	mockParser.EXPECT().
		ParseWithAlias(refsInput).
		Return(nil, expectedErr)

	// The overlay service should not be called if parsing fails

	// Act
	result := useCase.Execute(ctx, refsInput)

	// Assert - collect results from the iterator
	var overlays []*testtarget.Overlay
	var gotErr error

	result(func(overlay *testtarget.Overlay, err error) bool {
		if err != nil {
			gotErr = err
			return false
		}
		overlays = append(overlays, overlay)
		return true
	})

	assert.Error(t, gotErr)
	assert.Contains(t, gotErr.Error(), "parsing reference")
	assert.Contains(t, gotErr.Error(), expectedErr.Error())
	assert.Empty(t, overlays)
}

func TestExecute_NoMatchingOverlays(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParser := repository_mock.NewMockReferenceParser(ctrl)
	mockService := overlay_mock.NewMockOverlayService(ctrl)
	useCase := testtarget.NewUseCase(mockParser, mockService)

	ctx := context.Background()
	refsInput := "github.com/kyoh86/gogh"

	// Create a reference with alias
	ref := &repository.ReferenceWithAlias{
		Reference: repository.NewReference("github.com", "kyoh86", "gogh"),
	}

	// Set up mock to successfully parse the reference
	mockParser.EXPECT().
		ParseWithAlias(refsInput).
		Return(ref, nil)

	// Set up empty overlays list
	mockService.EXPECT().
		ListOverlays().
		Return(func(yield func(*overlay.Overlay, error) bool) {
			// Return no overlays
		})

	// Act
	result := useCase.Execute(ctx, refsInput)

	// Assert - collect results from the iterator
	var overlays []*testtarget.Overlay
	var gotErr error

	result(func(overlay *testtarget.Overlay, err error) bool {
		if err != nil {
			gotErr = err
			return false
		}
		overlays = append(overlays, overlay)
		return true
	})

	assert.NoError(t, gotErr)
	assert.Empty(t, overlays)
}

func TestExecute_WithMatchingOverlays(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParser := repository_mock.NewMockReferenceParser(ctrl)
	mockService := overlay_mock.NewMockOverlayService(ctrl)
	useCase := testtarget.NewUseCase(mockParser, mockService)

	ctx := context.Background()
	refsInput := "github.com/kyoh86/gogh"

	// Create a reference with alias
	ref := &repository.ReferenceWithAlias{
		Reference: repository.NewReference("github.com", "kyoh86", "gogh"),
	}

	// Test overlays
	overlays := []*overlay.Overlay{
		{
			RepoPattern:  "github.com/kyoh86/*",
			RelativePath: "setup.sh",
			ForInit:      true,
		},
		{
			RepoPattern:  "github.com/other/*",
			RelativePath: "config.yml",
			ForInit:      true,
		},
		{
			RepoPattern:  "github.com/kyoh86/gogh",
			RelativePath: "specific.txt",
			ForInit:      false,
		},
	}

	// Set up mock to successfully parse the reference
	mockParser.EXPECT().
		ParseWithAlias(refsInput).
		Return(ref, nil)

	// Set up mock to return all overlays
	mockService.EXPECT().
		ListOverlays().
		Return(func(yield func(*overlay.Overlay, error) bool) {
			for _, ov := range overlays {
				if !yield(ov, nil) {
					return
				}
			}
		})

	// Act
	result := useCase.Execute(ctx, refsInput)

	// Assert - collect results from the iterator
	var matchingOverlays []*testtarget.Overlay
	var gotErr error

	result(func(overlay *testtarget.Overlay, err error) bool {
		if err != nil {
			gotErr = err
			return false
		}
		matchingOverlays = append(matchingOverlays, overlay)
		return true
	})

	assert.NoError(t, gotErr)
	require.Len(t, matchingOverlays, 2)

	// Verify the matching overlays are returned (1st and 3rd should match)
	assert.Equal(t, "setup.sh", matchingOverlays[0].RelativePath)
	assert.Equal(t, "specific.txt", matchingOverlays[1].RelativePath)
}

func TestExecute_WithErrorInOverlays(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockParser := repository_mock.NewMockReferenceParser(ctrl)
	mockService := overlay_mock.NewMockOverlayService(ctrl)
	useCase := testtarget.NewUseCase(mockParser, mockService)

	ctx := context.Background()
	refsInput := "github.com/kyoh86/gogh"

	// Create a reference with alias
	ref := &repository.ReferenceWithAlias{
		Reference: repository.NewReference("github.com", "kyoh86", "gogh"),
	}

	expectedErr := errors.New("overlay list error")

	// Set up mock to successfully parse the reference
	mockParser.EXPECT().
		ParseWithAlias(refsInput).
		Return(ref, nil)

	// Set up mock to return an error
	mockService.EXPECT().
		ListOverlays().
		Return(func(yield func(*overlay.Overlay, error) bool) {
			// Return one overlay, then an error
			overlay1 := &overlay.Overlay{
				RepoPattern:  "github.com/kyoh86/*",
				RelativePath: "setup.sh",
				ForInit:      true,
			}

			if !yield(overlay1, nil) {
				return
			}

			yield(nil, expectedErr)
		})

	// Act
	result := useCase.Execute(ctx, refsInput)

	// Assert - collect results from the iterator
	var overlays []*testtarget.Overlay
	var gotErr error

	result(func(overlay *testtarget.Overlay, err error) bool {
		if err != nil {
			gotErr = err
			return false
		}
		overlays = append(overlays, overlay)
		return true
	})

	assert.Error(t, gotErr)
	assert.Equal(t, expectedErr, gotErr)
	assert.Len(t, overlays, 1) // We should get the first overlay before the error
}
