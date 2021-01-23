package gogh_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	testtarget "github.com/kyoh86/gogh/v2"
)

func TestFormat(t *testing.T) {
	ctx := context.Background()
	local := testtarget.NewLocalController(ctx, filepath.Join("/", "tmp"))
	desc, err := testtarget.ValidateDescription("github.com", "kyoh86", "gogh")
	if err != nil {
		t.Fatalf("failed to init Description: %s", err)
	}
	// NOTE: when the path is checked, it should be passed with filepath.Clean.
	// Because windows uses '\' for path separator.
	project, err := local.Get(ctx, *desc)
	if err != nil {
		t.Fatalf("failed to get project from Description: %s", err)
	}
	for _, testcase := range []struct {
		title  string
		format testtarget.Format
		expect string
	}{
		{
			title:  "FullPath",
			format: testtarget.FormatFullPath,
			expect: filepath.Clean("/tmp/github.com/kyoh86/gogh"),
		},
		{
			title:  "RelPath",
			format: testtarget.FormatRelPath,
			expect: filepath.Clean("github.com/kyoh86/gogh"),
		},
		{
			title:  "URL",
			format: testtarget.FormatURL,
			expect: "https://github.com/kyoh86/gogh",
		},
		{
			title:  "Fields",
			format: testtarget.FormatFields,
			expect: strings.Join([]string{
				filepath.Clean("/tmp/github.com/kyoh86/gogh"),
				"https://github.com/kyoh86/gogh",
				"github.com",
				"kyoh86",
				"gogh",
			}, " "),
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			actual := format.Format(project)
			if testcase.expect != actual {
				t.Errorf("expect %q but %q is gotten", testcase.expect, actual)
			}
		})
	}
	// UNDONE: Test JSON Formatter: it should be checked with assert.JSONEq
}
