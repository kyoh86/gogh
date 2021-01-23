package gogh_test

import (
	"context"
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
					desc, err := descriptor.Parse(testcase.source)
					if err != nil {
						t.Fatalf("failed to parse %q: %s", testcase.source, err)
					}
					if testcase.expect.Host() != desc.Host() {
						t.Errorf("expect host %q but %q gotten", testcase.expect.Host(), desc.Host())
					}
					if testcase.expect.User() != desc.User() {
						t.Errorf("expect user %q but %q gotten", testcase.expect.User(), desc.User())
					}
					if testcase.expect.Name() != desc.Name() {
						t.Errorf("expect name %q but %q gotten", testcase.expect.Name(), desc.Name())
					}
				})
			}
		})

		t.Run("InvalidInput", func(t *testing.T) {
			for _, source := range []string{
				"",
				"invalid name",

				"/",
				"/gogh",
				"invalid user/gogh",
				"kyoh86/",
				"kyoh86/invalid name",

				"/kyoh86/gogh",
				"invalid host/kyoh86/gogh",
				"github.com//gogh",
				"github.com/invalid user/gogh",
				"github.com/kyoh86/",
				"github.com/kyoh86/invalid name",
				"github.com//",
				"/kyoh86/",
				"//gogh",
				"//",
				"github.com/kyoh86/gogh/",
				"/github.com/kyoh86/gogh/",

				"gogh", // shortage
			} {
				t.Run(source, func(t *testing.T) {
					desc, err := descriptor.Parse(source)
					if err == nil {
						t.Errorf("expect failure for parse %q but parsed to %+v", source, desc)
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
				desc, err := descriptor.Parse(testcase.source)
				if err != nil {
					t.Fatalf("failed to parse %q: %s", testcase.source, err)
				}
				if testcase.expect.Host() != desc.Host() {
					t.Errorf("expect host %q but %q gotten", testcase.expect.Host(), desc.Host())
				}
				if testcase.expect.User() != desc.User() {
					t.Errorf("expect user %q but %q gotten", testcase.expect.User(), desc.User())
				}
				if testcase.expect.Name() != desc.Name() {
					t.Errorf("expect name %q but %q gotten", testcase.expect.Name(), desc.Name())
				}
			})
		}
	})
}
