package logout_test

import (
	"context"
	"errors"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/auth/logout"
	"github.com/kyoh86/gogh/v4/core/auth_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		host        string
		owner       string
		setupMock   func(*auth_mock.MockTokenService)
		expectError bool
	}{
		{
			name:  "Normal case: Token is deleted correctly",
			host:  "github.com",
			owner: "kyoh86",
			setupMock: func(mockService *auth_mock.MockTokenService) {
				mockService.EXPECT().Delete("github.com", "kyoh86").Return(nil)
			},
			expectError: false,
		},
		{
			name:  "Error case: Error occurs during token deletion",
			host:  "github.com",
			owner: "kyoh86",
			setupMock: func(mockService *auth_mock.MockTokenService) {
				mockService.EXPECT().Delete("github.com", "kyoh86").Return(errors.New("delete token error"))
			},
			expectError: true,
		},
	}

	// Execute each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up mock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create TokenService mock
			mockTokenService := auth_mock.NewMockTokenService(ctrl)
			tc.setupMock(mockTokenService)

			// Create UseCase under test
			useCase := testtarget.NewUseCase(mockTokenService)

			// Execute
			err := useCase.Execute(context.Background(), tc.host, tc.owner)

			// Verify result
			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
