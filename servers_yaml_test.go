package gogh_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	testtarget "github.com/kyoh86/gogh/v2"
)

func TestServersYAML(t *testing.T) {
	const (
		user1  = "kyoh86"
		user2  = "anonymous"
		host1  = "example.com" // host a not default
		host2  = "kyoh86.dev"  // host a not default
		token1 = "1111111111111111111111111111111111111111"
		token2 = "2222222222222222222222222222222222222222"
	)
	un := func(t *testing.T, file string) (servers testtarget.Servers, retErr error) {
		t.Helper()
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := yaml.Unmarshal(buf, &servers); err != nil {
			retErr = err
		}
		return
	}
	t.Run("UnmarshalInvalidYAML", func(t *testing.T) {
		const (
			invalidHost = "invalid host"
			invalidUser = "invalid user"
		)
		for _, testcase := range []struct {
			title  string
			input  string
			expect error
		}{
			{
				title:  "invalid-yaml",
				input:  "NaN",
				expect: errors.New("String node found where MapNode is expected"),
			},
			{
				title:  "invalid-value",
				input:  fmt.Sprintf(`{"%s":1}`, invalidHost),
				expect: testtarget.ErrInvalidHost("invalid value: 1"),
			},
			{
				title:  "invalid-host",
				input:  fmt.Sprintf(`{"%s":{"user":"%s","token":"%s"}}`, invalidHost, user1, token1),
				expect: testtarget.ErrInvalidHost("invalid host: " + invalidHost),
			},
			{
				title:  "invalid-host-type",
				input:  fmt.Sprintf(`{%d:{"user":"%s","token":"%s"}}`, 1, user1, token1),
				expect: testtarget.ErrInvalidHost("invalid host: 1"),
			},
			{
				title:  "invalid-user",
				input:  fmt.Sprintf(`{"%s":{"user":"%s","token":"%s"}}`, host1, invalidUser, token1),
				expect: testtarget.ErrInvalidUser("invalid user: " + invalidUser),
			},
			{
				title:  "invalid-user-type",
				input:  fmt.Sprintf(`{"%s":{"user":1,"token":"%s"}}`, host1, token1),
				expect: testtarget.ErrInvalidUser("invalid user: 1"),
			},
			{
				title:  "invalid-token-type",
				input:  fmt.Sprintf(`{"%s":{"user":"%s","token":1}}`, host1, user1),
				expect: testtarget.ErrInvalidUser("invalid token: 1"),
			},
		} {
			t.Run(testcase.title, func(t *testing.T) {
				var s testtarget.Servers
				actual := yaml.Unmarshal([]byte(testcase.input), &s)
				if actual == nil {
					t.Fatal("expect error, but actual nil is gotten")
				}
				if testcase.expect.Error() != actual.Error() {
					t.Errorf("expect error %s, but actual %s", testcase.expect, actual)
				}
			})
		}
	})
	t.Run("UnmarshalSingleYAML", func(t *testing.T) {
		servers, err := un(t, "testdata/servers_single.yaml")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		def, err := servers.Default()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def.Host() != host1 {
			t.Errorf("expect host %q, actual: %q", host1, def.Host())
		}
		if def.User() != user1 {
			t.Errorf("expect user %q, actual: %q", user1, def.User())
		}
		if def.Token() != token1 {
			t.Errorf("expect token %q, actual: %q", token1, def.Token())
		}
	})
	t.Run("UnmarshalMultipleYAML", func(t *testing.T) {
		servers, err := un(t, "testdata/servers_multiple.yaml")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		def, err := servers.Default()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def.Host() != host1 {
			t.Errorf("expect host %q, actual: %q", host1, def.Host())
		}
		if def.User() != user1 {
			t.Errorf("expect user %q, actual: %q", user1, def.User())
		}
		if def.Token() != token1 {
			t.Errorf("expect token %q, actual: %q", token1, def.Token())
		}
		second, err := servers.Find(host2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if second.Host() != host2 {
			t.Errorf("expect host %q, actual: %q", host2, second.Host())
		}
		if second.User() != user2 {
			t.Errorf("expect user %q, actual: %q", user2, second.User())
		}
		if second.Token() != token2 {
			t.Errorf("expect token %q, actual: %q", token2, second.Token())
		}
	})

	t.Run("UnmarshalEmptyYAML", func(t *testing.T) {
		servers, err := un(t, "testdata/servers_empty.yaml")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if _, err := servers.Default(); !errors.Is(err, testtarget.ErrNoServer) {
			t.Errorf("expect error: %v, acutal: %v", testtarget.ErrNoServer, err)
		}
	})

	t.Run("Marshal", func(t *testing.T) {
		var servers testtarget.Servers
		data := `
example.com:
  user: kyoh86
  token: "1111111111111111111111111111111111111111"
kyoh86.dev:
  user: anonymous
  token: "2222222222222222222222222222222222222222"`
		if err := yaml.Unmarshal([]byte(data), &servers); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		def, err := servers.Default()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if def.Host() != host1 {
			t.Errorf("expect host %q, actual: %q", host1, def.Host())
		}
		if def.User() != user1 {
			t.Errorf("expect user %q, actual: %q", user1, def.User())
		}
		if def.Token() != token1 {
			t.Errorf("expect token %q, actual: %q", token1, def.Token())
		}
		second, err := servers.Find(host2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if second.Host() != host2 {
			t.Errorf("expect host %q, actual: %q", host2, second.Host())
		}
		if second.User() != user2 {
			t.Errorf("expect user %q, actual: %q", user2, second.User())
		}
		if second.Token() != token2 {
			t.Errorf("expect token %q, actual: %q", token2, second.Token())
		}

		var buffer bytes.Buffer
		encoder := yaml.NewEncoder(&buffer, yaml.Indent(2))
		if err := encoder.Encode(&servers); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if strings.TrimSpace(data) != strings.TrimSpace(buffer.String()) {
			t.Errorf("expect marshalled data: %v, actual: %v", data, buffer.String())
		}
	})
	t.Run("MarshalEmpty", func(t *testing.T) {
		var empty testtarget.Servers
		var buffer bytes.Buffer
		encoder := yaml.NewEncoder(&buffer, yaml.Indent(2))
		if err := encoder.Encode(&empty); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if strings.TrimSpace(buffer.String()) != "null" {
			t.Errorf("expect marshalled data: %v, actual: %v", "null", buffer.String())
		}
	})
}
