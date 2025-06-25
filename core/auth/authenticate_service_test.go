package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kyoh86/gogh/v4/core/auth"
	"golang.org/x/oauth2"
)

// Mock implementation of AuthenticateService
type mockAuthenticateService struct {
	authenticateFunc func(ctx context.Context, host string, verify auth.Verify) (string, *auth.Token, error)
}

func (m *mockAuthenticateService) Authenticate(ctx context.Context, host string, verify auth.Verify) (string, *auth.Token, error) {
	if m.authenticateFunc != nil {
		return m.authenticateFunc(ctx, host, verify)
	}
	return "", nil, errors.New("not implemented")
}

func TestDeviceAuthResponse(t *testing.T) {
	// Test DeviceAuthResponse struct
	response := auth.DeviceAuthResponse{
		VerificationURI: "https://github.com/login/device",
		UserCode:        "ABCD-1234",
	}

	if response.VerificationURI != "https://github.com/login/device" {
		t.Errorf("expected VerificationURI %q, got %q", "https://github.com/login/device", response.VerificationURI)
	}

	if response.UserCode != "ABCD-1234" {
		t.Errorf("expected UserCode %q, got %q", "ABCD-1234", response.UserCode)
	}
}

func TestVerifyFunction(t *testing.T) {
	// Test Verify function type
	ctx := context.Background()
	response := auth.DeviceAuthResponse{
		VerificationURI: "https://example.com/verify",
		UserCode:        "TEST-CODE",
	}

	t.Run("successful verification", func(t *testing.T) {
		var called bool
		var receivedResponse auth.DeviceAuthResponse

		verify := auth.Verify(func(_ context.Context, response auth.DeviceAuthResponse) error {
			called = true
			receivedResponse = response
			return nil
		})

		err := verify(ctx, response)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if !called {
			t.Error("verify function was not called")
		}

		if receivedResponse.VerificationURI != response.VerificationURI {
			t.Errorf("expected VerificationURI %q, got %q", response.VerificationURI, receivedResponse.VerificationURI)
		}

		if receivedResponse.UserCode != response.UserCode {
			t.Errorf("expected UserCode %q, got %q", response.UserCode, receivedResponse.UserCode)
		}
	})

	t.Run("verification with error", func(t *testing.T) {
		verify := auth.Verify(func(_ context.Context, _ auth.DeviceAuthResponse) error {
			return errors.New("verification failed")
		})

		err := verify(ctx, response)
		if err == nil {
			t.Error("expected error from verify function")
		}
	})
}

