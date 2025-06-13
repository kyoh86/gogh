package overlay_show_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/overlay_show"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"go.uber.org/mock/gomock"
)

type readCloserMock struct {
	*bytes.Reader
	closed bool
}

func (r *readCloserMock) Close() error {
	r.closed = true
	return nil
}

func TestUseCase_Execute(t *testing.T) {
	tests := []struct {
		name         string
		repoPattern  string
		forInit      bool
		relativePath string
		mockSetup    func(*overlay_mock.MockOverlayService) (io.ReadCloser, error)
		wantErr      bool
	}{
		{
			name:         "正常系：オーバーレイコンテンツを表示",
			repoPattern:  "example/repo",
			forInit:      false,
			relativePath: "path/to/file",
			mockSetup: func(m *overlay_mock.MockOverlayService) (io.ReadCloser, error) {
				content := &readCloserMock{
					Reader: bytes.NewReader([]byte("overlay content")),
				}
				m.EXPECT().OpenOverlayContent(gomock.Any(), overlay.Overlay{
					RepoPattern:  "example/repo",
					ForInit:      false,
					RelativePath: "path/to/file",
				}).Return(content, nil)
				return content, nil
			},
			wantErr: false,
		},
		{
			name:         "エラー系：オーバーレイを開けない",
			repoPattern:  "error/repo",
			forInit:      true,
			relativePath: "file.txt",
			mockSetup: func(m *overlay_mock.MockOverlayService) (io.ReadCloser, error) {
				m.EXPECT().OpenOverlayContent(gomock.Any(), overlay.Overlay{
					RepoPattern:  "error/repo",
					ForInit:      true,
					RelativePath: "file.txt",
				}).Return(nil, errors.New("overlay not found"))
				return nil, errors.New("overlay not found")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := overlay_mock.NewMockOverlayService(ctrl)
			content, _ := tt.mockSetup(mockService)

			w := &bytes.Buffer{}
			uc := testtarget.NewUseCaseContent(mockService, w, 60)
			err := uc.Execute(context.Background(), &overlay.Overlay{RepoPattern: tt.repoPattern, ForInit: tt.forInit, RelativePath: tt.relativePath})

			if (err != nil) != tt.wantErr {
				t.Errorf("UseCase.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			// ReadCloserが閉じられたことを確認（エラーがなかった場合）
			if err == nil && content != nil {
				if mock, ok := content.(*readCloserMock); ok && !mock.closed {
					t.Error("ReadCloser was not closed")
				}
			}
		})
	}
}
