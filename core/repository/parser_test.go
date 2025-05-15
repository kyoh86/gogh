package repository_test

import (
	"errors"
	"reflect"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/core/repository"
)

func TestReferenceParser(t *testing.T) {
	const (
		host0  = "github.com"
		owner1 = "kyoh86"
		owner2 = "anonymous"
		host1  = "example.com" // host a not default
		host2  = "kyoh86.dev"  // host a not default
		name   = "gogh"
	)
	t.Run("Empty", func(t *testing.T) {
		parser := testtarget.NewReferenceParser("", "")
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
		parser := testtarget.NewReferenceParser(host0, owner1)
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
						t.Fatalf("invalid %q: %s", testcase.source, err)
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
					expect: errors.New("'.' is reserved name"),
				},
				{
					title:  "valid-owner,dotdot",
					input:  owner1 + "/..",
					expect: errors.New("'..' is reserved name"),
				},
				{
					title:  "invalid-owner,valid-name",
					input:  "space in the owner/" + name,
					expect: errors.New("invalid owner: space in the owner"),
				},
				{
					title:  "valid-owner,empty-name",
					input:  owner1 + "/",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "valid-owner,invalid-name",
					input:  owner1 + "/space in the name",
					expect: errors.New("invalid name: space in the name"),
				},

				{
					title:  "empty-host,valid-owner,valid-name",
					input:  "/" + owner1 + "/" + name,
					expect: testtarget.ErrEmptyHost,
				},
				{
					title:  "invalid-host,valid-owner,valid-name",
					input:  "space in the host/" + owner1 + "/" + name,
					expect: errors.New("invalid host: space in the host"),
				},
				{
					title:  "valid-host,empty-owner,valid-name",
					input:  host1 + "//" + name,
					expect: testtarget.ErrEmptyOwner,
				},
				{
					title:  "valid-host,invalid-owner,valid-name",
					input:  host1 + "/space in the owner/" + name,
					expect: errors.New("invalid owner: space in the owner"),
				},
				{
					title:  "valid-host,valid-owner,empty-name",
					input:  host1 + "/" + owner1 + "/",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "valid-host,valid-owner,invalid-name",
					input:  host1 + "/" + owner1 + "/space in the name",
					expect: errors.New("invalid name: space in the name"),
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
					res, err := parser.ParseWithAlias(testcase.input)
					if err == nil {
						t.Fatalf(
							"expect failure to parse %q but parsed to %+v",
							testcase.input,
							res.Reference,
						)
					}
					if reflect.TypeOf(testcase.expect) != reflect.TypeOf(err) {
						t.Fatalf(
							"expect error %T to parse %q but %T gotten",
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
		parser := testtarget.NewReferenceParser(host1, owner1)
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
						t.Fatalf("invalid %q: %s", testcase.source, err)
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
		parser := testtarget.NewReferenceParser("default-host", "default-owner")
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
				res, err := parser.ParseWithAlias(testcase.source)
				if err != nil {
					t.Fatalf("invalid %q: %s", testcase.source, err)
				}
				if res.Alias != nil {
					t.Errorf("want alias is nil but %#v gotten", res.Alias)
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
					res, err := parser.ParseWithAlias(testcase.source)
					if err != nil {
						t.Fatalf("invalid %q: %s", testcase.source, err)
					}
					if res.Alias == nil {
						t.Fatal("want valid alias but got nil")
					}
					if testcase.wantHost != res.Alias.Host() {
						t.Errorf("want host %q but %q gotten", testcase.wantHost, res.Alias.Host())
					}
					if testcase.wantRepoUser != res.Reference.Owner() {
						t.Errorf("want repo owner %q but %q gotten", testcase.wantRepoUser, res.Reference.Owner())
					}
					if testcase.wantRepoName != res.Reference.Name() {
						t.Errorf("want repo name %q but %q gotten", testcase.wantRepoName, res.Reference.Name())
					}
					if testcase.wantAliasUser != res.Alias.Owner() {
						t.Errorf("want res.Alias owner %q but %q gotten", testcase.wantAliasUser, res.Alias.Owner())
					}
					if testcase.wantAliasName != res.Alias.Name() {
						t.Errorf("want res.Alias name %q but %q gotten", testcase.wantAliasName, res.Alias.Name())
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
					_, err := parser.ParseWithAlias(testcase.source)
					if err == nil {
						t.Fatal("want error, but got nil")
					}
					t.Log(err)
				})
			}
		})
	})
}
