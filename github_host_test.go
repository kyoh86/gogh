package gogh_test

import (
	"encoding/json"
	"testing"

	"github.com/goccy/go-yaml"
	testtarget "github.com/kyoh86/gogh/v2"
)

func TestGithubHost(t *testing.T) {
	var h testtarget.GithubHost
	if err := h.SetName("invalid host"); err == nil {
		t.Error("expect failure to set invalid host, but not")
	}

	host := "example.com"
	if err := h.SetName(host); err != nil {
		t.Fatalf("failed to set host: %s", err)
	}
	if host != h.Name() {
		t.Errorf("expect host %q but %q", host, h.Name())
	}
	user := "kyoh86"
	if err := h.SetUser(user); err != nil {
		t.Fatalf("failed to set user: %s", err)
	}
	if user != h.User() {
		t.Errorf("expect user %q but %q", user, h.User())
	}
	token := "xxxxxxxxxxx"
	if err := h.SetToken(token); err != nil {
		t.Fatalf("failed to set token: %s", err)
	}
	if token != h.Token() {
		t.Errorf("expect token %q but %q", token, h.Token())
	}

	t.Run("YAML", func(t *testing.T) {
		buf, err := yaml.Marshal(h)
		if err != nil {
			t.Fatalf("failed to marshal: %s", err)
		}
		var actual testtarget.GithubHost
		if err := yaml.Unmarshal(buf, &actual); err != nil {
			t.Fatalf("failed to unmarshal: %s", err)
		}
		if host != actual.Name() {
			t.Errorf("expect host %q but %q", host, actual.Name())
		}
		if user != actual.User() {
			t.Errorf("expect user %q but %q", user, actual.User())
		}
		if token != actual.Token() {
			t.Errorf("expect token %q but %q", token, actual.Token())
		}
	})

	t.Run("JSON", func(t *testing.T) {
		buf, err := json.Marshal(h)
		if err != nil {
			t.Fatalf("failed to marshal: %s", err)
		}
		var actual testtarget.GithubHost
		if err := json.Unmarshal(buf, &actual); err != nil {
			t.Fatalf("failed to unmarshal: %s", err)
		}
		if host != actual.Name() {
			t.Errorf("expect host %q but %q", host, actual.Name())
		}
		if user != actual.User() {
			t.Errorf("expect user %q but %q", user, actual.User())
		}
		if token != actual.Token() {
			t.Errorf("expect token %q but %q", token, actual.Token())
		}
	})

	//UNDONE: set invalid host name
	//UNDONE: unmarshal JSON with invalid host
	//UNDONE: unmarshal YAML with invalid host
}
