package gogh_test

import (
	"testing"

	testtarget "github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/github"
)

func TestTokenManager(t *testing.T) {
	t.Run("Empty token manager always returns empty", func(t *testing.T) {
		var tm testtarget.TokenManager
		if got := tm.Get("host", "owner"); got != "" {
			t.Errorf("TokenManager.Get() = %v, want %v", got, "")
		}
		if gotHost, gotToken := tm.GetDefaultKey(); gotHost != github.DefaultHost || gotToken != "" {
			t.Errorf("TokenManager.GetDefaultKey() = %v, %v, want %v, %v", gotHost, gotToken, "", "")
		}
	})

	t.Run("Set and get token", func(t *testing.T) {
		var tm testtarget.TokenManager
		tm.Set("host", "owner", "token")
		if got := tm.Get("host", "owner"); got != "token" {
			t.Errorf("TokenManager.Get() = %v, want %v", got, "token")
		}
		t.Run("Delete token", func(t *testing.T) {
			tm.Delete("host", "owner")
			if got := tm.Get("host", "owner"); got != "" {
				t.Errorf("TokenManager.Get() = %v, want %v", got, "")
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
		tm.Set("host1", "owner1-1", "token")
		if gotHost, gotToken := tm.GetDefaultKey(); gotHost != "host1" || gotToken != "owner1-1" {
			t.Errorf("TokenManager.GetDefaultKey() = %v, %v, want %v, %v", gotHost, gotToken, "host1", "owner1-1")
		}
		t.Run("Set host and owner second time, default should not be changed", func(t *testing.T) {
			tm.Set("host1", "owner1-2", "token")
			if gotHost, gotToken := tm.GetDefaultKey(); gotHost != "host1" || gotToken != "owner1-1" {
				t.Errorf("TokenManager.GetDefaultKey() = %v, %v, want %v, %v", gotHost, gotToken, "host1", "owner1-1")
			}
		})
		t.Run("Set another host and owner, the owner should be default", func(t *testing.T) {
			tm.Set("host2", "owner2-1", "token")
			if gotHost, gotToken := tm.GetDefaultKey(); gotHost != "host1" || gotToken != "owner1-1" {
				t.Errorf("TokenManager.GetDefaultKey() = %v, %v, want %v, %v", gotHost, gotToken, "host1", "owner1-1")
			}
			if got := tm.Hosts.Get("host2").DefaultOwner; got != "owner2-1" {
				t.Errorf("TokenManager.Hosts.Get() = %v, want %v", got, "owner2-1")
			}
		})
	})

	t.Run("Set default host and owner", func(t *testing.T) {
		var tm testtarget.TokenManager
		tm.Set("host1", "owner1-1", "token")
		tm.Set("host1", "owner1-2", "token")
		tm.Set("host2", "owner2-1", "token")
		tm.Set("host2", "owner2-2", "token")
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
