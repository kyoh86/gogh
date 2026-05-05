// Package auth provides authentication and token management for repository hosting services.
//
// This package handles device authentication flows and token storage for accessing
// remote repositories. It provides two main services:
//
//   - AuthenticateService: Manages device authentication flows with hosting services
//   - TokenService: Stores and retrieves OAuth tokens for authenticated access
//
// # Main Types
//
//   - DeviceAuthResponse: Contains verification URI and user code for device auth
//   - Token: OAuth2 token (alias for oauth2.Token)
//   - TokenEntry: Represents a stored token with its host and owner
//
// # Architecture
//
// The TokenService implements the store.Content interface for change tracking and
// persistence. Tokens are indexed by host/owner pairs and can be retrieved for
// specific repositories.
package auth
