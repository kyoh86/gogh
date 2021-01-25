package gogh_test

import (
	"testing"

	testtarget "github.com/kyoh86/gogh/v2"
)

func TestDescription(t *testing.T) {
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
			title:  "slash-name",
			host:   validHost,
			user:   validUser,
			name:   "go/gh",
			expect: testtarget.ErrInvalidName("invalid name: go/gh"),
		},
		{
			title:  "slash-user",
			host:   validHost,
			user:   "kyoh/86",
			name:   validName,
			expect: testtarget.ErrInvalidUser("invalid user: kyoh/86"),
		},
		{
			title:  "slash-host",
			host:   "example.com/example",
			user:   validUser,
			name:   validName,
			expect: testtarget.ErrInvalidHost("invalid host: example.com/example"),
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
			d, err := testtarget.NewDescription(testcase.host, testcase.user, testcase.name)
			if testcase.expect == nil {
				if err != nil {
					t.Fatalf("failed to create new description: %s", err)
				}
				if testcase.host != d.Host() {
					t.Errorf("expect host %q but %q gotten", testcase.host, d.Host())
				}
				if testcase.user != d.User() {
					t.Errorf("expect user %q but %q gotten", testcase.user, d.User())
				}
				if testcase.name != d.Name() {
					t.Errorf("expect name %q but %q gotten", testcase.name, d.Name())
				}
			} else {
				if err == nil {
					t.Fatal("expect failure to create new description, but not")
				}
				if testcase.expect.Error() != err.Error() {
					t.Fatalf("expect error %s, but actual %s is gottten", testcase.expect, err)
				}
			}
		})
	}
}
