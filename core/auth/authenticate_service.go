package auth

import "context"

type DeviceAuthResponse struct {
	VerificationURI string
	UserCode        string
}

type Verify func(
	ctx context.Context,
	response DeviceAuthResponse,
) error

type AuthenticateService interface {
	// Authenticate the user with the given host.
	// The function will return a user name and a token if the authentication is successful.
	Authenticate(ctx context.Context, host string, verify Verify) (string, *Token, error)
}
