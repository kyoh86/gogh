package gogh_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/goccy/go-yaml"
	testtarget "github.com/kyoh86/gogh/v2"
)

func TestServer(t *testing.T) {
	var h testtarget.Server
	if err := h.SetHost("invalid host"); err == nil {
		t.Error("expect failure to set invalid host, but not")
	}

	host := "example.com"
	user := "kyoh86"
	token := "xxxxxxxxxxx"
	t.Run("SetValidValue", func(t *testing.T) {
		if err := h.SetHost(host); err != nil {
			t.Fatalf("failed to set host: %s", err)
		}
		if host != h.Host() {
			t.Errorf("expect host %q but %q", host, h.Host())
		}
		if err := h.SetUser(user); err != nil {
			t.Fatalf("failed to set user: %s", err)
		}
		if user != h.User() {
			t.Errorf("expect user %q but %q", user, h.User())
		}
		if err := h.SetToken(token); err != nil {
			t.Fatalf("failed to set token: %s", err)
		}
		if token != h.Token() {
			t.Errorf("expect token %q but %q", token, h.Token())
		}
	})

	t.Run("YAML", func(t *testing.T) {
		buf, err := yaml.Marshal(h)
		if err != nil {
			t.Fatalf("failed to marshal: %s", err)
		}
		var actual testtarget.Server
		if err := yaml.Unmarshal(buf, &actual); err != nil {
			t.Fatalf("failed to unmarshal: %s", err)
		}
		if host != actual.Host() {
			t.Errorf("expect host %q but %q", host, actual.Host())
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
		var actual testtarget.Server
		if err := json.Unmarshal(buf, &actual); err != nil {
			t.Fatalf("failed to unmarshal: %s", err)
		}
		if host != actual.Host() {
			t.Errorf("expect host %q but %q", host, actual.Host())
		}
		if user != actual.User() {
			t.Errorf("expect user %q but %q", user, actual.User())
		}
		if token != actual.Token() {
			t.Errorf("expect token %q but %q", token, actual.Token())
		}
	})

	invalidHost := "invalid host"
	invalidUser := "invalid user"
	t.Run("SetInvalidValue", func(t *testing.T) {
		if err := h.SetHost(invalidHost); err == nil {
			t.Fatalf("expect failure to set invalid host, but not")
		}
		if host != h.Host() { // expect no changing
			t.Errorf("expect host %q but %q", host, h.Host())
		}
		if err := h.SetUser(invalidUser); err == nil {
			t.Fatalf("expect failure to set invalid user, but not")
		}
		if user != h.User() { // expect no changing
			t.Errorf("expect user %q but %q", user, h.User())
		}
	})

	t.Run("UnmarshalInvalid", func(t *testing.T) {
		for _, testcase := range []struct {
			title  string
			input  string
			expect error
		}{
			{
				title:  "invalid-json",
				input:  "{}}}",
				expect: errors.New("invalid character '}' after top-level value"),
			},
			{
				title:  "empty-host",
				input:  `{"host":""}`,
				expect: testtarget.ErrEmptyHost,
			},
			{
				title:  "invalid-user",
				input:  fmt.Sprintf(`{"host":"%s","user":"%s"}`, host, invalidUser),
				expect: testtarget.ErrInvalidUser("invalid user: " + invalidUser),
			},
		} {
			t.Run(testcase.title, func(t *testing.T) {
				var s testtarget.Server
				actual := json.Unmarshal([]byte(testcase.input), &s)
				if actual == nil {
					t.Fatal("ecpect error, but actual nil is gotten")
				}
				if testcase.expect.Error() != actual.Error() {
					t.Errorf("ecpect error %s, but actual %s", testcase.expect, actual)
				}
			})
		}
	})
}
