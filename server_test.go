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

	t.Run("YAML", func(t *testing.T) {
		buf, err := yaml.Marshal(s)
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
		buf, err := json.Marshal(s)
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

	t.Run("UnmarshalInvalidJSON", func(t *testing.T) {
		for _, testcase := range []struct {
			title  string
			input  string
			expect error
		}{
			{
				title:  "invalid-json",
				input:  "{}}",
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

	t.Run("UnmarshalInvalidYAML", func(t *testing.T) {
		for _, testcase := range []struct {
			title  string
			input  string
			expect error
		}{
			{
				title:  "invalid-json",
				input:  "NaN",
				expect: errors.New("String node found where MapNode is expected"),
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
				actual := yaml.Unmarshal([]byte(testcase.input), &s)
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
