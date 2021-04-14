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
					t.Fatal("want error, but nil")
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

				wantExplicitHost string
				wantExplicitUser string
				wantHost         string
				wantUser         string
				wantName         string

				wantServerUser string
				wantToken      string
			}{{
				title:            "valid-name",
				source:           name,
				wantExplicitHost: "",
				wantExplicitUser: "",
				wantHost:         testtarget.DefaultHost,
				wantUser:         owner1,
				wantName:         name,
				wantServerUser:   owner1,
				wantToken:        "token",
			}, {
				title:            "default-owner,valid-name",
				source:           owner1 + "/" + name,
				wantExplicitHost: "",
				wantExplicitUser: owner1,
				wantHost:         testtarget.DefaultHost,
				wantUser:         owner1,
				wantName:         name,
				wantServerUser:   owner1,
				wantToken:        "token",
			}, {
				title:            "default-host,default-owner,valid-name",
				source:           testtarget.DefaultHost + "/" + owner1 + "/" + name,
				wantExplicitHost: testtarget.DefaultHost,
				wantExplicitUser: owner1,
				wantHost:         testtarget.DefaultHost,
				wantUser:         owner1,
				wantName:         name,
				wantServerUser:   owner1,
				wantToken:        "token",
			}, {
				title:            "valid-host,valid-owner,valid-name",
				source:           host1 + "/" + owner2 + "/" + name,
				wantExplicitHost: host1,
				wantExplicitUser: owner2,
				wantHost:         host1,
				wantUser:         owner2,
				wantName:         name,
				wantServerUser:   owner2,
				wantToken:        "",
			}} {
				t.Run(testcase.title, func(t *testing.T) {
					spec, server, err := parser.Parse(testcase.source)
					if err != nil {
						t.Fatalf("failed to parse %q: %s", testcase.source, err)
					}
					if testcase.wantExplicitHost != spec.ExplicitHost() {
						t.Errorf("want explicit host %q but %q gotten", testcase.wantExplicitHost, spec.ExplicitHost())
					}
					if testcase.wantExplicitUser != spec.ExplicitOwner() {
						t.Errorf("want explicit owner %q but %q gotten", testcase.wantExplicitUser, spec.ExplicitOwner())
					}
					if testcase.wantHost != spec.Host() {
						t.Errorf("want host %q but %q gotten", testcase.wantHost, spec.Host())
					}
					if testcase.wantUser != spec.Owner() {
						t.Errorf("want owner %q but %q gotten", testcase.wantUser, spec.Owner())
					}
					if testcase.wantName != spec.Name() {
						t.Errorf("want name %q but %q gotten", testcase.wantName, spec.Name())
					}
					if testcase.wantHost != server.Host() {
						t.Errorf("want host %q but %q gotten", testcase.wantHost, server.Host())
					}
					if testcase.wantServerUser != server.User() {
						t.Errorf("want owner %q but %q gotten", testcase.wantServerUser, server.User())
					}
					if testcase.wantToken != server.Token() {
						t.Errorf("want token %q but %q gotten", testcase.wantToken, server.Token())
					}
				})
			}
		})

		t.Run("InvalidInput", func(t *testing.T) {
			for _, testcase := range []struct {
				title string
				input string
				want  error
			}{
				{
					title: "empty",
					input: "",
					want:  testtarget.ErrEmptyName,
				},
				{
					title: "empty-owner,empty-name",
					input: "/",
					want:  testtarget.ErrEmptyName,
				},
				{
					title: "empty-owner,valid-name",
					input: "/" + name,
					want:  testtarget.ErrEmptyOwner,
				},
				{
					title: "valid-owner,dot",
					input: owner1 + "/.",
					want:  testtarget.ErrInvalidName("'.' is reserved name"),
				},
				{
					title: "valid-owner,dotdot",
					input: owner1 + "/..",
					want:  testtarget.ErrInvalidName("'..' is reserved name"),
				},
				{
					title: "invalid-owner,valid-name",
					input: "space in the owner/" + name,
					want:  testtarget.ErrInvalidOwner("invalid owner: space in the owner"),
				},
				{
					title: "valid-owner,empty-name",
					input: owner1 + "/",
					want:  testtarget.ErrEmptyName,
				},
				{
					title: "valid-owner,invalid-name",
					input: owner1 + "/space in the name",
					want:  testtarget.ErrInvalidName("invalid name: space in the name"),
				},

				{
					title: "empty-host,valid-owner,valid-name",
					input: "/" + owner1 + "/" + name,
					want:  testtarget.ErrEmptyHost,
				},
				{
					title: "invalid-host,valid-owner,valid-name",
					input: "space in the host/" + owner1 + "/" + name,
					want:  testtarget.ErrInvalidHost("invalid host: space in the host"),
				},
				{
					title: "valid-host,empty-owner,valid-name",
					input: host1 + "//" + name,
					want:  testtarget.ErrEmptyOwner,
				},
				{
					title: "valid-host,invalid-owner,valid-name",
					input: host1 + "/space in the owner/" + name,
					want:  testtarget.ErrInvalidOwner("invalid owner: space in the owner"),
				},
				{
					title: "valid-host,valid-owner,empty-name",
					input: host1 + "/" + owner1 + "/",
					want:  testtarget.ErrEmptyName,
				},
				{
					title: "valid-host,valid-owner,invalid-name",
					input: host1 + "/" + owner1 + "/space in the name",
					want:  testtarget.ErrInvalidName("invalid name: space in the name"),
				},
				{
					title: "valid-host,empty-owner,empty-name",
					input: host1 + "//",
					want:  testtarget.ErrEmptyName,
				},
				{
					title: "empty-host,valid-owner,empty-name",
					input: "/" + owner1 + "/",
					want:  testtarget.ErrEmptyName,
				},
				{
					title: "empty-host,empty-owner,valid-name",
					input: "//" + name,
					want:  testtarget.ErrEmptyOwner,
				},
				{
					title: "empty-host,empty-owner,empty-name",
					input: "//",
					want:  testtarget.ErrEmptyName,
				},
				{
					title: "unnecessary-following-slash",
					input: host1 + "/" + owner1 + "/" + name + "/",
					want:  testtarget.ErrTooManySlashes,
				},
				{
					title: "unnecessary-heading-slash",
					input: "/" + host1 + "/" + owner1 + "/" + name + "/",
					want:  testtarget.ErrTooManySlashes,
				},
			} {
				t.Run(testcase.title, func(t *testing.T) {
					spec, _, _, err := parser.ParseWithAlias(testcase.input)
					if err == nil {
						t.Fatalf("want failure to parse %q but parsed to %+v", testcase.input, spec)
					}
					if reflect.TypeOf(testcase.want) != reflect.TypeOf(err) {
						t.Fatalf("want error %t to parse %q but %t gotten", testcase.want, testcase.input, err)
					}
					if testcase.want.Error() != err.Error() {
						t.Fatalf("want error value %q to parse %q but %q gotten", testcase.want, testcase.input, err)
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

				wantHost string
				wantUser string
				wantName string

				wantServerUser string
				wantToken      string
			}{{
				title:          "valid-name",
				source:         name,
				wantHost:       host1,
				wantUser:       owner1,
				wantName:       name,
				wantServerUser: owner1,
				wantToken:      "token",
			}, {
				title:          "default-owner,valid-name",
				source:         owner1 + "/" + name,
				wantHost:       host1,
				wantUser:       owner1,
				wantName:       name,
				wantServerUser: owner1,
				wantToken:      "token",
			}, {
				title:          "default-host,default-owner,valid-name",
				source:         host1 + "/" + owner1 + "/" + name,
				wantHost:       host1,
				wantUser:       owner1,
				wantName:       name,
				wantServerUser: owner1,
				wantToken:      "token",
			}, {
				title:          "valid-host,valid-owner,valid-name",
				source:         host2 + "/" + owner2 + "/" + name,
				wantHost:       host2,
				wantUser:       owner2,
				wantServerUser: owner2,
				wantName:       name,
			}} {
				t.Run(testcase.title, func(t *testing.T) {
					spec, server, err := parser.Parse(testcase.source)
					if err != nil {
						t.Fatalf("failed to parse %q: %s", testcase.source, err)
					}
					if testcase.wantHost != spec.Host() {
						t.Errorf("want host %q but %q gotten", testcase.wantHost, spec.Host())
					}
					if testcase.wantUser != spec.Owner() {
						t.Errorf("want owner %q but %q gotten", testcase.wantUser, spec.Owner())
					}
					if testcase.wantName != spec.Name() {
						t.Errorf("want name %q but %q gotten", testcase.wantName, spec.Name())
					}
					if testcase.wantHost != server.Host() {
						t.Errorf("want host %q but %q gotten", testcase.wantHost, server.Host())
					}
					if testcase.wantServerUser != server.User() {
						t.Errorf("want owner %q but %q gotten", testcase.wantServerUser, server.User())
					}
					if testcase.wantToken != server.Token() {
						t.Errorf("want token %q but %q gotten", testcase.wantToken, server.Token())
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

			wantHost string
			wantUser string
			wantName string

			wantServerUser string
			wantToken      string
		}{{
			title:          "valid-name",
			source:         name,
			wantHost:       testtarget.DefaultHost,
			wantUser:       owner1,
			wantName:       name,
			wantServerUser: owner1,
			wantToken:      "token1",
		}, {
			title:          "valid-name,valid-owner",
			source:         owner2 + "/" + name,
			wantHost:       testtarget.DefaultHost,
			wantUser:       owner2,
			wantServerUser: owner1,
			wantName:       name,
			wantToken:      "token1",
		}, {
			title:          "full-name",
			source:         testtarget.DefaultHost + "/" + owner2 + "/" + name,
			wantHost:       testtarget.DefaultHost,
			wantUser:       owner2,
			wantName:       name,
			wantServerUser: owner1,
			wantToken:      "token1",
		}, {
			title:          "other-host",
			source:         host1 + "/" + owner2 + "/" + name,
			wantHost:       host1,
			wantUser:       owner2,
			wantName:       name,
			wantServerUser: owner2,
			wantToken:      "token2",
		}, {
			title:          "not-matched",
			source:         host2 + "/" + owner2 + "/" + name,
			wantHost:       host2,
			wantUser:       owner2,
			wantServerUser: owner2,
			wantName:       name,
			wantToken:      "",
		}} {
			t.Run(testcase.title, func(t *testing.T) {
				spec, server, err := parser.Parse(testcase.source)
				if err != nil {
					t.Fatalf("failed to parse %q: %s", testcase.source, err)
				}
				if testcase.wantHost != spec.Host() {
					t.Errorf("want host %q but %q gotten", testcase.wantHost, spec.Host())
				}
				if testcase.wantUser != spec.Owner() {
					t.Errorf("want owner %q but %q gotten", testcase.wantUser, spec.Owner())
				}
				if testcase.wantName != spec.Name() {
					t.Errorf("want name %q but %q gotten", testcase.wantName, spec.Name())
				}
				if testcase.wantHost != server.Host() {
					t.Errorf("want host %q but %q gotten", testcase.wantHost, server.Host())
				}
				if testcase.wantServerUser != server.User() {
					t.Errorf("want owner %q but %q gotten", testcase.wantServerUser, server.User())
				}
				if testcase.wantToken != server.Token() {
					t.Errorf("want token %q but %q gotten", testcase.wantToken, server.Token())
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
				title:    "with-alias-owner",
				source:   host1 + "/" + owner1 + "/" + name + "=" + owner2 + "/alias",
				wantHost: host1,
				wantUser: owner2,
				wantName: "alias",
			}, {
				title:    "with-alias-host",
				source:   host1 + "/" + owner1 + "/" + name + "=" + host2 + "/" + owner2 + "/alias",
				wantHost: host2,
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

		t.Run("WithInvalidAlias", func(t *testing.T) {
			for _, testcase := range []struct {
				title  string
				source string
			}{{
				title:  "empty",
				source: name + "=",
			}, {
				title:  "space",
				source: name + "= ",
			}, {
				title:  "space-in-the-alias",
				source: name + "=splitted name",
			}, {
				title:  "double-alias",
				source: name + "=alias1=alias2",
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
