package overlay_remove_test

import (
	"context"
	"errors"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/overlay_remove"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		mockSetup func(*overlay_mock.MockOverlayService) error
		wantErr   bool
	}{
		{
			name: "正常系：オーバーレイを削除",
			id:   "overlay1",
			mockSetup: func(m *overlay_mock.MockOverlayService) error {
				m.EXPECT().Remove(gomock.Any(), "overlay1").Return(nil)
				return nil
			},
			wantErr: false,
		},
		{
			name: "エラー系：削除に失敗",
			id:   "overlay2",
			mockSetup: func(m *overlay_mock.MockOverlayService) error {
				err := errors.New("overlay not found")
				m.EXPECT().Remove(gomock.Any(), "overlay2").Return(err)
				return err
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
			err := uc.Execute(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("UseCase.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && expectedErr != nil && err != expectedErr {
				t.Errorf("Expected error doesn't match: got %v", err)
			}
		})
	}
}
