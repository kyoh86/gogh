package gogh_test

import (
	"testing"

	testtarget "github.com/kyoh86/gogh/v2"
)

func TestSpec(t *testing.T) {
	const (
		validHost = "example.com"
		validUser = "kyoh86"
		validName = "gogh"
	)
	for _, testcase := range []struct {
		title  string
		host   string
		user   string
		name   string
		expect error
	}{
		{
			title:  "valid",
			host:   validHost,
			user:   validUser,
			name:   validName,
			expect: nil,
		},
		{
			title:  "empty-name",
			host:   validHost,
			user:   validUser,
			name:   "",
			expect: testtarget.ErrEmptyName,
		},
		{
			title:  "empty-user",
			host:   validHost,
			user:   "",
			name:   validName,
			expect: testtarget.ErrEmptyUser,
		},
		{
			title:  "empty-host",
			host:   "",
			user:   validUser,
			name:   validName,
			expect: testtarget.ErrEmptyHost,
		},
		{
			title:  "slashed-name",
			host:   validHost,
			user:   validUser,
			name:   "go/gh",
			expect: testtarget.ErrInvalidName("invalid name: go/gh"),
		},
		{
			title:  "slashed-user",
			host:   validHost,
			user:   "kyoh/86",
			name:   validName,
			expect: testtarget.ErrInvalidUser("invalid user: kyoh/86"),
		},
		{
			title:  "slashed-host",
			host:   "example.com/example",
			user:   validUser,
			name:   validName,
			expect: testtarget.ErrInvalidHost("invalid host: example.com/example"),
		},
		{
			title:  "dotted-user",
			host:   validHost,
			user:   "kyoh.86",
			name:   validName,
			expect: testtarget.ErrInvalidUser("invalid user: kyoh.86"),
		},
		{
			title:  "dot-name",
			host:   validHost,
			user:   validUser,
			name:   ".",
			expect: testtarget.ErrInvalidName("'.' is reserved name"),
		},
		{
			title:  "dot-user",
			host:   validHost,
			user:   ".",
			name:   validName,
			expect: testtarget.ErrInvalidUser("invalid user: ."),
		},
		{
			title:  "dotdot-name",
			host:   validHost,
			user:   validUser,
			name:   "..",
			expect: testtarget.ErrInvalidName("'..' is reserved name"),
		},
		{
			title:  "dotdot-user",
			host:   validHost,
			user:   "..",
			name:   validName,
			expect: testtarget.ErrInvalidUser("invalid user: .."),
		},
		{
			title:  "ported-host",
			host:   "127.0.0.1:9000",
			user:   validUser,
			name:   validName,
			expect: nil,
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			spec, err := testtarget.NewSpec(testcase.host, testcase.user, testcase.name)
			if testcase.expect == nil {
				if err != nil {
					t.Fatalf("failed to create new spec: %s", err)
				}
				if testcase.host != spec.Host() {
					t.Errorf("expect host %q but %q gotten", testcase.host, spec.Host())
				}
				if testcase.user != spec.User() {
					t.Errorf("expect user %q but %q gotten", testcase.user, spec.User())
				}
				if testcase.name != spec.Name() {
					t.Errorf("expect name %q but %q gotten", testcase.name, spec.Name())
				}
			} else {
				if err == nil {
					t.Fatal("expect failure to create new spec, but not")
				}
				if testcase.expect.Error() != err.Error() {
					t.Fatalf("expect error %s, but actual %s is gottten", testcase.expect, err)
				}
			}
		})
	}
}
