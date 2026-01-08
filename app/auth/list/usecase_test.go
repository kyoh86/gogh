package list_test

import (
	"context"
	"reflect"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/auth/list"
	"github.com/kyoh86/gogh/v4/core/auth"
	"github.com/kyoh86/gogh/v4/core/auth_mock"
	"go.uber.org/mock/gomock"
	"golang.org/x/oauth2"
)

func TestUsecase_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(mockTokenService *auth_mock.MockTokenService)
		want    []auth.TokenEntry
		wantErr bool
	}{
		{
			name: "success: return token entries",
			setup: func(mockTokenService *auth_mock.MockTokenService) {
				entries := []auth.TokenEntry{
					{Host: "github.com", Owner: "user1", Token: oauth2.Token{AccessToken: "token1"}},
					{Host: "github.com", Owner: "user2", Token: oauth2.Token{AccessToken: "token2"}},
				}
				mockTokenService.EXPECT().Entries().Return(entries)
			},
			want: []auth.TokenEntry{
				{Host: "github.com", Owner: "user1", Token: oauth2.Token{AccessToken: "token1"}},
				{Host: "github.com", Owner: "user2", Token: oauth2.Token{AccessToken: "token2"}},
			},
			wantErr: false,
		},
		{
			name: "success: return empty entries",
			setup: func(mockTokenService *auth_mock.MockTokenService) {
				mockTokenService.EXPECT().Entries().Return([]auth.TokenEntry{})
			},
			want:    []auth.TokenEntry{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockTokenService := auth_mock.NewMockTokenService(ctrl)
			tt.setup(mockTokenService)

			uc := testtarget.NewUsecase(mockTokenService)
			got, err := uc.Execute(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Usecase.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
