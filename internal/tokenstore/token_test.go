package tokenstore_test

import (
	"testing"

	"github.com/kyoh86/gogh/v3/internal/github"
	testtarget "github.com/kyoh86/gogh/v3/internal/tokenstore"
)

func TestTokenManager(t *testing.T) {
	wantToken := github.Token{RefreshToken: "refresh-token", AccessToken: "access-token"}
	t.Run("Empty token manager always returns error", func(t *testing.T) {
		var tm testtarget.TokenManager
		if _, err := tm.Get("host", "owner"); err != testtarget.ErrNoHost {
			t.Errorf("TokenManager.Get() returns an error %v, want %v", err, testtarget.ErrNoHost)
		}
		if gotHost, gotToken := tm.GetDefaultKey(); gotHost != github.DefaultHost || gotToken != "" {
			t.Errorf("TokenManager.GetDefaultKey() = %v, %v, want %v, %v", gotHost, gotToken, "", "")
		}
	})

	t.Run("Set and get token", func(t *testing.T) {
		var tm testtarget.TokenManager
		tm.Set("host", "owner", wantToken)
		got, err := tm.Get("host", "owner")
		if err != nil {
			t.Errorf("TokenManager.Get() returns an error %v, want to succeed", err)
		}
		if got != wantToken {
			t.Errorf("TokenManager.Get() = %v, want %v", got, wantToken)
		}
		t.Run("Delete token", func(t *testing.T) {
			tm.Delete("host", "owner")
			if _, err := tm.Get("host", "owner"); err != testtarget.ErrNoHost {
				t.Errorf("TokenManager.Get() returns an error %v, want %v", err, testtarget.ErrNoHost)
			}
			t.Run("Deleted default host / owner should be empty", func(t *testing.T) {
				if gotHost, gotToken := tm.GetDefaultKey(); gotHost != github.DefaultHost || gotToken != "" {
					t.Errorf("TokenManager.GetDefaultKey() = %v, %v, want %v, %v", gotHost, gotToken, "", "")
				}
			})
		})
	})

	t.Run("Set host and owner first time, they should be default", func(t *testing.T) {
		var tm testtarget.TokenManager
		tm.Set("host1", "owner1-1", wantToken)
		if gotHost, gotToken := tm.GetDefaultKey(); gotHost != "host1" || gotToken != "owner1-1" {
			t.Errorf("TokenManager.GetDefaultKey() = %v, %v, want %v, %v", gotHost, gotToken, "host1", "owner1-1")
		}
		t.Run("Set host and owner second time, default should not be changed", func(t *testing.T) {
			tm.Set("host1", "owner1-2", wantToken)
			if gotHost, gotToken := tm.GetDefaultKey(); gotHost != "host1" || gotToken != "owner1-1" {
				t.Errorf("TokenManager.GetDefaultKey() = %v, %v, want %v, %v", gotHost, gotToken, "host1", "owner1-1")
			}
		})
		t.Run("Set another host and owner, the owner should be default", func(t *testing.T) {
			tm.Set("host2", "owner2-1", wantToken)
			if gotHost, gotToken := tm.GetDefaultKey(); gotHost != "host1" || gotToken != "owner1-1" {
				t.Errorf("TokenManager.GetDefaultKey() = %v, %v, want %v, %v", gotHost, gotToken, "host1", "owner1-1")
			}
			host, _ := tm.Hosts.TryGet("host2")
			if got := host.DefaultOwner; got != "owner2-1" {
				t.Errorf("TokenManager.Hosts.Get() = %v, want %v", got, "owner2-1")
			}
		})
	})

	t.Run("Set default host and owner", func(t *testing.T) {
		var tm testtarget.TokenManager
		tm.Set("host1", "owner1-1", wantToken)
		tm.Set("host1", "owner1-2", wantToken)
		tm.Set("host2", "owner2-1", wantToken)
		tm.Set("host2", "owner2-2", wantToken)
		if err := tm.SetDefaultOwner("host1", "owner1-2"); err != nil {
			t.Errorf("TokenManager.SetDefaultOwner() = %v, want %v", err, nil)
		}
		if gotHost, gotToken := tm.GetDefaultKey(); gotHost != "host1" || gotToken != "owner1-2" {
			t.Errorf("TokenManager.GetDefaultKey() = %v, %v, want %v, %v", gotHost, gotToken, "host1", "owner1-2")
		}
		if err := tm.SetDefaultOwner("host2", "owner2-2"); err != nil {
			t.Errorf("TokenManager.SetDefaultOwner() = %v, want %v", err, nil)
		}
		if gotHost, gotToken := tm.GetDefaultKey(); gotHost != "host1" || gotToken != "owner1-2" {
			t.Errorf("TokenManager.GetDefaultKey() = %v, %v, want %v, %v", gotHost, gotToken, "host1", "owner1-2")
		}
		if err := tm.SetDefaultHost("host2"); err != nil {
			t.Errorf("TokenManager.SetDefaultHost() = %v, want %v", err, nil)
		}
		if gotHost, gotToken := tm.GetDefaultKey(); gotHost != "host2" || gotToken != "owner2-2" {
			t.Errorf("TokenManager.GetDefaultKey() = %v, %v, want %v, %v", gotHost, gotToken, "host2", "owner2-2")
		}
	})
}
