package auth_test

import (
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/core/auth"
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

func TestTokenServiceHas(t *testing.T) {
	// Create a new TokenService
	tokenService := testtarget.NewTokenService()

	// Initially should not have any tokens
	if tokenService.Has("github.com", "testuser") {
		t.Error("TokenService.Has() returned true for non-existent token")
	}

	// Add a token
	token := testtarget.Token{AccessToken: "test-token"}
	err := tokenService.Set("github.com", "testuser", token)
	if err != nil {
		t.Fatalf("TokenService.Set() failed: %v", err)
	}

	// Should now have the token
	if !tokenService.Has("github.com", "testuser") {
		t.Error("TokenService.Has() returned false for existing token")
	}

	// Should not have token for different host/owner
	if tokenService.Has("gitlab.com", "testuser") {
		t.Error("TokenService.Has() returned true for wrong host")
	}
	if tokenService.Has("github.com", "otheruser") {
		t.Error("TokenService.Has() returned true for wrong owner")
	}
}

func TestTokenServiceEntries(t *testing.T) {
	// Create a new TokenService
	tokenService := testtarget.NewTokenService()

	// Initially should have no entries
	entries := tokenService.Entries()
	if len(entries) != 0 {
		t.Errorf("TokenService.Entries() returned %d entries, expected 0", len(entries))
	}

	// Add some tokens
	token1 := testtarget.Token{AccessToken: "token1"}
	token2 := testtarget.Token{AccessToken: "token2"}

	err := tokenService.Set("github.com", "user1", token1)
	if err != nil {
		t.Fatalf("TokenService.Set() failed: %v", err)
	}

	err = tokenService.Set("gitlab.com", "user2", token2)
	if err != nil {
		t.Fatalf("TokenService.Set() failed: %v", err)
	}

	// Should now have 2 entries
	entries = tokenService.Entries()
	if len(entries) != 2 {
		t.Errorf("TokenService.Entries() returned %d entries, expected 2", len(entries))
	}

	// Verify entry contents
	found1, found2 := false, false
	for _, entry := range entries {
		if entry.Host == "github.com" && entry.Owner == "user1" {
			if entry.Token != token1 {
				t.Errorf("Token mismatch for github.com/user1")
			}
			found1 = true
		}
		if entry.Host == "gitlab.com" && entry.Owner == "user2" {
			if entry.Token != token2 {
				t.Errorf("Token mismatch for gitlab.com/user2")
			}
			found2 = true
		}
	}

	if !found1 {
		t.Error("Entry for github.com/user1 not found")
	}
	if !found2 {
		t.Error("Entry for gitlab.com/user2 not found")
	}
}

func TestTokenServiceChanges(t *testing.T) {
	// Create a new TokenService
	tokenService := testtarget.NewTokenService()

	// Initially should not have changes
	if tokenService.HasChanges() {
		t.Error("New TokenService should not have changes")
	}

	// Set a token should cause changes
	token := testtarget.Token{AccessToken: "test-token"}
	err := tokenService.Set("github.com", "testuser", token)
	if err != nil {
		t.Fatalf("TokenService.Set() failed: %v", err)
	}

	if !tokenService.HasChanges() {
		t.Error("TokenService should have changes after Set")
	}

	// Mark as saved
	tokenService.MarkSaved()
	if tokenService.HasChanges() {
		t.Error("TokenService should not have changes after MarkSaved")
	}

	// Delete a token should cause changes
	err = tokenService.Delete("github.com", "testuser")
	if err != nil {
		t.Fatalf("TokenService.Delete() failed: %v", err)
	}

	if !tokenService.HasChanges() {
		t.Error("TokenService should have changes after Delete")
	}
}

func TestTokenEntryString(t *testing.T) {
	// Create a token entry
	entry := testtarget.TokenEntry{
		Host:  "github.com",
		Owner: "testuser",
		Token: testtarget.Token{AccessToken: "sensitive-token-value"},
	}

	// Get string representation
	str := entry.String()

	// Should mask the token
	if str != "*****@github.com/testuser" {
		t.Errorf("TokenEntry.String() returned unexpected value: %s", str)
	}

	// Should not contain the actual token
	if str != "*****@github.com/testuser" && str != "sensitive-token-value" {
		t.Errorf("TokenEntry.String() might expose token value: %s", str)
	}
}
