package auth

import "context"

// DeviceAuthResponse represents the response from a device authentication request.
type DeviceAuthResponse struct {
	// VerificationURI is the URI where the user can verify the authentication.
	VerificationURI string
	// UserCode is the code that the user needs to enter to verify the authentication.
	UserCode string
}

// Verify is a function type that takes a context and a DeviceAuthResponse
type Verify func(
	ctx context.Context,
	response DeviceAuthResponse,
) error

// AuthenticateService is an interface that defines a method for authenticating users.
type AuthenticateService interface {
	// Authenticate the user with the given host.
	// The function will return a user name and a token if the authentication is successful.
	Authenticate(ctx context.Context, host string, verify Verify) (string, *Token, error)
}
