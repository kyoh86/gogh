package auth_logout_test

import (
	"context"
	"errors"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/auth_logout"
	"github.com/kyoh86/gogh/v4/core/auth_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	// テストケースの定義
	testCases := []struct {
		name        string
		host        string
		owner       string
		setupMock   func(*auth_mock.MockTokenService)
		expectError bool
	}{
		{
			name:  "正常系: トークンが正しく削除される",
			host:  "github.com",
			owner: "kyoh86",
			setupMock: func(mockService *auth_mock.MockTokenService) {
				mockService.EXPECT().Delete("github.com", "kyoh86").Return(nil)
			},
			expectError: false,
		},
		{
			name:  "異常系: トークン削除でエラーが発生",
			host:  "github.com",
			owner: "kyoh86",
			setupMock: func(mockService *auth_mock.MockTokenService) {
				mockService.EXPECT().Delete("github.com", "kyoh86").Return(errors.New("delete token error"))
			},
			expectError: true,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モックコントローラのセットアップ
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// TokenServiceのモックを作成
			mockTokenService := auth_mock.NewMockTokenService(ctrl)
			tc.setupMock(mockTokenService)

			// テスト対象のUseCaseを作成
			useCase := testtarget.NewUseCase(mockTokenService)

			// Execute実行
			err := useCase.Execute(context.Background(), tc.host, tc.owner)

			// 結果検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーは期待されませんでしたが、エラーが発生しました: %v", err)
			}
		})
	}
}
