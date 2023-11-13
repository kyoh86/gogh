package migrate_test

import (
	"fmt"
	"testing"

	testtarget "github.com/kyoh86/gogh/v2/cmd/gogh/migrate/internal"
)

func TestServer(t *testing.T) {
	host := "example.com"
	user := "kyoh86"
	token := "xxxxxxxxxxx"
	s, err := testtarget.NewServerFor(host, user, token)
	if err != nil {
		t.Fatal("failed to create new server")
	}
	if host != s.Host() {
		t.Fatalf("expect host %q but %q", host, s.Host())
	}
	if user != s.User() {
		t.Fatalf("expect user %q but %q", user, s.User())
	}
	if token != s.Token() {
		t.Errorf("expect token %q but %q", token, s.Token())
	}
	wantStr := "kyoh86@example.com"
	if wantStr != s.String() {
		t.Errorf("expect the string notation %s, but %q", wantStr, s.String())
	}
	wantFormatted := "server: kyoh86@example.com"
	if wantFormatted != fmt.Sprintf("server: %s", s) {
		t.Errorf("expect hide the token in the string notation, but %q", s.String())
	}
	invalidHost := "invalid host"
	invalidUser := "invalid user"

	t.Run("NewServer", func(t *testing.T) {
		if _, err := testtarget.NewServer(invalidUser, token); err == nil {
			t.Error("expect failure to create new server with invalid host, but not")
		}
		s, err := testtarget.NewServer(user, token)
		if err != nil {
			t.Fatal("failed to create new server")
		}
		if user != s.User() {
			t.Errorf("expect user %q but %q", user, s.User())
		}
		if token != s.Token() {
			t.Errorf("expect token %q but %q", token, s.Token())
		}
	})

	t.Run("NewServerFor", func(t *testing.T) {
		if _, err := testtarget.NewServerFor(invalidHost, user, token); err == nil {
			t.Error("expect failure to create new server with invalid host, but not")
		}
		if _, err := testtarget.NewServerFor(host, invalidUser, token); err == nil {
			t.Error("expect failure to create new server with invalid user, but not")
		}
		s, err := testtarget.NewServerFor(host, user, token)
		if err != nil {
			t.Fatal("failed to create new server")
		}
		if host != s.Host() {
			t.Errorf("expect host %q but %q", host, s.Host())
		}
		if user != s.User() {
			t.Errorf("expect user %q but %q", user, s.User())
		}
		if token != s.Token() {
			t.Errorf("expect token %q but %q", token, s.Token())
		}
	})
}