func TestAuthenticateService_Authenticate(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name             string
		host             string
		setupMock        func() *mockAuthenticateService
		setupVerify      func() auth.Verify
		wantUser         string
		wantTokenPresent bool
		wantErr          bool
		errMsg           string
	}{
		{
			name: "successful authentication",
			host: "github.com",
			setupMock: func() *mockAuthenticateService {
				return &mockAuthenticateService{
					authenticateFunc: func(ctx context.Context, host string, verify auth.Verify) (string, *auth.Token, error) {
						// Simulate calling verify
						response := auth.DeviceAuthResponse{
							VerificationURI: "https://github.com/login/device",
							UserCode:        "TEST-1234",
						}
						if err := verify(ctx, response); err != nil {
							return "", nil, err
						}

						token := &oauth2.Token{
							AccessToken: "test-token",
							TokenType:   "Bearer",
						}
						return "testuser", token, nil
					},
				}
			},
			setupVerify: func() auth.Verify {
				return func(ctx context.Context, response auth.DeviceAuthResponse) error {
					// Simulate successful verification
					return nil
				}
			},
			wantUser:         "testuser",
			wantTokenPresent: true,
			wantErr:          false,
		},
		{
			name: "authentication failure",
			host: "github.com",
			setupMock: func() *mockAuthenticateService {
				return &mockAuthenticateService{
					authenticateFunc: func(ctx context.Context, host string, verify auth.Verify) (string, *auth.Token, error) {
						return "", nil, errors.New("authentication failed")
					},
				}
			},
			setupVerify: func() auth.Verify {
				return func(ctx context.Context, response auth.DeviceAuthResponse) error {
					return nil
				}
			},
			wantUser:         "",
			wantTokenPresent: false,
			wantErr:          true,
			errMsg:           "authentication failed",
		},
		{
			name: "verify error",
			host: "github.com",
			setupMock: func() *mockAuthenticateService {
				return &mockAuthenticateService{
					authenticateFunc: func(ctx context.Context, host string, verify auth.Verify) (string, *auth.Token, error) {
						response := auth.DeviceAuthResponse{
							VerificationURI: "https://github.com/login/device",
							UserCode:        "TEST-1234",
						}
						if err := verify(ctx, response); err != nil {
							return "", nil, err
						}
						return "testuser", nil, nil
					},
				}
			},
			setupVerify: func() auth.Verify {
				return func(ctx context.Context, response auth.DeviceAuthResponse) error {
					return errors.New("verification failed")
				}
			},
			wantUser:         "",
			wantTokenPresent: false,
			wantErr:          true,
			errMsg:           "verification failed",
		},
		{
			name: "different hosts",
			host: "gitlab.com",
			setupMock: func() *mockAuthenticateService {
				return &mockAuthenticateService{
					authenticateFunc: func(ctx context.Context, host string, verify auth.Verify) (string, *auth.Token, error) {
						if host != "gitlab.com" {
							return "", nil, errors.New("unexpected host")
						}
						token := &oauth2.Token{
							AccessToken: "gitlab-token",
							TokenType:   "Bearer",
						}
						return "gitlabuser", token, nil
					},
				}
			},
			setupVerify: func() auth.Verify {
				return func(ctx context.Context, response auth.DeviceAuthResponse) error {
					return nil
				}
			},
			wantUser:         "gitlabuser",
			wantTokenPresent: true,
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := tt.setupMock()
			verify := tt.setupVerify()

			user, token, err := service.Authenticate(ctx, tt.host, verify)

			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.errMsg != "" && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Authenticate() error = %v, want error %q", err, tt.errMsg)
			}

			if user != tt.wantUser {
				t.Errorf("Authenticate() user = %q, want %q", user, tt.wantUser)
			}

			if (token != nil) != tt.wantTokenPresent {
				t.Errorf("Authenticate() token present = %v, want %v", token != nil, tt.wantTokenPresent)
			}

			if token != nil && token.AccessToken == "" {
				t.Error("Authenticate() returned token with empty AccessToken")
			}
		})
	}
}

func TestAuthenticateService_Interface(t *testing.T) {
	// This test ensures that mockAuthenticateService implements AuthenticateService
	var _ auth.AuthenticateService = (*mockAuthenticateService)(nil)
}

func TestAuthenticateWithContext(t *testing.T) {
	// Test context cancellation
	service := &mockAuthenticateService{
		authenticateFunc: func(ctx context.Context, host string, verify auth.Verify) (string, *auth.Token, error) {
			select {
			case <-ctx.Done():
				return "", nil, ctx.Err()
			default:
				token := &oauth2.Token{
					AccessToken: "test-token",
				}
				return "testuser", token, nil
			}
		},
	}

	// Test with cancelled context
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	verify := func(ctx context.Context, response auth.DeviceAuthResponse) error {
		return nil
	}

	_, _, err := service.Authenticate(cancelledCtx, "github.com", verify)
	if err == nil {
		t.Error("expected error with cancelled context")
	}

	// Test with valid context
	validCtx := context.Background()
	user, token, err := service.Authenticate(validCtx, "github.com", verify)
	if err != nil {
		t.Errorf("unexpected error with valid context: %v", err)
	}
	if user != "testuser" {
		t.Errorf("expected user %q, got %q", "testuser", user)
	}
	if token == nil || token.AccessToken != "test-token" {
		t.Error("unexpected token value")
	}
}
