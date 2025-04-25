package reporef_test

import (
	"reflect"
	"testing"

	testtarget "github.com/kyoh86/gogh/v3/domain/reporef"
)

func TestRepoRefParser(t *testing.T) {
	const (
		host0  = "github.com"
		owner1 = "kyoh86"
		owner2 = "anonymous"
		host1  = "example.com" // host a not default
		host2  = "kyoh86.dev"  // host a not default
		name   = "gogh"
	)
	t.Run("Empty", func(t *testing.T) {
		var parser testtarget.RepoRefParser
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
				if _, err := parser.Parse(testcase.source); err == nil {
					t.Fatal("expect error, but nil")
				}
			})
		}
	})

	t.Run("DefaultHost", func(t *testing.T) {
		parser := testtarget.NewRepoRefParser(host0, owner1)
		t.Run("ValidInput", func(t *testing.T) {
			for _, testcase := range []struct {
				title  string
				source string

				expectHost string
				expectUser string
				expectName string
			}{{
				title:      "valid-name",
				source:     name,
				expectHost: host0,
				expectUser: owner1,
				expectName: name,
			}, {
				title:      "default-owner,valid-name",
				source:     owner1 + "/" + name,
				expectHost: host0,
				expectUser: owner1,
				expectName: name,
			}, {
				title:      "default-host,default-owner,valid-name",
				source:     host0 + "/" + owner1 + "/" + name,
				expectHost: host0,
				expectUser: owner1,
				expectName: name,
			}, {
				title:      "valid-host,valid-owner,valid-name",
				source:     host1 + "/" + owner2 + "/" + name,
				expectHost: host1,
				expectUser: owner2,
				expectName: name,
			}} {
				t.Run(testcase.title, func(t *testing.T) {
					ref, err := parser.Parse(testcase.source)
					if err != nil {
						t.Fatalf("failed to parse %q: %s", testcase.source, err)
					}
					if testcase.expectHost != ref.Host() {
						t.Errorf("expect host %q but %q gotten", testcase.expectHost, ref.Host())
					}
					if testcase.expectUser != ref.Owner() {
						t.Errorf("expect owner %q but %q gotten", testcase.expectUser, ref.Owner())
					}
					if testcase.expectName != ref.Name() {
						t.Errorf("expect name %q but %q gotten", testcase.expectName, ref.Name())
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
					ref, _, err := parser.ParseWithAlias(testcase.input)
					if err == nil {
						t.Fatalf(
							"expect failure to parse %q but parsed to %+v",
							testcase.input,
							ref,
						)
					}
					if reflect.TypeOf(testcase.expect) != reflect.TypeOf(err) {
						t.Fatalf(
							"expect error %t to parse %q but %t gotten",
							testcase.expect,
							testcase.input,
							err,
						)
					}
					if testcase.expect.Error() != err.Error() {
						t.Fatalf(
							"expect error value %q to parse %q but %q gotten",
							testcase.expect,
							testcase.input,
							err,
						)
					}
				})
			}
		})
	})

	t.Run("WithHost", func(t *testing.T) {
		parser := testtarget.NewRepoRefParser(host1, owner1)
		t.Run("ValidInput", func(t *testing.T) {
			for _, testcase := range []struct {
				title  string
				source string

				expectHost string
				expectUser string
				expectName string
			}{{
				title:      "valid-name",
				source:     name,
				expectHost: host1,
				expectUser: owner1,
				expectName: name,
			}, {
				title:      "default-owner,valid-name",
				source:     owner1 + "/" + name,
				expectHost: host1,
				expectUser: owner1,
				expectName: name,
			}, {
				title:      "default-host,default-owner,valid-name",
				source:     host1 + "/" + owner1 + "/" + name,
				expectHost: host1,
				expectUser: owner1,
				expectName: name,
			}, {
				title:      "valid-host,valid-owner,valid-name",
				source:     host2 + "/" + owner2 + "/" + name,
				expectHost: host2,
				expectUser: owner2,
				expectName: name,
			}} {
				t.Run(testcase.title, func(t *testing.T) {
					ref, err := parser.Parse(testcase.source)
					if err != nil {
						t.Fatalf("failed to parse %q: %s", testcase.source, err)
					}
					if testcase.expectHost != ref.Host() {
						t.Errorf("expect host %q but %q gotten", testcase.expectHost, ref.Host())
					}
					if testcase.expectUser != ref.Owner() {
						t.Errorf("expect owner %q but %q gotten", testcase.expectUser, ref.Owner())
					}
					if testcase.expectName != ref.Name() {
						t.Errorf("expect name %q but %q gotten", testcase.expectName, ref.Name())
					}
				})
			}
		})
	})

	t.Run("ParseWithAlias", func(t *testing.T) {
		parser := testtarget.NewRepoRefParser("default-host", "default-owner")
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
				_, alias, err := parser.ParseWithAlias(testcase.source)
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

				wantHost      string
				wantRepoUser  string
				wantRepoName  string
				wantAliasUser string
				wantAliasName string
			}{{
				title:         "with-alias-name",
				source:        host1 + "/" + owner1 + "/" + name + "=alias",
				wantHost:      host1,
				wantRepoUser:  owner1,
				wantRepoName:  name,
				wantAliasUser: owner1,
				wantAliasName: "alias",
			}, {
				title:         "with-not-default-repo-with-alias-name",
				source:        host2 + "/" + owner2 + "/" + name + "=alias",
				wantHost:      host2,
				wantRepoUser:  owner2,
				wantRepoName:  name,
				wantAliasUser: owner2,
				wantAliasName: "alias",
			}, {
				title:         "with-alias-owner",
				source:        host1 + "/" + owner1 + "/" + name + "=" + owner2 + "/alias",
				wantHost:      host1,
				wantRepoUser:  owner1,
				wantRepoName:  name,
				wantAliasUser: owner2,
				wantAliasName: "alias",
			}} {
				t.Run(testcase.title, func(t *testing.T) {
					repo, alias, err := parser.ParseWithAlias(testcase.source)
					if err != nil {
						t.Fatalf("failed to parse %q: %s", testcase.source, err)
					}
					if alias == nil {
						t.Fatal("want valid alisa but got nil")
					}
					if testcase.wantHost != alias.Host() {
						t.Errorf("want host %q but %q gotten", testcase.wantHost, alias.Host())
					}
					if testcase.wantRepoUser != repo.Owner() {
						t.Errorf("want repo owner %q but %q gotten", testcase.wantRepoUser, alias.Owner())
					}
					if testcase.wantRepoName != repo.Name() {
						t.Errorf("want repo name %q but %q gotten", testcase.wantRepoName, alias.Name())
					}
					if testcase.wantAliasUser != alias.Owner() {
						t.Errorf("want alias owner %q but %q gotten", testcase.wantAliasUser, alias.Owner())
					}
					if testcase.wantAliasName != alias.Name() {
						t.Errorf("want alias name %q but %q gotten", testcase.wantAliasName, alias.Name())
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
					_, _, err := parser.ParseWithAlias(testcase.source)
					if err == nil {
						t.Fatal("want error, but got nil")
					}
					t.Log(err)
				})
			}
		})
	})
}
