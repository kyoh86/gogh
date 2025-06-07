package overlay_remove_test

import (
	"context"
	"errors"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/overlay_remove"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	tests := []struct {
		name         string
		forInit      bool
		relativePath string
		repoPattern  string
		mockSetup    func(*overlay_mock.MockOverlayService) error
		wantErr      bool
	}{
		{
			name:         "正常系：オーバーレイを削除",
			forInit:      false,
			relativePath: "path/to/file",
			repoPattern:  "example/repo",
			mockSetup: func(m *overlay_mock.MockOverlayService) error {
				m.EXPECT().RemoveOverlay(gomock.Any(), overlay.Overlay{
					RepoPattern:  "example/repo",
					ForInit:      false,
					RelativePath: "path/to/file",
				}).Return(nil)
				return nil
			},
			wantErr: false,
		},
		{
			name:         "正常系：初期化用オーバーレイを削除",
			forInit:      true,
			relativePath: "config.yaml",
			repoPattern:  "example/repo",
			mockSetup: func(m *overlay_mock.MockOverlayService) error {
				m.EXPECT().RemoveOverlay(gomock.Any(), overlay.Overlay{
					RepoPattern:  "example/repo",
					ForInit:      true,
					RelativePath: "config.yaml",
				}).Return(nil)
				return nil
			},
			wantErr: false,
		},
		{
			name:         "エラー系：削除に失敗",
			forInit:      false,
			relativePath: "non-existent.txt",
			repoPattern:  "error/repo",
			mockSetup: func(m *overlay_mock.MockOverlayService) error {
				m.EXPECT().RemoveOverlay(gomock.Any(), overlay.Overlay{
					RepoPattern:  "error/repo",
					ForInit:      false,
					RelativePath: "non-existent.txt",
				}).Return(errors.New("overlay not found"))
				return errors.New("overlay not found")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := overlay_mock.NewMockOverlayService(ctrl)
			expectedErr := tt.mockSetup(mockService)

			uc := testtarget.NewUseCase(mockService)
			err := uc.Execute(context.Background(), tt.forInit, tt.relativePath, tt.repoPattern)

			if (err != nil) != tt.wantErr {
				t.Errorf("UseCase.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && expectedErr != nil && err.Error() != "removing entry "+tt.relativePath+" for "+tt.repoPattern+": "+expectedErr.Error() {
				t.Errorf("Expected error message doesn't match: got %v", err)
			}
		})
	}
}
