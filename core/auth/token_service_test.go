package auth_test

import (
	"testing"

	testtarget "github.com/kyoh86/gogh/v3/core/auth"
)

func TestTokenServiceDeleteNonExistentToken(t *testing.T) {
	// Create a new TokenService
	tokenService := testtarget.NewTokenService()

	// Attempt to delete a token that doesn't exist
	err := tokenService.Delete("nonexistent-host", "nonexistent-owner")

	// Should not error when deleting non-existent token
	if err != nil {
		t.Errorf("TokenService.Delete() on non-existent token returned error %v, want nil", err)
	}

	// Verify that Get still returns not found
	_, err = tokenService.Get("nonexistent-host", "nonexistent-owner")
	if err != testtarget.ErrTokenNotFound {
		t.Errorf("TokenService.Get() after Delete() returned error %v, want %v", err, testtarget.ErrTokenNotFound)
	}
}

func TestTokenServiceOverwriteExistingToken(t *testing.T) {
	// Create a new TokenService
	tokenService := testtarget.NewTokenService()

	// Set initial token
	initialToken := testtarget.Token{
		AccessToken: "initial-token",
		TokenType:   "bearer",
	}
	err := tokenService.Set("github.com", "testuser", initialToken)
	if err != nil {
		t.Fatalf("TokenService.Set() failed: %v", err)
	}

	// Get and verify initial token
	retrievedToken, err := tokenService.Get("github.com", "testuser")
	if err != nil {
		t.Fatalf("TokenService.Get() failed: %v", err)
	}
	if retrievedToken.AccessToken != initialToken.AccessToken {
		t.Errorf("TokenService.Get() returned AccessToken %v, want %v",
			retrievedToken.AccessToken, initialToken.AccessToken)
	}

	// Set new token with same host/owner
	newToken := testtarget.Token{
		AccessToken: "new-token",
		TokenType:   "bearer",
	}
	err = tokenService.Set("github.com", "testuser", newToken)
	if err != nil {
		t.Fatalf("TokenService.Set() for overwrite failed: %v", err)
	}

	// Get and verify new token
	retrievedToken, err = tokenService.Get("github.com", "testuser")
	if err != nil {
		t.Fatalf("TokenService.Get() after overwrite failed: %v", err)
	}
	if retrievedToken.AccessToken != newToken.AccessToken {
		t.Errorf("TokenService.Get() after overwrite returned AccessToken %v, want %v",
			retrievedToken.AccessToken, newToken.AccessToken)
	}
}

func TestTokenService(t *testing.T) {
	wantToken := testtarget.Token{RefreshToken: "refresh-token", AccessToken: "access-token"}
	t.Run("Empty token manager always returns error", func(t *testing.T) {
		tokenService := testtarget.NewTokenService()
		if _, err := tokenService.Get("host", "owner"); err != testtarget.ErrTokenNotFound {
			t.Errorf("TokenService.Get() returns an error %v, want %v", err, testtarget.ErrTokenNotFound)
		}
	})

	t.Run("Set and get token", func(t *testing.T) {
		tokenService := testtarget.NewTokenService()
		if err := tokenService.Set("host", "owner", wantToken); err != nil {
			t.Errorf("TokenService.Set() returns an error %v, want to succeed", err)
		}
		got, err := tokenService.Get("host", "owner")
		if err != nil {
			t.Errorf("TokenService.Get() returns an error %v, want to succeed", err)
		}
		if got != wantToken {
			t.Errorf("TokenService.Get() = %v, want %v", got, wantToken)
		}
		t.Run("Delete token", func(t *testing.T) {
			if err := tokenService.Delete("host", "owner"); err != nil {
				t.Errorf("TokenService.Delete() returns an error %v, want to succeed", err)
			}
			if _, err := tokenService.Get("host", "owner"); err != testtarget.ErrTokenNotFound {
				t.Errorf("TokenService.Get() returns an error %v, want %v", err, testtarget.ErrTokenNotFound)
			}
		})
	})
}
