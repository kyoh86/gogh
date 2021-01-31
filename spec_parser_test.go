package gogh_test

import (
	"reflect"
	"testing"

	testtarget "github.com/kyoh86/gogh/v2"
)

func TestSpecParser(t *testing.T) {
	const (
		user1 = "kyoh86"
		user2 = "anonymous"
		host1 = "example.com" // host a not default
		host2 = "kyoh86.dev"  // host a not default
		name  = "gogh"
	)
	t.Run("Empty", func(t *testing.T) {
		var parser testtarget.SpecParser
		for _, testcase := range []struct {
			title  string
			source string
		}{{
			title:  "valid-name",
			source: name,
		}, {
			title:  "default-user,valid-name",
			source: user1 + "/" + name,
		}} {
			t.Run(testcase.title, func(t *testing.T) {
				if _, _, err := parser.Parse(testcase.source); err == nil {
					t.Fatal("expect error, but nil")
				}
			})
		}
	})

	t.Run("DefaultHost", func(t *testing.T) {
		server, err := testtarget.NewServer(user1, "token")
		if err != nil {
			t.Fatalf("failed to create new server for %s: %v", user1, err)
		}
		parser := testtarget.NewSpecParser(server)
		t.Run("ValidInput", func(t *testing.T) {
			for _, testcase := range []struct {
				title  string
				source string

				expectHost string
				expectUser string
				expectName string

				expectServerUser string
				expectToken      string
			}{{
				title:            "valid-name",
				source:           name,
				expectHost:       testtarget.DefaultHost,
				expectUser:       user1,
				expectName:       name,
				expectServerUser: user1,
				expectToken:      "token",
			}, {
				title:            "default-user,valid-name",
				source:           user1 + "/" + name,
				expectHost:       testtarget.DefaultHost,
				expectUser:       user1,
				expectName:       name,
				expectServerUser: user1,
				expectToken:      "token",
			}, {
				title:            "default-host,default-user,valid-name",
				source:           testtarget.DefaultHost + "/" + user1 + "/" + name,
				expectHost:       testtarget.DefaultHost,
				expectUser:       user1,
				expectName:       name,
				expectServerUser: user1,
				expectToken:      "token",
			}, {
				title:            "valid-host,valid-user,valid-name",
				source:           host1 + "/" + user2 + "/" + name,
				expectHost:       host1,
				expectUser:       user2,
				expectName:       name,
				expectServerUser: user2,
				expectToken:      "",
			}} {
				t.Run(testcase.title, func(t *testing.T) {
					spec, server, err := parser.Parse(testcase.source)
					if err != nil {
						t.Fatalf("failed to parse %q: %s", testcase.source, err)
					}
					if testcase.expectHost != spec.Host() {
						t.Errorf("expect host %q but %q gotten", testcase.expectHost, spec.Host())
					}
					if testcase.expectUser != spec.User() {
						t.Errorf("expect user %q but %q gotten", testcase.expectUser, spec.User())
					}
					if testcase.expectName != spec.Name() {
						t.Errorf("expect name %q but %q gotten", testcase.expectName, spec.Name())
					}
					if testcase.expectHost != server.Host() {
						t.Errorf("expect host %q but %q gotten", testcase.expectHost, server.Host())
					}
					if testcase.expectServerUser != server.User() {
						t.Errorf("expect user %q but %q gotten", testcase.expectServerUser, server.User())
					}
					if testcase.expectToken != server.Token() {
						t.Errorf("expect token %q but %q gotten", testcase.expectToken, server.Token())
					}
				})
			}
		})

		t.Run("InvalidInput", func(t *testing.T) {
			for _, testcase := range []struct {
				title  string
				input  string
				expect error
			}{
				{
					title:  "empty",
					input:  "",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "empty-user,empty-name",
					input:  "/",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "empty-user,valid-name",
					input:  "/" + name,
					expect: testtarget.ErrEmptyUser,
				},
				{
					title:  "valid-user,dot",
					input:  user1 + "/.",
					expect: testtarget.ErrInvalidName("'.' is reserved name"),
				},
				{
					title:  "valid-user,dotdot",
					input:  user1 + "/..",
					expect: testtarget.ErrInvalidName("'..' is reserved name"),
				},
				{
					title:  "invalid-user,valid-name",
					input:  "space in the user/" + name,
					expect: testtarget.ErrInvalidUser("invalid user: space in the user"),
				},
				{
					title:  "valid-user,empty-name",
					input:  user1 + "/",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "valid-user,invalid-name",
					input:  user1 + "/space in the name",
					expect: testtarget.ErrInvalidName("invalid name: space in the name"),
				},

				{
					title:  "empty-host,valid-user,valid-name",
					input:  "/" + user1 + "/" + name,
					expect: testtarget.ErrEmptyHost,
				},
				{
					title:  "invalid-host,valid-user,valid-name",
					input:  "space in the host/" + user1 + "/" + name,
					expect: testtarget.ErrInvalidHost("invalid host: space in the host"),
				},
				{
					title:  "valid-host,empty-user,valid-name",
					input:  host1 + "//" + name,
					expect: testtarget.ErrEmptyUser,
				},
				{
					title:  "valid-host,invalid-user,valid-name",
					input:  host1 + "/space in the user/" + name,
					expect: testtarget.ErrInvalidUser("invalid user: space in the user"),
				},
				{
					title:  "valid-host,valid-user,empty-name",
					input:  host1 + "/" + user1 + "/",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "valid-host,valid-user,invalid-name",
					input:  host1 + "/" + user1 + "/space in the name",
					expect: testtarget.ErrInvalidName("invalid name: space in the name"),
				},
				{
					title:  "valid-host,empty-user,empty-name",
					input:  host1 + "//",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "empty-host,valid-user,empty-name",
					input:  "/" + user1 + "/",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "empty-host,empty-user,valid-name",
					input:  "//" + name,
					expect: testtarget.ErrEmptyUser,
				},
				{
					title:  "empty-host,empty-user,empty-name",
					input:  "//",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "unnecessary-following-slash",
					input:  host1 + "/" + user1 + "/" + name + "/",
					expect: testtarget.ErrTooManySlashes,
				},
				{
					title:  "unnecessary-heading-slash",
					input:  "/" + host1 + "/" + user1 + "/" + name + "/",
					expect: testtarget.ErrTooManySlashes,
				},
			} {
				t.Run(testcase.title, func(t *testing.T) {
					spec, _, err := parser.Parse(testcase.input)
					if err == nil {
						t.Fatalf("expect failure to parse %q but parsed to %+v", testcase.input, spec)
					}
					if reflect.TypeOf(testcase.expect) != reflect.TypeOf(err) {
						t.Fatalf("expect error %t to parse %q but %t gotten", testcase.expect, testcase.input, err)
					}
					if testcase.expect.Error() != err.Error() {
						t.Fatalf("expect error value %q to parse %q but %q gotten", testcase.expect, testcase.input, err)
					}
				})
			}
		})
	})

	t.Run("WithHost", func(t *testing.T) {
		server, err := testtarget.NewServerFor(host1, user1, "token")
		if err != nil {
			t.Fatalf("failed to create new server for %s: %v", user1, err)
		}
		parser := testtarget.NewSpecParser(server)
		t.Run("ValidInput", func(t *testing.T) {
			for _, testcase := range []struct {
				title  string
				source string

				expectHost string
				expectUser string
				expectName string

				expectServerUser string
				expectToken      string
			}{{
				title:            "valid-name",
				source:           name,
				expectHost:       host1,
				expectUser:       user1,
				expectName:       name,
				expectServerUser: user1,
				expectToken:      "token",
			}, {
				title:            "default-user,valid-name",
				source:           user1 + "/" + name,
				expectHost:       host1,
				expectUser:       user1,
				expectName:       name,
				expectServerUser: user1,
				expectToken:      "token",
			}, {
				title:            "default-host,default-user,valid-name",
				source:           host1 + "/" + user1 + "/" + name,
				expectHost:       host1,
				expectUser:       user1,
				expectName:       name,
				expectServerUser: user1,
				expectToken:      "token",
			}, {
				title:            "valid-host,valid-user,valid-name",
				source:           host2 + "/" + user2 + "/" + name,
				expectHost:       host2,
				expectUser:       user2,
				expectServerUser: user2,
				expectName:       name,
			}} {
				t.Run(testcase.title, func(t *testing.T) {
					spec, server, err := parser.Parse(testcase.source)
					if err != nil {
						t.Fatalf("failed to parse %q: %s", testcase.source, err)
					}
					if testcase.expectHost != spec.Host() {
						t.Errorf("expect host %q but %q gotten", testcase.expectHost, spec.Host())
					}
					if testcase.expectUser != spec.User() {
						t.Errorf("expect user %q but %q gotten", testcase.expectUser, spec.User())
					}
					if testcase.expectName != spec.Name() {
						t.Errorf("expect name %q but %q gotten", testcase.expectName, spec.Name())
					}
					if testcase.expectHost != server.Host() {
						t.Errorf("expect host %q but %q gotten", testcase.expectHost, server.Host())
					}
					if testcase.expectServerUser != server.User() {
						t.Errorf("expect user %q but %q gotten", testcase.expectServerUser, server.User())
					}
					if testcase.expectToken != server.Token() {
						t.Errorf("expect token %q but %q gotten", testcase.expectToken, server.Token())
					}
				})
			}
		})
	})

	t.Run("WithMultipeServers", func(t *testing.T) {
		// (default) github.com/kyoh86
		server1, err := testtarget.NewServerFor(testtarget.DefaultHost, user1, "token1")
		if err != nil {
			t.Fatalf("failed to create new server for %s@%s: %q", user1, testtarget.DefaultHost, err)
		}
		// example.com/anonymous
		server2, err := testtarget.NewServerFor(host1, user2, "token2")
		if err != nil {
			t.Fatalf("failed to create new server for %s@%s: %q", user2, host1, err)
		}
		// kyoh86.dev/anonymous
		server3, err := testtarget.NewServerFor(testtarget.DefaultHost, user2, "token3")
		if err != nil {
			t.Fatalf("failed to create new server for %s@%s: %q", user2, testtarget.DefaultHost, err)
		}

		parser := testtarget.NewSpecParser(server1, server2, server3)

		for _, testcase := range []struct {
			title  string
			source string

			expectHost string
			expectUser string
			expectName string

			expectServerUser string
			expectToken      string
		}{{
			title:            "valid-name",
			source:           name,
			expectHost:       testtarget.DefaultHost,
			expectUser:       user1,
			expectName:       name,
			expectServerUser: user1,
			expectToken:      "token1",
		}, {
			title:            "valid-name,valid-user",
			source:           user2 + "/" + name,
			expectHost:       testtarget.DefaultHost,
			expectUser:       user2,
			expectServerUser: user1,
			expectName:       name,
			expectToken:      "token1",
		}, {
			title:            "full-name",
			source:           testtarget.DefaultHost + "/" + user2 + "/" + name,
			expectHost:       testtarget.DefaultHost,
			expectUser:       user2,
			expectName:       name,
			expectServerUser: user1,
			expectToken:      "token1",
		}, {
			title:            "other-host",
			source:           host1 + "/" + user2 + "/" + name,
			expectHost:       host1,
			expectUser:       user2,
			expectName:       name,
			expectServerUser: user2,
			expectToken:      "token2",
		}, {
			title:            "not-matched",
			source:           host2 + "/" + user2 + "/" + name,
			expectHost:       host2,
			expectUser:       user2,
			expectServerUser: user2,
			expectName:       name,
			expectToken:      "",
		}} {
			t.Run(testcase.title, func(t *testing.T) {
				spec, server, err := parser.Parse(testcase.source)
				if err != nil {
					t.Fatalf("failed to parse %q: %s", testcase.source, err)
				}
				if testcase.expectHost != spec.Host() {
					t.Errorf("expect host %q but %q gotten", testcase.expectHost, spec.Host())
				}
				if testcase.expectUser != spec.User() {
					t.Errorf("expect user %q but %q gotten", testcase.expectUser, spec.User())
				}
				if testcase.expectName != spec.Name() {
					t.Errorf("expect name %q but %q gotten", testcase.expectName, spec.Name())
				}
				if testcase.expectHost != server.Host() {
					t.Errorf("expect host %q but %q gotten", testcase.expectHost, server.Host())
				}
				if testcase.expectServerUser != server.User() {
					t.Errorf("expect user %q but %q gotten", testcase.expectServerUser, server.User())
				}
				if testcase.expectToken != server.Token() {
					t.Errorf("expect token %q but %q gotten", testcase.expectToken, server.Token())
				}
			})
		}
	})
}
