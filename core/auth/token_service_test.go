package auth_test

import (
	"testing"

	testtarget "github.com/kyoh86/gogh/v3/core/auth"
	"github.com/kyoh86/gogh/v3/infra/github"
)

func TestTokenService(t *testing.T) {
	wantToken := github.Token{RefreshToken: "refresh-token", AccessToken: "access-token"}
	t.Run("Empty token manager always returns error", func(t *testing.T) {
		tm := testtarget.NewTokenService()
		if _, err := tm.Get("host", "owner"); err != testtarget.ErrTokenNotFound {
			t.Errorf("TokenManager.Get() returns an error %v, want %v", err, testtarget.ErrTokenNotFound)
		}
	})

	t.Run("Set and get token", func(t *testing.T) {
		tm := testtarget.NewTokenService()
		if err := tm.Set("host", "owner", wantToken); err != nil {
			t.Errorf("TokenManager.Set() returns an error %v, want to succeed", err)
		}
		got, err := tm.Get("host", "owner")
		if err != nil {
			t.Errorf("TokenManager.Get() returns an error %v, want to succeed", err)
		}
		if got != wantToken {
			t.Errorf("TokenManager.Get() = %v, want %v", got, wantToken)
		}
		t.Run("Delete token", func(t *testing.T) {
			if err := tm.Delete("host", "owner"); err != nil {
				t.Errorf("TokenManager.Delete() returns an error %v, want to succeed", err)
			}
			if _, err := tm.Get("host", "owner"); err != testtarget.ErrTokenNotFound {
				t.Errorf("TokenManager.Get() returns an error %v, want %v", err, testtarget.ErrTokenNotFound)
			}
		})
	})
}
