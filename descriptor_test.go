package gogh_test

import (
	"context"
	"reflect"
	"testing"

	testtarget "github.com/kyoh86/gogh/v2"
)

func TestDescriptor(t *testing.T) {
	ctx := context.Background()
	t.Run("Empty", func(t *testing.T) {
		descriptor := testtarget.NewDescriptor(ctx)
		t.Run("ValidInput", func(t *testing.T) {
			for _, testcase := range []struct {
				source string
				expect testtarget.Description
			}{{
				source: "kyoh86/gogh",
				expect: description(t, "github.com", "kyoh86", "gogh"),
			}, {
				source: "github.com/kyoh86/gogh",
				expect: description(t, "github.com", "kyoh86", "gogh"),
			}, {
				source: "example.com/kyoh86/gogh",
				expect: description(t, "example.com", "kyoh86", "gogh"),
			}} {
				t.Run(testcase.source, func(t *testing.T) {
					description, err := descriptor.Parse(testcase.source)
					if err != nil {
						t.Fatalf("failed to parse %q: %s", testcase.source, err)
					}
					if testcase.expect.Host() != description.Host() {
						t.Errorf("expect host %q but %q gotten", testcase.expect.Host(), description.Host())
					}
					if testcase.expect.User() != description.User() {
						t.Errorf("expect user %q but %q gotten", testcase.expect.User(), description.User())
					}
					if testcase.expect.Name() != description.Name() {
						t.Errorf("expect name %q but %q gotten", testcase.expect.Name(), description.Name())
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
					expect: testtarget.ErrEmptyUser,
				},
				{
					title:  "valid-name",
					input:  "gogh", // shortage
					expect: testtarget.ErrEmptyUser,
				},
				{
					title:  "empty-user,empty-name",
					input:  "/",
					expect: testtarget.ErrEmptyUser,
				},
				{
					title:  "empty-user,valid-name",
					input:  "/gogh",
					expect: testtarget.ErrEmptyUser,
				},
				{
					title:  "valid-user,dot",
					input:  "kyoh86/.",
					expect: testtarget.ErrInvalidName("'.' is reserved name"),
				},
				{
					title:  "valid-user,dotdot",
					input:  "kyoh86/..",
					expect: testtarget.ErrInvalidName("'..' is reserved name"),
				},
				{
					title:  "invalid-user,valid-name",
					input:  "space in the user/gogh",
					expect: testtarget.ErrInvalidUser("invalid user: space in the user"),
				},
				{
					title:  "valid-user,empty-name",
					input:  "kyoh86/",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "valid-user,invalid-name",
					input:  "kyoh86/space in the name",
					expect: testtarget.ErrInvalidName("invalid name: space in the name"),
				},

				{
					title:  "empty-host,valid-user,valid-name",
					input:  "/kyoh86/gogh",
					expect: testtarget.ErrEmptyHost,
				},
				{
					title:  "invalid-host,valid-user,valid-name",
					input:  "space in the host/kyoh86/gogh",
					expect: testtarget.ErrInvalidHost("invalid host: space in the host"),
				},
				{
					title:  "valid-host,empty-user,valid-name",
					input:  "example.com//gogh",
					expect: testtarget.ErrEmptyUser,
				},
				{
					title:  "valid-host,invalid-user,valid-name",
					input:  "example.com/space in the user/gogh",
					expect: testtarget.ErrInvalidUser("invalid user: space in the user"),
				},
				{
					title:  "valid-host,valid-user,empty-name",
					input:  "example.com/kyoh86/",
					expect: testtarget.ErrEmptyName,
				},
				{
					title:  "valid-host,valid-user,invalid-name",
					input:  "example.com/kyoh86/space in the name",
					expect: testtarget.ErrInvalidName("invalid name: space in the name"),
				},
				{
					title:  "valid-host,empty-user,empty-name",
					input:  "example.com//",
					expect: testtarget.ErrEmptyUser,
				},
				{
					title:  "empty-host,valid-user,empty-name",
					input:  "/kyoh86/",
					expect: testtarget.ErrEmptyHost,
				},
				{
					title:  "empty-host,empty-user,valid-name",
					input:  "//gogh",
					expect: testtarget.ErrEmptyHost,
				},
				{
					title:  "empty-host,empty-user,empty-name",
					input:  "//",
					expect: testtarget.ErrEmptyHost,
				},
				{
					title:  "unnecessary-following-slash",
					input:  "example.com/kyoh86/gogh/",
					expect: testtarget.ErrTooManySlashes,
				},
				{
					title:  "unnecessary-heading-slash",
					input:  "/example.com/kyoh86/gogh/",
					expect: testtarget.ErrTooManySlashes,
				},
			} {
				t.Run(testcase.title, func(t *testing.T) {
					description, err := descriptor.Parse(testcase.input)
					if err == nil {
						t.Fatalf("expect failure to parse %q but parsed to %+v", testcase.input, description)
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

	t.Run("WithDefaultUser", func(t *testing.T) {
		descriptor := testtarget.NewDescriptor(ctx)
		if err := descriptor.SetDefaultUser("invalid user"); err == nil {
			t.Error("expect failure for set invalid default user")
		}

		if err := descriptor.SetDefaultUser("kyoh86"); err != nil {
			t.Fatalf("failed to set default user: %q", err)
		}

		for _, testcase := range []struct {
			source string
			expect testtarget.Description
		}{{
			source: "gogh",
			expect: description(t, "github.com", "kyoh86", "gogh"),
		}, {
			source: "kyoh86/gogh",
			expect: description(t, "github.com", "kyoh86", "gogh"),
		}, {
			source: "example/gogh",
			expect: description(t, "github.com", "example", "gogh"),
		}, {
			source: "github.com/example/gogh",
			expect: description(t, "github.com", "example", "gogh"),
		}, {
			source: "example.com/example/gogh",
			expect: description(t, "example.com", "example", "gogh"),
		}} {
			t.Run(testcase.source, func(t *testing.T) {
				description, err := descriptor.Parse(testcase.source)
				if err != nil {
					t.Fatalf("failed to parse %q: %s", testcase.source, err)
				}
				if testcase.expect.Host() != description.Host() {
					t.Errorf("expect host %q but %q gotten", testcase.expect.Host(), description.Host())
				}
				if testcase.expect.User() != description.User() {
					t.Errorf("expect user %q but %q gotten", testcase.expect.User(), description.User())
				}
				if testcase.expect.Name() != description.Name() {
					t.Errorf("expect name %q but %q gotten", testcase.expect.Name(), description.Name())
				}
			})
		}
	})
}
