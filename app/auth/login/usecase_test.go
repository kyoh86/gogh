package login_test

import (
	"context"
	"errors"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/auth/login"
	"github.com/kyoh86/gogh/v4/core/auth_mock"
	"github.com/kyoh86/gogh/v4/core/hosting_mock"
	"go.uber.org/mock/gomock"
	"golang.org/x/oauth2"
)

func TestUsecase_Execute(t *testing.T) {
	type mocks struct {
		tokenService   *auth_mock.MockTokenService
		authService    *auth_mock.MockAuthenticateService
		hostingService *hosting_mock.MockHostingService
	}

	tests := []struct {
		name      string
		setupMock func(m mocks)
		host      string
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success",
			setupMock: func(m mocks) {
				token := oauth2.Token{}
				m.authService.EXPECT().Authenticate(gomock.Any(), "github.com", gomock.Any()).Return("test-user", &token, nil)
				m.tokenService.EXPECT().Set("github.com", "test-user", token).Return(nil)
			},
			host:    "github.com",
			wantErr: false,
		},
		{
			name: "authentication fails",
			setupMock: func(m mocks) {
				m.authService.EXPECT().Authenticate(gomock.Any(), "github.com", gomock.Any()).Return("", nil, errors.New("auth failed"))
			},
			host:    "github.com",
			wantErr: true,
			errMsg:  "authenticating: auth failed",
		},
		{
			name: "nil token",
			setupMock: func(m mocks) {
				m.authService.EXPECT().Authenticate(gomock.Any(), "github.com", gomock.Any()).Return("test-user", nil, nil)
			},
			host:    "github.com",
			wantErr: true,
			errMsg:  "token is nil",
		},
		{
			name: "token set fails",
			setupMock: func(m mocks) {
				token := oauth2.Token{}
				m.authService.EXPECT().Authenticate(gomock.Any(), "github.com", gomock.Any()).Return("test-user", &token, nil)
				m.tokenService.EXPECT().Set("github.com", "test-user", token).Return(errors.New("set failed"))
			},
			host:    "github.com",
			wantErr: true,
			errMsg:  "setting token: set failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := mocks{
				tokenService:   auth_mock.NewMockTokenService(ctrl),
				authService:    auth_mock.NewMockAuthenticateService(ctrl),
				hostingService: hosting_mock.NewMockHostingService(ctrl),
			}

			if tt.setupMock != nil {
				tt.setupMock(m)
			}

			uc := testtarget.NewUsecase(m.tokenService, m.authService, m.hostingService)
			err := uc.Execute(context.Background(), tt.host, func(_ context.Context, resp testtarget.DeviceAuthResponse) error { return nil })

			if (err != nil) != tt.wantErr {
				t.Errorf("Usecase.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Usecase.Execute() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
