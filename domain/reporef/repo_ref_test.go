package reporef_test

import (
	"testing"

	testtarget "github.com/kyoh86/gogh/v3/domain/reporef"
)

func TestRepoRef(t *testing.T) {
	const (
		validHost  = "example.com"
		validOwner = "kyoh86"
		validName  = "gogh"
	)
	for _, testcase := range []struct {
		title  string
		host   string
		owner  string
		name   string
		expect error
	}{
		{
			title:  "valid",
			host:   validHost,
			owner:  validOwner,
			name:   validName,
			expect: nil,
		},
		{
			title:  "empty-name",
			host:   validHost,
			owner:  validOwner,
			name:   "",
			expect: testtarget.ErrEmptyName,
		},
		{
			title:  "empty-owner",
			host:   validHost,
			owner:  "",
			name:   validName,
			expect: testtarget.ErrEmptyOwner,
		},
		{
			title:  "empty-host",
			host:   "",
			owner:  validOwner,
			name:   validName,
			expect: testtarget.ErrEmptyHost,
		},
		{
			title:  "slashed-name",
			host:   validHost,
			owner:  validOwner,
			name:   "go/gh",
			expect: testtarget.ErrInvalidName("invalid name: go/gh"),
		},
		{
			title:  "slashed-owner",
			host:   validHost,
			owner:  "kyoh/86",
			name:   validName,
			expect: testtarget.ErrInvalidOwner("invalid owner: kyoh/86"),
		},
		{
			title:  "slashed-host",
			host:   "example.com/example",
			owner:  validOwner,
			name:   validName,
			expect: testtarget.ErrInvalidHost("invalid host: example.com/example"),
		},
		{
			title:  "dotted-owner",
			host:   validHost,
			owner:  "kyoh.86",
			name:   validName,
			expect: testtarget.ErrInvalidOwner("invalid owner: kyoh.86"),
		},
		{
			title:  "dot-name",
			host:   validHost,
			owner:  validOwner,
			name:   ".",
			expect: testtarget.ErrInvalidName("'.' is reserved name"),
		},
		{
			title:  "dot-owner",
			host:   validHost,
			owner:  ".",
			name:   validName,
			expect: testtarget.ErrInvalidOwner("invalid owner: ."),
		},
		{
			title:  "dotdot-name",
			host:   validHost,
			owner:  validOwner,
			name:   "..",
			expect: testtarget.ErrInvalidName("'..' is reserved name"),
		},
		{
			title:  "dotdot-owner",
			host:   validHost,
			owner:  "..",
			name:   validName,
			expect: testtarget.ErrInvalidOwner("invalid owner: .."),
		},
		{
			title:  "ported-host",
			host:   "127.0.0.1:9000",
			owner:  validOwner,
			name:   validName,
			expect: nil,
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			ref, err := testtarget.NewRepoRef(testcase.host, testcase.owner, testcase.name)
			if testcase.expect == nil {
				if err != nil {
					t.Fatalf("failed to create new ref: %s", err)
				}
				if testcase.host != ref.Host() {
					t.Errorf("expect host %q but %q gotten", testcase.host, ref.Host())
				}
				if testcase.owner != ref.Owner() {
					t.Errorf("expect owner %q but %q gotten", testcase.owner, ref.Owner())
				}
				if testcase.name != ref.Name() {
					t.Errorf("expect name %q but %q gotten", testcase.name, ref.Name())
				}
			} else {
				if err == nil {
					t.Fatal("expect failure to create new ref, but not")
				}
				if testcase.expect.Error() != err.Error() {
					t.Fatalf("expect error %s, but actual %s is gottten", testcase.expect, err)
				}
			}
		})
	}
}
