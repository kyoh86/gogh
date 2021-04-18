package gogh_test

import (
	"reflect"
	"testing"

	testtarget "github.com/kyoh86/gogh/v2"
)

func TestSpecParser(t *testing.T) {
	const (
		owner1 = "kyoh86"
		owner2 = "anonymous"
		host1  = "example.com" // host a not default
		host2  = "kyoh86.dev"  // host a not default
		name   = "gogh"
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
			title:  "default-owner,valid-name",
			source: owner1 + "/" + name,
		}} {
			t.Run(testcase.title, func(t *testing.T) {
				if _, _, err := parser.Parse(testcase.source); err == nil {
					t.Fatal("expect error, but nil")
				}
			})
		}
	})

	t.Run("DefaultHost", func(t *testing.T) {
		server, err := testtarget.NewServer(owner1, "token")
		if err != nil {
			t.Fatalf("failed to create new server for %s: %v", owner1, err)
		}
		parser := testtarget.NewSpecParser(testtarget.NewServers(server))
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
				expectUser:       owner1,
				expectName:       name,
				expectServerUser: owner1,
				expectToken:      "token",
			}, {
				title:            "default-owner,valid-name",
				source:           owner1 + "/" + name,
				expectHost:       testtarget.DefaultHost,
				expectUser:       owner1,
				expectName:       name,
				expectServerUser: owner1,
				expectToken:      "token",
			}, {
				title:            "default-host,default-owner,valid-name",
				source:           testtarget.DefaultHost + "/" + owner1 + "/" + name,
				expectHost:       testtarget.DefaultHost,
				expectUser:       owner1,
				expectName:       name,
				expectServerUser: owner1,
				expectToken:      "token",
			}, {
				title:            "valid-host,valid-owner,valid-name",
				source:           host1 + "/" + owner2 + "/" + name,
				expectHost:       host1,
				expectUser:       owner2,
				expectName:       name,
				expectServerUser: owner2,
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
					if testcase.expectUser != spec.Owner() {
						t.Errorf("expect owner %q but %q gotten", testcase.expectUser, spec.Owner())
					}
					if testcase.expectName != spec.Name() {
						t.Errorf("expect name %q but %q gotten", testcase.expectName, spec.Name())
					}
					if testcase.expectHost != server.Host() {
						t.Errorf("expect host %q but %q gotten", testcase.expectHost, server.Host())
					}
					if testcase.expectServerUser != server.User() {
						t.Errorf("expect owner %q but %q gotten", testcase.expectServerUser, server.User())
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
					title:  "empty-owner,empty-name",
					input:  "/",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "empty-owner,valid-name",
					input:  "/" + name,
					expect: testtarget.ErrEmptyOwner,
				},
				{
					title:  "valid-owner,dot",
					input:  owner1 + "/.",
					expect: testtarget.ErrInvalidName("'.' is reserved name"),
				},
				{
					title:  "valid-owner,dotdot",
					input:  owner1 + "/..",
					expect: testtarget.ErrInvalidName("'..' is reserved name"),
				},
				{
					title:  "invalid-owner,valid-name",
					input:  "space in the owner/" + name,
					expect: testtarget.ErrInvalidOwner("invalid owner: space in the owner"),
				},
				{
					title:  "valid-owner,empty-name",
					input:  owner1 + "/",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "valid-owner,invalid-name",
					input:  owner1 + "/space in the name",
					expect: testtarget.ErrInvalidName("invalid name: space in the name"),
				},

				{
					title:  "empty-host,valid-owner,valid-name",
					input:  "/" + owner1 + "/" + name,
					expect: testtarget.ErrEmptyHost,
				},
				{
					title:  "invalid-host,valid-owner,valid-name",
					input:  "space in the host/" + owner1 + "/" + name,
					expect: testtarget.ErrInvalidHost("invalid host: space in the host"),
				},
				{
					title:  "valid-host,empty-owner,valid-name",
					input:  host1 + "//" + name,
					expect: testtarget.ErrEmptyOwner,
				},
				{
					title:  "valid-host,invalid-owner,valid-name",
					input:  host1 + "/space in the owner/" + name,
					expect: testtarget.ErrInvalidOwner("invalid owner: space in the owner"),
				},
				{
					title:  "valid-host,valid-owner,empty-name",
					input:  host1 + "/" + owner1 + "/",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "valid-host,valid-owner,invalid-name",
					input:  host1 + "/" + owner1 + "/space in the name",
					expect: testtarget.ErrInvalidName("invalid name: space in the name"),
				},
				{
					title:  "valid-host,empty-owner,empty-name",
					input:  host1 + "//",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "empty-host,valid-owner,empty-name",
					input:  "/" + owner1 + "/",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "empty-host,empty-owner,valid-name",
					input:  "//" + name,
					expect: testtarget.ErrEmptyOwner,
				},
				{
					title:  "empty-host,empty-owner,empty-name",
					input:  "//",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "unnecessary-following-slash",
					input:  host1 + "/" + owner1 + "/" + name + "/",
					expect: testtarget.ErrTooManySlashes,
				},
				{
					title:  "unnecessary-heading-slash",
					input:  "/" + host1 + "/" + owner1 + "/" + name + "/",
					expect: testtarget.ErrTooManySlashes,
				},
			} {
				t.Run(testcase.title, func(t *testing.T) {
					spec, _, _, err := parser.ParseWithAlias(testcase.input)
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
		server, err := testtarget.NewServerFor(host1, owner1, "token")
		if err != nil {
			t.Fatalf("failed to create new server for %s: %v", owner1, err)
		}
		parser := testtarget.NewSpecParser(testtarget.NewServers(server))
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
				expectUser:       owner1,
				expectName:       name,
				expectServerUser: owner1,
				expectToken:      "token",
			}, {
				title:            "default-owner,valid-name",
				source:           owner1 + "/" + name,
				expectHost:       host1,
				expectUser:       owner1,
				expectName:       name,
				expectServerUser: owner1,
				expectToken:      "token",
			}, {
				title:            "default-host,default-owner,valid-name",
				source:           host1 + "/" + owner1 + "/" + name,
				expectHost:       host1,
				expectUser:       owner1,
				expectName:       name,
				expectServerUser: owner1,
				expectToken:      "token",
			}, {
				title:            "valid-host,valid-owner,valid-name",
				source:           host2 + "/" + owner2 + "/" + name,
				expectHost:       host2,
				expectUser:       owner2,
				expectServerUser: owner2,
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
					if testcase.expectUser != spec.Owner() {
						t.Errorf("expect owner %q but %q gotten", testcase.expectUser, spec.Owner())
					}
					if testcase.expectName != spec.Name() {
						t.Errorf("expect name %q but %q gotten", testcase.expectName, spec.Name())
					}
					if testcase.expectHost != server.Host() {
						t.Errorf("expect host %q but %q gotten", testcase.expectHost, server.Host())
					}
					if testcase.expectServerUser != server.User() {
						t.Errorf("expect owner %q but %q gotten", testcase.expectServerUser, server.User())
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
		server1, err := testtarget.NewServerFor(testtarget.DefaultHost, owner1, "token1")
		if err != nil {
			t.Fatalf("failed to create new server for %s@%s: %q", owner1, testtarget.DefaultHost, err)
		}
		// example.com/anonymous
		server2, err := testtarget.NewServerFor(host1, owner2, "token2")
		if err != nil {
			t.Fatalf("failed to create new server for %s@%s: %q", owner2, host1, err)
		}
		// kyoh86.dev/anonymous
		server3, err := testtarget.NewServerFor(testtarget.DefaultHost, owner2, "token3")
		if err != nil {
			t.Fatalf("failed to create new server for %s@%s: %q", owner2, testtarget.DefaultHost, err)
		}

		parser := testtarget.NewSpecParser(testtarget.NewServers(server1, server2, server3))

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
			expectUser:       owner1,
			expectName:       name,
			expectServerUser: owner1,
			expectToken:      "token1",
		}, {
			title:            "valid-name,valid-owner",
			source:           owner2 + "/" + name,
			expectHost:       testtarget.DefaultHost,
			expectUser:       owner2,
			expectServerUser: owner1,
			expectName:       name,
			expectToken:      "token1",
		}, {
			title:            "full-name",
			source:           testtarget.DefaultHost + "/" + owner2 + "/" + name,
			expectHost:       testtarget.DefaultHost,
			expectUser:       owner2,
			expectName:       name,
			expectServerUser: owner1,
			expectToken:      "token1",
		}, {
			title:            "other-host",
			source:           host1 + "/" + owner2 + "/" + name,
			expectHost:       host1,
			expectUser:       owner2,
			expectName:       name,
			expectServerUser: owner2,
			expectToken:      "token2",
		}, {
			title:            "not-matched",
			source:           host2 + "/" + owner2 + "/" + name,
			expectHost:       host2,
			expectUser:       owner2,
			expectServerUser: owner2,
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
				if testcase.expectUser != spec.Owner() {
					t.Errorf("expect owner %q but %q gotten", testcase.expectUser, spec.Owner())
				}
				if testcase.expectName != spec.Name() {
					t.Errorf("expect name %q but %q gotten", testcase.expectName, spec.Name())
				}
				if testcase.expectHost != server.Host() {
					t.Errorf("expect host %q but %q gotten", testcase.expectHost, server.Host())
				}
				if testcase.expectServerUser != server.User() {
					t.Errorf("expect owner %q but %q gotten", testcase.expectServerUser, server.User())
				}
				if testcase.expectToken != server.Token() {
					t.Errorf("expect token %q but %q gotten", testcase.expectToken, server.Token())
				}
			})
		}
	})

	t.Run("ParseWithAlias", func(t *testing.T) {
		server, err := testtarget.NewServerFor(host1, owner1, "token")
		if err != nil {
			t.Fatalf("failed to create new server for %s: %v", owner1, err)
		}
		parser := testtarget.NewSpecParser(testtarget.NewServers(server))
		t.Run("WithoutAlias", func(t *testing.T) {
			for _, testcase := range []struct {
				title  string
				source string
			}{{
				title:  "without-alias",
				source: name,
			}, {
				title:  "same-name",
				source: name + "=" + name,
			}} {
				_, alias, _, err := parser.ParseWithAlias(testcase.source)
				if err != nil {
					t.Fatalf("failed to parse %q: %s", testcase.source, err)
				}
				if alias != nil {
					t.Errorf("want alias is nil but %#v gotten", *alias)
				}
			}
		})

		t.Run("WithValidAlias", func(t *testing.T) {
			for _, testcase := range []struct {
				title  string
				source string

				wantHost string
				wantUser string
				wantName string
			}{{
				title:    "with-alias-name",
				source:   host1 + "/" + owner1 + "/" + name + "=alias",
				wantHost: host1,
				wantUser: owner1,
				wantName: "alias",
			}, {
				title:    "with-not-default-repo-with-alias-name",
				source:   host2 + "/" + owner2 + "/" + name + "=alias",
				wantHost: host2,
				wantUser: owner2,
				wantName: "alias",
			}, {
				title:    "with-alias-owner",
				source:   host1 + "/" + owner1 + "/" + name + "=" + owner2 + "/alias",
				wantHost: host1,
				wantUser: owner2,
				wantName: "alias",
			}} {
				t.Run(testcase.title, func(t *testing.T) {
					_, alias, _, err := parser.ParseWithAlias(testcase.source)
					if err != nil {
						t.Fatalf("failed to parse %q: %s", testcase.source, err)
					}
					if alias == nil {
						t.Fatal("want valid alisa but got nil")
					}
					if testcase.wantHost != alias.Host() {
						t.Errorf("want host %q but %q gotten", testcase.wantHost, alias.Host())
					}
					if testcase.wantUser != alias.Owner() {
						t.Errorf("want owner %q but %q gotten", testcase.wantUser, alias.Owner())
					}
					if testcase.wantName != alias.Name() {
						t.Errorf("want name %q but %q gotten", testcase.wantName, alias.Name())
					}
					if testcase.wantHost != alias.Host() {
						t.Errorf("want host %q but %q gotten", testcase.wantHost, alias.Host())
					}
				})
			}
		})

		t.Run("WithInvalid", func(t *testing.T) {
			for _, testcase := range []struct {
				title  string
				source string
			}{{
				title:  "invalid-name",
				source: ".=",
			}, {
				title:  "empty-alias",
				source: name + "=",
			}, {
				title:  "space-alias",
				source: name + "= ",
			}, {
				title:  "space-in-the-alias",
				source: name + "=splitted name",
			}, {
				title:  "double-alias",
				source: name + "=alias1=alias2",
			}, {
				title:  "too-many-shashes",
				source: name + "=example.com/baz/many",
			}} {
				t.Run(testcase.title, func(t *testing.T) {
					_, _, _, err := parser.ParseWithAlias(testcase.source)
					if err == nil {
						t.Fatal("want error, but got nil")
					}
					t.Log(err)
				})
			}
		})
	})
}
